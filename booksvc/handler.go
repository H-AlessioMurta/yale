package booksvc


import (
    //"io/ioutil"
    "encoding/json"
    "net/http"
    "database/sql"
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)


type Server struct {
	Router *mux.Router
	DB     *sql.DB
}


func (s *Server) initializeRoutes() {
    logInfo("Initializing http routes")
    s.Router.HandleFunc("/books", s.getBooks).Methods("GET")
    s.Router.HandleFunc("/books", s.postBook).Methods("POST")
    s.Router.HandleFunc("/books/{id}", s.getBook).Methods("GET")
    s.Router.HandleFunc("/books/{id}", s.putBook).Methods("PUT")
    s.Router.HandleFunc("/books/{id}", s.deleteBook).Methods("DELETE")
    logInfo("CRUD's path are ready ")
}




// BookService provides operations on DB.
type BookService interface {
	postBook(w http.ResponseWriter, r *http.Request) 
	getBook(w http.ResponseWriter, r *http.Request) 
	putBook(w http.ResponseWriter, r *http.Request) 
	deleteBook(w http.ResponseWriter, r *http.Request)
	getBooks(w http.ResponseWriter, r *http.Request)
}

type bookService struct{}

func (s *Server) getBook(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    vars := mux.Vars(r)
    id:= vars["id"]
    b := Book{ID: id}
    if err := b.getBook(s.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "Book not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }
    respondWithJSON(w, http.StatusOK, b)
}

func (s *Server) getBooks(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    books, err := getBooks(s.DB)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusOK, books)
}

func (s *Server) postBook(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    var b Book
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&b); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()
    if err := b.postBook(s.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusCreated, b)
}

func (s *Server) putBook(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    vars := mux.Vars(r)
    id := vars["id"]
    if id == "" {
        respondWithError(w, http.StatusBadRequest, "Invalid Book ID")
        return
    }
    var b Book
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&b); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
        return
    }
    defer r.Body.Close()
    b.ID= id
    if err := b.putBook(s.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, b)
}

func (s *Server) deleteBook(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    vars := mux.Vars(r)
    id := vars["id"]
    if id == "" {
        respondWithError(w, http.StatusBadRequest, "Invalid Book ID")
        return
    }

    b := Book{ID: id}
    if err := b.deleteBook(s.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

/* v1
func GetBooks(w http.ResponseWriter, r *http.Request) {
    db := setupDB()
    printMessage("Getting all Books...")
    // Get all books from books table 
    rows, err := db.Query("SELECT * FROM books")
    // check errors
    checkErr(err)
    // var response []JsonResponse
    var booksList []Book
    // Foreach book
    for rows.Next() {
        var bookID string
        var bookTitle string
		var bookAuthors string
        err = rows.Scan(&bookID, &bookTitle, &bookAuthors)
        // check errors
        checkErr(err)
        booksList = append(booksList, Book{ID: bookID, Title: bookTitle, Authors: bookAuthors })
    }
    var response = JsonResponse{Type: "success", Data: booksList}

    json.NewEncoder(w).Encode(response)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
    db := setupDB()
	
	bookID := r.FormValue("id")
    // Get all books from books table 
    row, err := db.Query("SELECT * FROM books where id = $1", bookID)
    // check errors
    checkErr(err)
    // var response []JsonResponse
    var abook Book
    for row.Next(){
	    err = row.Scan(&abook.ID, &abook.Title, &abook.Authors)
        // check errors
	    checkErr(err)
        fmt.Println("porcazozza")
    }
	

	var bk []Book
	bk=append(bk,abook)

    var response = JsonResponse{Type: "success", Data: bk}

    json.NewEncoder(w).Encode(response)
}

func PostBook(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    bookID := uuid.New().String()
    bookTitle := r.Form.Get("title")
	bookAuthors := r.FormValue("authors")
    printMessage(bookID)
    printMessage(bookTitle)
    printMessage(bookAuthors)

    var response = JsonResponse{}

    if bookTitle == "" || bookAuthors == "" {
        response = JsonResponse{Type: "error", Message: "You are missing s valid ID for the Book or s for value title or author"}
    } else {
        db := setupDB()

        printMessage("Inserting the new book into DB")

        fmt.Println("Inserting new book with ID: " + bookID+ " and name: " + bookTitle + " written by:" + bookAuthors )
        var lastInsertID string
        err := db.QueryRow("INSERT INTO books(id,title,authors) VALUES($1, $2, $3) returning id;", bookID, bookTitle, bookAuthors ).Scan(&lastInsertID)
        // check errors
        checkErr(err)
        response = JsonResponse{Type: "success", Message: "The book has been inserted successfully!"}
    }
    json.NewEncoder(w).Encode(response)
}

func PutBook(w http.ResponseWriter, r *http.Request) {
    bookID := r.FormValue("id")
    bookTitle := r.FormValue("title")
	bookAuthors := r.FormValue("authors")

    var response = JsonResponse{}
    var con string
    if bookTitle == "" && bookAuthors == "" {
        response = JsonResponse{Type: "error", Message: "You are missing s valid ID for the Book or s for value title or author"}
    } else {
        db := setupDB()
        printMessage("Updating the book value into DB")
        //fmt.Println("Updating book with ID: " + bookID + " and title: " + bookTitle + " written by:" + bookAuthors )
        
        if bookTitle == "" && bookAuthors != ""{
            con = "UPDATE books set authors=$2 where id =$1"
            _, err := db.Exec(con, bookID, bookAuthors)
            checkErr(err)
            fmt.Println("Updating book with ID: " + bookID + " and title: " + bookTitle + " written by:" + bookAuthors )
        }
        if bookAuthors == "" && bookTitle != ""{
            con = "UPDATE books set title = $2 where id =$1"
            _, err := db.Exec(con, bookID, bookTitle )
            checkErr(err)
            fmt.Println("Updating book with ID: " + bookID + " and title: " + bookTitle + " written by:" + bookAuthors )
        }
        if bookAuthors != "" && bookTitle != ""{
            con = "UPDATE books set title = $2, authors=$3 where id =$1"
            _, err := db.Exec(con, bookID, bookTitle, bookAuthors )
            checkErr(err)
            fmt.Println("Updating book with ID: " + bookID + " and title: " + bookTitle + " written by:" + bookAuthors )
        }  
		// check errors
		
		response = JsonResponse{Type: "success", Message: "The book has been putd successfully!"}
    }
    json.NewEncoder(w).Encode(response)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    bookID := params["id"]
    var response = JsonResponse{}

    if bookID == "" {
        response = JsonResponse{Type: "error", Message: "You are missing bookID parameter."}
    } else {
        db := setupDB()
        printMessage("Deleting book from DB")
        _, err := db.Exec("DELETE FROM books where id = $1", bookID)
        // check errors
        checkErr(err)
        response = JsonResponse{Type: "success", Message: "The book has been deleted successfully!"}
    }
    json.NewEncoder(w).Encode(response)
}
*/