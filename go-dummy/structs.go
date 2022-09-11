package main

type Person struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Country_id  string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type PersonName struct {
	Name string
}

type Country struct {
	Country_id  string
	Probability float32
}

type ResponseObject struct {
	Name    string
	Country []Country
}
