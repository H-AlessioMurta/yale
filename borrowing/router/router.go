/*
* In this package, we gonna use all crud operation on booksvc and customer svc, as SAGA pattern
*/

package router

import (
	"yale/borrowing/graph/model"
	l "yale/borrowing/logger"
	"net/http"
    "encoding/json"
	"time"
	"fmt"
	"bytes"
	"errors"
	"io/ioutil"
)

const (
	BOOKAPI = "http://10.109.224.28:30088/books"
	CUSTOMERAPI = "http://10.105.151.100:32647/customers"
)


func parsingErr(b http.Response)string{
	body, err := ioutil.ReadAll(b.Body)
	l.CheckErr(err)
	jsonStr := string(body)
	return jsonStr
}

//Setting a new Client with timeout in getting request
var myClient =&http.Client{Timeout:10*time.Second}

func GetBooks()([]*model.Book, error) {
	fmt.Printf("")
	response, err := myClient.Get(BOOKAPI)
	l.CheckErr(err)
	defer response.Body.Close()
	var books []*model.Book
	err = json.NewDecoder(response.Body).Decode(&books)
	l.CheckErr(err)
	for i, _:= range books{
		l.LogResponseBook(books[i])
	}
	return books, err
}

func GetBook(id string)(*model.Book, error) {
	response, err := http.Get(BOOKAPI+"/"+id)
	l.CheckErr(err)
	defer response.Body.Close()
	var book *model.Book
	err = json.NewDecoder(response.Body).Decode(&book)
	l.CheckErr(err)
	l.LogResponseBook(book)
	return book, err
}

func GetCustomers()([]*model.Customer, error) {
	response, err := myClient.Get(CUSTOMERAPI)
	l.CheckErr(err)
	defer response.Body.Close()
	var customers []*model.Customer
	err = json.NewDecoder(response.Body).Decode(&customers)
	l.CheckErr(err)
	for i, _:= range customers{
		l.LogResponseCustomer(customers[i])
	}
	return customers, err
}

func GetCustomer(id string)(*model.Customer, error) {
	response, err := http.Get(CUSTOMERAPI+"/"+id)
	l.CheckErr(err)
	defer response.Body.Close()
	var customer *model.Customer
	err = json.NewDecoder(response.Body).Decode(&customer)
	l.CheckErr(err)
	l.LogResponseCustomer(customer)
	return customer, err
}

func CheckBook(id string) bool{
	check,_:= GetBook(id)
	if *check != (model.Book{}){
		return true
 	}else{
		 return false
	 } 
}

func CheckCustomer(id string) bool{
	check,_ := GetCustomer(id)
	if *check != (model.Customer{}){
		return true
 	}else{
		 return false
	 }  
}

func BookRequestHandler(payload map[string]string,method string, id string)(*model.Book,error){
	path := BOOKAPI
	if id != ""{
		path = path+"/"+id
	}
	body,err := json.Marshal(payload)
	l.CheckErr(err)
	request, err := http.NewRequest(method,path, bytes.NewBuffer(body))
	l.CheckErr(err)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	l.LogRequest(request)
	response, err := myClient.Do(request)
	l.CheckErr(err)
	defer request.Body.Close()
	if response.StatusCode == http.StatusCreated{
		var book *model.Book
		err = json.NewDecoder(response.Body).Decode(&book)
		l.CheckErr(err)
		l.LogResponseBook(book)
		return book, err 
	}else{
		 return &model.Book{}, errors.New(parsingErr(*response))
	}
}

func CustomerRequestHandler(payload map[string]string,method string, id string)(*model.Customer,error){
	path := CUSTOMERAPI
	if id != ""{
		path = path+"/"+id
	}
	body,err := json.Marshal(payload)
	l.CheckErr(err)
	request, err := http.NewRequest(method,path, bytes.NewBuffer(body))
	l.CheckErr(err)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	l.LogRequest(request)
	response, err := myClient.Do(request)
	l.CheckErr(err)
	defer request.Body.Close()
	if response.StatusCode == http.StatusCreated{
		var customer *model.Customer
		err = json.NewDecoder(response.Body).Decode(&customer)
		l.CheckErr(err)
		l.LogResponseCustomer(customer)
		return customer, err
	}else{
		 return &model.Customer{}, errors.New(parsingErr(*response))
	}
 }
 

func Eraser(apiUri string, id string)(*string,error){
	r :="Can't delete the resoruse with id:"+id
	apipath := apiUri+"/"+id
	request, err := http.NewRequest(http.MethodDelete,apipath,nil)
	l.CheckErr(err)
	l.LogRequest(request)
	response, err := myClient.Do(request)
	l.CheckErr(err)
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK{
		l.LogResponse("Delete complete")
		return &id,err
	}
	l.LogError(r)
	return &id,errors.New("Not found")
}

func GetIDBook(t string, a string)(*model.Book,error){
	books, err := GetBooks()
	l.CheckErr(err)
	for i, _:=range books{
		if books[i].Title == t && books[i].Authors == a{
			return books[i],nil
		} 
	}
	return nil, errors.New("No Book with this title and author")
}

func GetIDCustomer(n string, s string,nin string)(*model.Customer,error){
	cs, err := GetCustomers()
	l.CheckErr(err)
	for i, _:=range cs{
		if cs[i].Name == n && cs[i].Surname == s && cs[i].Nin == nin{
			return cs[i],nil
		} 
	}
	return nil, errors.New("No Customer with this name, surname, nin")
}