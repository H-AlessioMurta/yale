package booksvc

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
    
)


const (
    DB_USER     = "postgres"
    DB_PASSWORD = "12345678"
    DB_NAME     = "books"
)

// DB set up
func setupDB() *sql.DB {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
    db, err := sql.Open("postgres", dbinfo)
    checkErr(err)
    return db
}