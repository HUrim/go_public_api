package main

import (
	"database/sql"
)

func createTable(db *sql.DB) {
	// CreateTable
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS Name (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR(50), country_id CHAR(2), probability FLOAT)")
	if err != nil {
		panic(err)
	}
}

func getPersonsQuery(db *sql.DB) ([]Person, error) {
	// getPersons
	rows, err := db.Query("SELECT * FROM Name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []Person

	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.Id, &person.Name, &person.Country_id,
			&person.Probability); err != nil {
			return persons, nil
		}
		persons = append(persons, person)
	}
	if err = rows.Err(); err != nil {
		return persons, err
	}
	return persons, nil
}

func createPerson(db *sql.DB, person Person) {
	// Create
	_, err := db.Exec("INSERT INTO Name (name, country_id, probability) VALUES (?, ?, ?)", person.Name, person.Country_id, person.Probability)
	if err != nil {
		panic(err)
	}
}

func deletePerson(db *sql.DB, name string) (count int) {
	// Delete
	rows, er := db.Query("SELECT COUNT(*) as count FROM  Name where name=?", name)
	if er != nil {
		panic(er)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
	}
	if count != 0 {
		_, err := db.Exec("DELETE FROM Name WHERE name = ?", name)
		if err != nil {
			panic(err)
		}
	}
	return count
}
