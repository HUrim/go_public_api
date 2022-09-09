package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type Person struct {
	// id          int     `json:"id"`
	name        string  `json:"name"`
	country_id  string  `json:"country_id"`
	probability float32 `json:"probability"`
}

func dbConnection() *sql.DB {
	db, err := sql.Open("mysql", "user:Urim.123@tcp(127.0.0.1:3306)/go_exercise_db")
	if err != nil {
		panic(err)
	} else if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}

// func getPersons(c echo.Context) error {
// 	return c.String(http.StatusOK, fmt.Sprintf("person name: %s\nperson country_id: %s\nperson probability: %s", person.name, person.country_id, person.probability))
// }

func getPersonsQuery(db *sql.DB) ([]Person, error) {
	// Read
	rows, err := db.Query("SELECT * FROM Name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []Person

	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.name, &person.country_id,
			&person.probability); err != nil {
			return persons, nil
		}
		persons = append(persons, person)
	}
	if err = rows.Err(); err != nil {
		return persons, err
	}
	return persons, nil
}

func getPersons(c echo.Context) error {

	var persons []Person
	db := dbConnection()

	persons, err := getPersonsQuery(db)
	fmt.Println(persons)
	if err != nil {
		panic(err)
	}

	var jsonResponse string

	// for i := 0; i < len(persons); i++ {
	// 	append(jsonResponse, [string]string{
	// 		"id":          strconv.Itoa(persons[0].id),
	// 		"name":        persons[0].name,
	// 		"country_id":  persons[0].country_id,
	// 		"probability": fmt.Sprintf("%f", persons[0].probability),
	// 	})
	// }

	return c.JSON(http.StatusOK, string(jsonResponse)) ///
}

type PersonName struct {
	Name string
}

func createPerson(db *sql.DB, person Person) {
	// Create
	_, err := db.Exec("INSERT INTO Name (name, country_id, probability) VALUES (?, ?, ?)", person.name, person.country_id, person.probability)
	if err != nil {
		panic(err)
	}
}

func deletePerson(db *sql.DB, name string) {
	// Delete
	_, err := db.Exec("DELETE FROM Name WHERE name = ?", name)
	if err != nil {
		panic(err)
	}
}

func addPerson(c echo.Context) error {

	personName := PersonName{}
	defer c.Request().Body.Close()

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading: ", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(body, &personName)
	if err != nil {
		log.Printf("Failed unmarshalling: ", err)
		return c.String(http.StatusInternalServerError, "")
	}

	fmt.Println(personName.Name)

	responseData := responseFromApi(personName.Name)

	// stringResponseData := string(responseData)
	// splitResponseData(string(responseData))
	fmt.Println(string(responseData))

	var responseObject ResponseObject
	json.Unmarshal(responseData, &responseObject)
	if responseObject.Country == nil {
		return c.String(http.StatusBadRequest, "")
	}

	name, country_id, probability := highestPercentage(responseObject)

	db := dbConnection()

	// Create
	createPerson(db, Person{name: name, country_id: country_id, probability: probability})

	return c.String(http.StatusOK, "added fshije")
}

func deleteExistingPerson(c echo.Context) error {

	personName := PersonName{}
	defer c.Request().Body.Close()

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading: ", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(body, &personName)
	if err != nil {
		log.Printf("Failed unmarshalling: ", err)
		return c.String(http.StatusInternalServerError, "")
	}

	db := dbConnection()

	deletePerson(db, personName.Name)

	fmt.Println(personName.Name)

	return c.String(http.StatusOK, "fdfr")
}

func responseFromApi(nameForApi string) []byte {

	response, err := http.Get("https://api.nationalize.io/?name=" + nameForApi)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseData
}

type Country struct {
	Country_id  string
	Probability float32
}

type ResponseObject struct {
	Name    string
	Country []Country
}

func highestPercentage(responseObject ResponseObject) (string, string, float32) {
	maxN := responseObject.Name
	maxP := responseObject.Country[0].Probability
	maxC := responseObject.Country[0].Country_id
	for i := 1; i < len(responseObject.Country); i++ {
		if responseObject.Country[i].Probability > maxP {
			maxP = responseObject.Country[i].Probability
			maxC = responseObject.Country[i].Country_id
		}
	}

	return maxN, maxC, maxP
}

func main() {

	e := echo.New()

	e.GET("/api/names/", getPersons) // get
	// e.GET("https://api.nationalize.io/:name", readFromAPI)

	e.POST("/api/names/:", addPerson) // post

	e.DELETE("/api/name/:name", deleteExistingPerson) //delete

	e.Start(":9191")
}
