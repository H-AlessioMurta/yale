package booksvc

import (
    "fmt"
)

type JsonResponse struct {
    Type    string `json:"type"`
    Data    []Book `json:"data"`
    Message string `json:"message"`
}


func printMessage(message string) {
    fmt.Println("")
    fmt.Println(message)
    fmt.Println("")
}


func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}