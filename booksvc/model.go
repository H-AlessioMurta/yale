package booksvc


type Book struct {
	ID        string    `json:"id"`
	Title      string    `json:"title,omitempty"`
	Authors string `json:"authors,omitempty"`
}



