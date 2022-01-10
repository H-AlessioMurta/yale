package main

import (
    "os"
    "testing"   
    "log"
    "yale/booksvc/src"
    "net/http"
    "net/http/httptest"
    "bytes"
    "encoding/json"
)

var a booksvc.Server
var path string

func TestMain(m *testing.M) {
    a.Initialize()
    ensureTableExists()
    //fmt.Println("TestMain")
    code := m.Run()
    //clearTable()
    os.Exit(code)
}

func ensureTableExists() {
    if _, err := a.DB.Exec(tableCreationQuery); err != nil {
        log.Fatal(err)
    }
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS books_test
(
    id character varying,
    title character varying,
    authors character varying
    
)`

func TestNotEmptyTable(t *testing.T) {
    
    req, _ := http.NewRequest("GET", "/books", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    if body := response.Body.String(); body == "[]" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)
    return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}

func TestCreateBook(t *testing.T) {
    var jsonStr = []byte(`{"title":"test book internal test", "authors": "Ciro"}`)   
    req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonStr))    
    req.Header.Set("Content-Type", "application/json")
    response := executeRequest(req)
    checkResponseCode(t, http.StatusCreated, response.Code)
    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["title"] != "test book internal test" {
        t.Errorf("Expected book title to be 'test book'. Got '%v'", m["title"])
    }
    if m["authors"] != "Ciro" {
        t.Errorf("Expected book authors to be 'Ciro'. Got '%v'", m["authors"])
    }
    setEndpointLastID()
}

func setEndpointLastID(){
    var lastID string
    row := a.DB.QueryRow("Select id from books where title =$1","test book internal test")
    err := row.Scan(&lastID)
    if err != nil{
      panic(err)
    }
     path = "/books/" +lastID
}

func TestGetBook(t *testing.T) {
    req, _ := http.NewRequest("GET", path, nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateBook(t *testing.T) {

    req, _ := http.NewRequest("GET", path, nil)
    response := executeRequest(req)
    var originalBook map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &originalBook)
    var jsonStr = []byte(`{"title":"test book - updated title", "authors": "zimuel"}`)
    req, _ = http.NewRequest("PUT", path, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")
    response = executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)
    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["id"] != originalBook["id"] {
        t.Errorf("Expected the id to remain the same (%v). Got %v", originalBook["id"], m["id"])
    }

    if m["title"] == originalBook["title"] {
        t.Errorf("Expected the title to change from '%v' to '%v'. Got '%v'", originalBook["title"], m["title"], m["title"])
    }

    if m["authors"] == originalBook["authors"] {
        t.Errorf("Expected the authors to change from '%v' to '%v'. Got '%v'", originalBook["authors"], m["authors"], m["authors"])
    }
}

func TestDeleteBook(t *testing.T) {
    req, _ := http.NewRequest("GET", path, nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)
    req, _ = http.NewRequest("DELETE", path, nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)
    req, _ = http.NewRequest("GET", path, nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusNotFound, response.Code)
}