package main

import (
    "fmt"
    "log"
    "net/http"
	"github.com/gorilla/mux"
	"yale/booksvc"
)

func main(){
	router := mux.NewRouter()
	// Route handles & endpoints
	// Get all books
	router.HandleFunc("/books/", booksvc.GetBooks).Methods("GET")
	// Get a specifi books
	router.HandleFunc("/books/{id}", booksvc.GetBook).Methods("GET")
	// Create a book
	router.HandleFunc("/books/", booksvc.PostBook).Methods("POST")
	// Delete a specific book by the bookID
	router.HandleFunc("/books/{id}", booksvc.DeleteBook).Methods("DELETE")
	//Updating a book parameters by ID
	router.HandleFunc("/books/{id}", booksvc.PutBook).Methods("PUT")
	fmt.Println("Server at 8085")
	log.Fatal(http.ListenAndServe(":8085", router))
}