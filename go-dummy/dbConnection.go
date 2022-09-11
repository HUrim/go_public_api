package main

import "database/sql"

func dbConnection() *sql.DB {
	db, err := sql.Open("mysql", "user:Urim.123@tcp(127.0.0.1:3306)/go_exercise_db")
	if err != nil {
		panic(err)
	} else if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}
