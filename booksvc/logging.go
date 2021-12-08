package booksvc

import (
	"errors"
	"strings"
	"context"
    "fmt"
    "log"
)

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