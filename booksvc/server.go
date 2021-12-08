package booksvc


import (
	"errors"
	"strings"
	"context"
	"database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
	"github.com/google/uuid"
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type Book struct {
	ID        string    `json:"id"`
	Title      string    `json:"title,omitempty"`
	Authors string `json:"authors,omitempty"`
}


// BookService provides operations on DB.
type BookService interface {
	PostBook(w http.ResponseWriter, r *http.Request) 
	GetBook(w http.ResponseWriter, r *http.Request) 
	PutBook(w http.ResponseWriter, r *http.Request) 
	DeleteBook(w http.ResponseWriter, r *http.Request)
	GetBooks(w http.ResponseWriter, r *http.Request)
}

type bookService struct{}

func GetBooks(w http.ResponseWriter, r *http.Request) {
    db := setupDB()

    printMessage("Getting all Books...")

    // Get all books from books table 
    rows, err := db.Query("SELECT * FROM books")

    // check errors
    checkErr(err)
	
    // var response []JsonResponse
    var books []Book

    // Foreach book
    for rows.Next() {
        
        var bookID string
        var bookTitle string
		var bookAuthors string

        err = rows.Scan(&bookID, &bookTitle, &bookAuthors)

        // check errors
        checkErr(err)

        books = append(books, Book{id: bookID, title: bookTitle, authors: bookAuthors })
    }

    var response = JsonResponse{Type: "success", Data: books}

    json.NewEncoder(w).Encode(response)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
    db := setupDB()
	
    printMessage("Getting all Books...")
	bookID := r.FormValue("id")
    // Get all books from books table 
    row, err := db.Query("SELECT * FROM books where bookID = $1", bookID)

    // check errors
    checkErr(err)
	
    // var response []JsonResponse
    var abook Book
	var bookID string
	var bookTitle string
	var bookAuthors string

	err = row.Scan(&bookID, &bookTitle, &bookAuthors)

	// check errors
	checkErr(err)

	abook = append(books, Book{id: bookID, title: bookTitle, authors: bookAuthors })
    }

    var response = JsonResponse{Type: "success", Data: abook}

    json.NewEncoder(w).Encode(response)
}


func PostBook(w http.ResponseWriter, r *http.Request) {
    bookID := uuid.New()
    bookTitle := r.FormValue("title")
	bookAuthors := r.FormValue("authors")

    var response = JsonResponse{}

    if bookID == "" || bookTitle == "" || bookAuthors == "" {
        response = JsonResponse{Type: "error", Message: "You are missing a valid ID for the Book or a for value title or author"}
    } else {
        db := setupDB()

        printMessage("Inserting the new book into DB")

        fmt.Println("Inserting new book with ID: " + bookID + " and name: " + bookTitle + " written by:" + bookAuthors )

        var lastInsertID int
    err := db.QueryRow("INSERT INTO books(bookID, bookTitle,bookAuthors) VALUES($1, $2, $3) returning id;", bookID, bookTitle, bookAuthors ).Scan(&lastInsertID)
    // check errors
    checkErr(err)

    response = JsonResponse{Type: "success", Message: "The book has been inserted successfully!"}
    }
    json.NewEncoder(w).Encode(response)
}

func PostBook(w http.ResponseWriter, r *http.Request) {
    bookID := r.FormValue("id")
    bookTitle := r.FormValue("title")
	bookAuthors := r.FormValue("authors")

    var response = JsonResponse{}

    if bookID == "" || bookTitle == "" || bookAuthors == "" {
        response = JsonResponse{Type: "error", Message: "You are missing a valid ID for the Book or a for value title or author"}
    } else {
        db := setupDB()
        printMessage("Updating the book value into DB")
        fmt.Println("Updating book with ID: " + bookID + " and title: " + bookTitle + " written by:" + bookAuthors )
        var lastInsertID int
		err := db.QueryRow("UPDATE INTO books(bookID, bookTitle,bookAuthors) VALUES($1, $2, $3) returning id;", bookID, bookTitle, bookAuthors ).Scan(&lastInsertID)
		// check errors
		checkErr(err)
		response = JsonResponse{Type: "success", Message: "The book has been inserted successfully!"}
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
        _, err := db.Exec("DELETE FROM books where bookID = $1", bookID)
        // check errors
        checkErr(err)
        response = JsonResponse{Type: "success", Message: "The book has been deleted successfully!"}
    }
    json.NewEncoder(w).Encode(response)
}



/*

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)
*/