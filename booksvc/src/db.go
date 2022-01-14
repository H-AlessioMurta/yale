package booksvc

import (
    "log"
	"fmt"
	"database/sql"
    "net/http"
    "errors"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

//Parameters for connecting on a postgreSql server hosted on a docker image
const (
    DB_HOST = "book-postgres"
    DB_PORT = 5432
    DB_USER     = "postgres"
    DB_PASSWORD = "postgres"
    DB_NAME     = "books"
)

// DB and Handler set up
func (bsvc *Server) Initialize(){
    logInfo("Opening DB connection" )
    dbinfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",DB_USER, DB_PASSWORD,DB_HOST,DB_PORT, DB_NAME)
    var err error
    bsvc.DB, err = sql.Open("postgres", dbinfo)
    checkErr(err)
    bsvc.Router = mux.NewRouter()
    logInfo("Opening service router" )
    bsvc.initializeRoutes()
}

func (bsvc *Server) Run(addr string) {
    logInfo("Starting booksvc on port:" + addr )
    log.Fatal(http.ListenAndServe(addr, bsvc.Router))
    
}

func (b *Book) getBook(db *sql.DB) error {
    return db.QueryRow("SELECT title, authors FROM books WHERE id=$1",
        b.ID).Scan(&b.Title, &b.Authors)
}

func (b *Book) putBook(db *sql.DB) error {
    if b.Title != "" && b.Authors !="" {
        _, err := db.Exec("UPDATE books SET title=$1, authors=$2 WHERE id=$3",b.Title, b.Authors, b.ID)
        return err
    }
    if b.Title == "" && b.Authors !="" {
        _, err := db.Exec("UPDATE books SET authors=$1 WHERE id=$2", b.Authors, b.ID)
        return err
    }
    if b.Title != "" && b.Authors =="" {
        _, err := db.Exec("UPDATE books SET title =$1 WHERE id=$2", b.Title, b.ID)
        return err
    }
    return errors.New("Unexpected params for the update")
}

func (b *Book) deleteBook(db *sql.DB) error {
    _, err := db.Exec("DELETE FROM books WHERE id=$1", b.ID)
    return err
}

func (b *Book) postBook(db *sql.DB) error {
    var i int
    err := db.QueryRow("SELECT COUNT(*) FROM books WHERE title =$1 AND authors = $2",b.Title, b.Authors).Scan(&i)
    checkErr(err)
    if i < 1{      
        _, err := db.Exec("INSERT INTO books(title,authors) VALUES($1, $2)", b.Title, b.Authors)
        checkErr(err)
        return nil
    }
    return errors.New("Already insert a title with this author")
}

func getBooks(db *sql.DB) ([]Book, error) {
	//rows, err := db.Query("SELECT id, title,  authors FROM books LIMIT $1 OFFSET $2",count, start)
    rows, err := db.Query("SELECT * FROM books")
	checkErr(err)
    defer rows.Close()
    books := []Book{}
    for rows.Next() {
        var b Book
        var nullID, nullTitle, nullAuthors sql.NullString 
        err := rows.Scan(&nullID, &nullTitle, &nullAuthors)
        checkErr(err)
        if nullID.Valid && nullTitle.Valid && nullAuthors.Valid{
            b.ID = nullID.String
            b.Title = nullTitle.String
            b.Authors = nullAuthors.String
            books = append(books, b)
        }
    }
    return books, nil
  }