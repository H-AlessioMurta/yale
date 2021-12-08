package booksvc

import (
	"errors"
	"strings"
	"context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

router := mux.NewRouter()

// Route handles & endpoints

// Get all movies
router.HandleFunc("/books/", GetBooks).Methods("GET")

// Get all movies
router.HandleFunc("/books/{id}", GetBooks).Methods("GET")


// Create a book
router.HandleFunc("/books/", PostBook).Methods("POST")

// Delete a specific book by the bookID
router.HandleFunc("/books/{id}", Deletebook).Methods("DELETE")

