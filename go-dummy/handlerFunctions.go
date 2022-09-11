package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func getPersons(c echo.Context) error {

	var persons []Person
	db := dbConnection()
	//getPersons
	persons, err := getPersonsQuery(db)
	if err != nil {
		panic(err)
	}

	resp, _ := json.Marshal(persons)
	return c.JSON(http.StatusOK, string(resp)) ///
}

func addPerson(c echo.Context) error {

	personName := PersonName{}
	defer c.Request().Body.Close()

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(body, &personName)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}
	responseData := responseFromApi(personName.Name)

	var responseObject ResponseObject
	json.Unmarshal(responseData, &responseObject)
	if responseObject.Country == nil {
		return c.String(http.StatusBadRequest, "")
	}

	name, country_id, probability := highestPercentage(responseObject)
	if country_id == "" && probability == 0 {
		return c.String(http.StatusNotAcceptable, personName.Name+" country's not found")
	}

	db := dbConnection()
	//create
	createPerson(db, Person{Name: name, Country_id: country_id, Probability: probability})

	return c.String(http.StatusOK, personName.Name+" added succesfully")
}

func deleteExistingPerson(c echo.Context) error {

	personName := PersonName{}
	defer c.Request().Body.Close()

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(body, &personName)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	db := dbConnection()
	//delete
	count := deletePerson(db, personName.Name)
	if count == 0 {
		return c.String(http.StatusNotFound, personName.Name+" doesn't exist in the database")
	}

	return c.String(http.StatusOK, personName.Name+" deleted successfully")
}

func responseFromApi(nameForApi string) []byte {

	response, err := http.Get("https://api.nationalize.io/?name=" + nameForApi)

	if err != nil {
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseData
}

func highestPercentage(responseObject ResponseObject) (string, string, float32) {
	maxN := responseObject.Name
	var maxP float32
	var maxC string
	if len(responseObject.Country) != 0 {
		maxP = responseObject.Country[0].Probability
		maxC = responseObject.Country[0].Country_id
		for i := 1; i < len(responseObject.Country); i++ {
			if responseObject.Country[i].Probability > maxP {
				maxP = responseObject.Country[i].Probability
				maxC = responseObject.Country[i].Country_id
			}
		}
	} else {
		maxP = 0
		maxC = ""
	}

	return maxN, maxC, maxP
}
