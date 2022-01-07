package customersvc


type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
	Nin string `json:"nin,omitempty"`//national insurance number
}



