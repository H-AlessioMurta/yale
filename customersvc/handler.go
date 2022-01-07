package customersvc


import (
    //"io/ioutil"
    "encoding/json"
    "net/http"
    "database/sql"
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)


type Server struct {
	Router *mux.Router
	DB     *sql.DB
}


func (s *Server) initializeRoutes() {
    logInfo("Initializing http routes")
    s.Router.HandleFunc("/customers", s.getCustomers).Methods("GET")
    s.Router.HandleFunc("/customers", s.postCustomer).Methods("POST")
    s.Router.HandleFunc("/customers/{id}", s.getCustomer).Methods("GET")
    s.Router.HandleFunc("/customers/{id}", s.putCustomer).Methods("PUT")
    s.Router.HandleFunc("/customers/{id}", s.deleteCustomer).Methods("DELETE")
    logInfo("CRUD's path are ready ")
}




// CustomerService provides operations on DB.
type CustomerService interface {
	postCustomer(w http.ResponseWriter, r *http.Request) 
	getCustomer(w http.ResponseWriter, r *http.Request) 
	putCustomer(w http.ResponseWriter, r *http.Request) 
	deleteCustomer(w http.ResponseWriter, r *http.Request)
	getCustomers(w http.ResponseWriter, r *http.Request)
}

type customerService struct{}

func (s *Server) getCustomer(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    vars := mux.Vars(r)
    id:= vars["id"]
    b := Customer{ID: id}
    if err := b.getCustomer(s.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "Customer not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }
    respondWithJSON(w, http.StatusOK, b)
}

func (s *Server) getCustomers(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    customers, err := getCustomers(s.DB)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusOK, customers)
}

func (s *Server) postCustomer(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    var b Customer
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&b); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()
    if err := b.postCustomer(s.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusCreated, b)
}

func (s *Server) putCustomer(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    vars := mux.Vars(r)
    id := vars["id"]
    if id == "" {
        respondWithError(w, http.StatusBadRequest, "Invalid Customer ID")
        return
    }

    var b Customer
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&b); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
        return
    }
    defer r.Body.Close()
    b.ID= id
    if err := b.putCustomer(s.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, b)
}

func (s *Server) deleteCustomer(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    vars := mux.Vars(r)
    id := vars["id"]
    if id == "" {
        respondWithError(w, http.StatusBadRequest, "Invalid Customer ID")
        return
    }
    b := Customer{ID: id}
    if err := b.deleteCustomer(s.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

