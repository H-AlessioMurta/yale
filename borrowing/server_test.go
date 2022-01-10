package main

import (
    "os"
    "testing"
	router "yale/borrowing/router"   
     l "yale/borrowing/logger"
    "net/http"
   // "net/http/httptest"
   // "bytes"
   	//"encoding/json"
	"yale/borrowing/graph"
	"yale/borrowing/graph/model"
	mongoDB "yale/borrowing/db"
	"yale/borrowing/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/client"
	"github.com/mitchellh/mapstructure"
)
//Variables for testing:
//gqlgen's server for handling request in Graphql 
var srv = handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
//A new connection to our borrowing's mongo collection Borrows
var db = mongoDB.Connect()
//Storing the id of the new element when created, because uuid is unprevedible
var idB, idC, idBorr string

//Return map transform the return of a grapql's post in a map for golang, like json.Unmarshal work with Json responses
func returnMap( qm string)(map[string]map[string]string){
	m := make(map[string]map[string]string)
	graphClient:= client.New(srv)
	raw, err :=graphClient.RawPost(qm)
	l.CheckErr(err)
	err= mapstructure.Decode(raw.Data,&m)
	l.CheckErr(err)
	return m
}
// returning the new id book while testing the Query Get_ID_Book
func GetIdBook()string{
	query := `
	query tornaIDlibro{
		Get_ID_Book(title:"TestingGoPostBook",authors:"Mr Murta"){
		id
		title
		authors
		}
	}
  `
  m := returnMap(query)// exec query
  idB= m["Get_ID_Book"]["id"]//fetching
  return idB
}
// returning the new id Customer while testing the Query Get_ID_Customer
func GetIdCustomer()string{
	query := `
	query tornaIDC{
		Get_ID_Customer(name:"Alessio",surname:"Murta",nin:"Mr Murta"){
		id
		name
		surname
		nin
		}
	}
  `
  m := returnMap(query)//exec
  idC=m["Get_ID_Customer"]["id"]//fetching
  return idC
}

func TestMain(m *testing.M) {
	idB,idC, idBorr="","",""	// not neeeded, golang default unassign variables to their 0 status, "" for strings
	http.Handle("/query", srv) // routing http to grapql server
    code := m.Run() // this command will execute all Test func
    os.Exit(code) //Returning fail or success
}

func TestCreateBook(t *testing.T){
	mutation := `
	mutation PostBook{
		Post_Book(title:"TestingGoPostBook",authors:"Mr Murta"){
		  title
		  authors
		}
	  }
	`
	m := returnMap(mutation)
	if  m["Post_Book"]["authors"] != "Mr Murta"{
		t.Errorf("Expected book's authors to be 'Mr Murta'. Got '%v'", m["Post_Book"]["authors"])
	}
	if m["Post_Book"]["title"] != "TestingGoPostBook" {
		t.Errorf("Expected book's title to be 'TestingGoPostBook'. Got '%v'", m["Post_Book"]["title"])
	}
}

func TestCreateCustomer(t *testing.T){
	mutation := `
	mutation PostCliente{
		Post_Customer(name:"Alessio",surname:"Murta",nin:"Mr Murta"){
		  name
		  surname
		  nin
		}
	  }
	`
	m := returnMap(mutation)
	if  m["Post_Customer"]["name"] != "Alessio"{
		t.Errorf("Expected name to be 'Mr Murta'. Got '%v'", m["Post_Customer"]["name"])
	}
	if m["Post_Customer"]["surname"] != "Murta" {
		t.Errorf("Expected surname to be 'Murta'. Got '%v'", m["Post_Customer"]["surname"])
	}
	if m["Post_Customer"]["nin"] != "Mr Murta" {
		t.Errorf("Expected nin to be 'Mr Murta'. Got '%v'", m["Post_Customer"]["nin"])
	}
}

func TestBorrow_create(t *testing.T){
	GetIdBook()
	GetIdCustomer()
	mutation :=`
	mutation TestBorrow_create{
		Borrow_create(data:{idBook:"`+idB+`",idCustomer:"`+idC+`"}){
			idBorrowing
			idCustomer
			idBook
		}
	  }
	`
	m := returnMap(mutation)
	idBorr=m["Borrow_create"]["idBorrowing"]
	if m["Borrow_create"]["idBook"] != idB{
		t.Errorf("Expected different book id. Got '%v'", m["Borrow_create"]["idBook"])
	}
	if m["Borrow_create"]["idCustomer"] != idC{
		t.Errorf("Expected different book id. Got '%v'", m["Borrow_create"]["idCustomer"])
	}
}


