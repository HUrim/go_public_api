package main

import "database/sql"

type Database struct {
	NAME          string
	PASSWORD      string
	DATABASE_NAME string
}

func dbConnection() *sql.DB {
	// change user password and database_name to your implementing
	dbConnection := Database{
		"user",
		"password",
		"db_name",
	}
	db, err := sql.Open("mysql", dbConnection.NAME+":"+dbConnection.PASSWORD+"@tcp(127.0.1:3306)/"+dbConnection.DATABASE_NAME)
	if err != nil {
		panic(err)
	} else if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}
