package customersvc

import (
    "log"
	"fmt"
	"database/sql"
    "net/http"
    "errors"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

//Parameters for connecting on a postgreSql server hosted on a docker image
const (
    DB_HOST = "customer-postgres"
    DB_PORT = 5433
    DB_USER     = "postgres"
    DB_PASSWORD = "postgres"
    DB_NAME     = "customers"
)

// DB and Handler set up
func (csvc *Server) Initialize(){
    logInfo("Opening DB connection" )
    dbinfo := fmt.Sprintf("host=%s  port=%d user=%s password=%s dbname=%s sslmode=disable", DB_HOST,DB_PORT,DB_USER, DB_PASSWORD, DB_NAME)
    var err error
    csvc.DB, err = sql.Open("postgres", dbinfo)
    checkErr(err)
    csvc.Router = mux.NewRouter()
    logInfo("Opening service router" )
    csvc.initializeRoutes()
}

func (csvc *Server) Run(addr string) {
    logInfo("Starting customersvc on port:" + addr )
    log.Fatal(http.ListenAndServe(addr, csvc.Router))
    
}

func (c *Customer) getCustomer(db *sql.DB) error {
    return db.QueryRow("SELECT name, surname, nin FROM customers WHERE id=$1",
        c.ID).Scan(&c.Name, &c.Surname, &c.Nin)
}

func (c *Customer) putCustomer(db *sql.DB) error {
    if c.Name != "" && c.Surname !="" && c.Nin != "" {
        _, err := db.Exec("UPDATE customers SET name=$1, surname=$2, nin =$3 WHERE id=$4",c.Name, c.Surname, c.Nin, c.ID)
        return err
    }
    if c.Name == "" && c.Surname !="" && c.Nin == ""{
        _, err := db.Exec("UPDATE customers SET surname=$1 WHERE id=$2", c.Surname, c.ID)
        return err
    }
    if c.Name != "" && c.Surname =="" && c.Nin == ""{
        _, err := db.Exec("UPDATE customers SET name =$1 WHERE id=$2", c.Name, c.ID)
        return err
    }
    if c.Name != "" && c.Surname !="" && c.Nin == ""{
        _, err := db.Exec("UPDATE customers SET nin =$1 WHERE id=$2", c.Nin, c.ID)
        return err
    }

    return errors.New("Unexpected params for the update")
}

func (c *Customer) deleteCustomer(db *sql.DB) error {
    _, err := db.Exec("DELETE FROM customers WHERE id=$1", c.ID)
    return err
}

func (c *Customer) postCustomer(db *sql.DB) error {
    var i int
    err := db.QueryRow("SELECT COUNT(*) FROM cusomers WHERE name =$1 AND surname = $2 and nin =3",c.Name, c.Surname, c.Nin).Scan(&i)
    checkErr(err)
    if i < 1{      
        _, err := db.Exec("INSERT INTO customers(name,surname,nin) VALUES($1, $2,$3)", c.Name, c.Surname, c.Nin)
        checkErr(err)
        return nil
    }
    return errors.New("Already insert a Customer  with this personal information")
}

func getCustomers(db *sql.DB) ([]Customer, error) {
	//rows, err := db.Query("SELECT id, name,  surname FROM customers LIMIT $1 OFFSET $2",count, start)
    rows, err := db.Query("SELECT * FROM customers")
	checkErr(err)
    defer rows.Close()
    customers := []Customer{}
    for rows.Next() {
        var c Customer
        var nullID, nullName, nullSurname, nullNin sql.NullString 
        err := rows.Scan(&nullID, &nullName, &nullSurname, &nullNin)
        checkErr(err)
        if nullID.Valid && nullName.Valid && nullSurname.Valid && nullNin.Valid{
            c.ID = nullID.String
            c.Name = nullName.String
            c.Surname = nullSurname.String
            c.Nin = nullNin.String
            customers = append(customers, c)
        }
    }
    return customers, nil
  }
