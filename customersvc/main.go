package main

import (
	"yale/customersvc/src"
	"log"
	"os"
	//log "github.com/sirupsen/logrus"
)

/*colorReset := "\033[0m"

colorRed := "\033[31m"
colorGreen := "\033[32m"
colorYellow := "\033[33m"
colorBlue := "\033[34m"
colorPurple := "\033[35m"
colorCyan := "\033[36m"
colorWhite := "\033[37m"*/

func main(){
	//var log = logrus.New()
	log.SetPrefix("\033[36mcustomersvc\033[0m: ")
	log.SetOutput(os.Stdout)
	bsvc := customersvc.Server{}
	bsvc.Initialize()
	bsvc.Run(":8014")
}