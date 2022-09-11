package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

func main() {

	e := echo.New()

	e.GET("/api/names/", getPersons) // get

	e.POST("/api/names/", addPerson) // post

	e.DELETE("/api/name/", deleteExistingPerson) //delete

	e.Start(":9191")
}
