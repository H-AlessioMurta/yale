package main

import (
	//"~/GitHubRepo/yale/booksvc"
	"github.com/H-AlessioMurta/yale/booksvc"
)

func main(){
	// serve the app
	fmt.Println("Server at 8085")
	log.Fatal(http.ListenAndServe(":8085", router))
}