func TestPutBook(t *testing.T){
	mutation :=`mutation MettoLibro{
		Put_Book(id: "`+idB+`",title: "TestingGoPutBook",
			  ,authors: "J.K. Rowling"){
		  id
		  title
		  authors
		}
	  }
	  `
	  returnMap(mutation)
	  b,err:=router.GetBook(idB)
	  l.CheckErr(err)
	  if b.Title != "TestingGoPutBook"{
		  t.Errorf("Expected titile to be 'TestingGoPutBook'. Got '%v'",b.Title)
	  }
	  if b.Authors != "J.K. Rowling"{
		  t.Errorf("Expected authors to be 'J.K. Rowling'. Got '%v'",b.Authors )
	  }
}

func TestPutCustomer(t *testing.T){
	mutation :=`mutation MettoCliente{
		Put_Customer(id: "`+idC+`",name: "Rio",
			  ,surname: "Murta", nin:"r10"){
		  id
		  name
		  surname
		  nin
		}
	  }
	  `
	returnMap(mutation)
	c,err:=router.GetCustomer(idC)
	l.CheckErr(err)
	if c.Name != "Rio"{
		t.Errorf("Expected name to be 'Rio'. Got '%v'",c.Name )
	}
	if c.Surname != "Murta"{
		t.Errorf("Expected surname to be 'Murta'. Got '%v'",c.Surname )
	}
	if c.Nin != "r10"{
		t.Errorf("Expected nin to be 'Murta'. Got '%v'",c.Nin)
	}
}

func TestBorrowsnotreturned(t *testing.T){
	//Because there aren't easy ways to decode an array of pointers using grapql's response, we intercept our core logic calling the fetcher function on db
	borrows := db.Borrowsnotreturned()
	var bs *model.Borrowed
	for i, _ := range borrows{
		err:= mapstructure.Decode(borrows[i],&bs)
		l.CheckErr(err)
		if bs.IDBorrowing == idBorr && bs.Returned != false{
			t.Errorf("Expected  id.returned to be false. Got '%v'", bs.Returned)
		}
	}
}

func TestBorrowforCustomer(t *testing.T){
	//Because there aren't easy ways to decode an array of pointers using grapql's response, we intercept our core logic calling the fetcher function on db
	borrows := db.Borrowsnotreturned()
	var bs *model.Borrowed
	for i, _ := range borrows{
		err:= mapstructure.Decode(borrows[i],&bs)
		l.CheckErr(err)
		if bs.IDCustomer != idC{
			t.Errorf("Expected  %v.returned to be false. Got '%v'",idC, bs.IDCustomer)
		}
	}
}

func TestBorrowforBook(t *testing.T){
	//Because there aren't easy ways to decode an array of pointers using grapql's response, we intercept our core logic calling the fetcher function on db
	borrows := db.Borrowsnotreturned()
	var bs *model.Borrowed
	for i, _ := range borrows{
		err:= mapstructure.Decode(borrows[i],&bs)
		l.CheckErr(err)
		if bs.IDBook != idB{
			t.Errorf("Expected  %v.returned to be false. Got '%v'",idB, bs.IDBook)
		}
	}
}

func TestReturnedbookBorrow(t *testing.T){
	graphClient:= client.New(srv)
	mutation :=`mutation Borr{
		Returnedbook(id:"`+idBorr+`"){
		  returned
		}
	  }
	`
	graphClient.RawPost(mutation)//exec without storing the mutation returnedbook, cause this mutation will fetch the old value of returned
	m := make(map[string]map[string]bool)//because we need to map a bool value in the next rows we will emulate returnMap's logic but for this purpose only.
	//Call a query forour idborrowing
	query:=`
	query tuttiiborri{
		borrow(id:"`+idBorr+`"){
		  returned
		}
	}
	`
	raw, err :=graphClient.RawPost(query)
	l.CheckErr(err)
	err= mapstructure.Decode(raw.Data,&m)
	if m["borrow"]["returned"] != true{
		  t.Errorf("Expected true. Got '%v'", m["borrow"]["returned"])
	}
}

func TestDeleteBook(t *testing.T){
	graphClient:= client.New(srv)
	mutation :=
	`mutation cancello{
		Delete_Book(id: "`+idB+`")
	  }
	`
	graphClient.RawPost(mutation)
	_,err:=router.GetBook(idB)
	if err != nil{
		t.Errorf("No deleted item %v", idB)
	}
}

func TestDeleteCustomer(t *testing.T){
	graphClient:= client.New(srv)
	mutation :=
	`mutation cancello{
		Delete_Customer(id: "`+idC+`")
	  }
	`
	graphClient.RawPost(mutation)
	_,err:=router.GetCustomer(idC)
	if err != nil{
		t.Errorf("No deleted item %v", idC)
	}
}

