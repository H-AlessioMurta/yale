package main

import (
	"yale/booksvc/src"
	"log"
	"os"
)

func main(){
	//var log = logrus.New()
	log.SetPrefix("\033[36mbooksvc\033[0m: ")
	log.SetOutput(os.Stdout)
	bsvc := booksvc.Server{}
	bsvc.Initialize()
	bsvc.Run(":8013")
}