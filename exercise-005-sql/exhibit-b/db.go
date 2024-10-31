package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func PanicOn(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := sql.Open("postgres", "dbname=test sslmode=disable")
	PanicOn(err)
	defer db.Close()

	// Start by clearing any previosu entries in the database
	_, err = db.Exec("TRUNCATE TABLE people RESTART IDENTITY CASCADE")
	PanicOn(err)

	_, err = db.Exec("INSERT INTO people(name, ssn) VALUES ($1, $2)", "Bruce Leroy", 111223333)
	PanicOn(err)

	_, err = db.Exec("INSERT INTO people(name, ssn) VALUES ($1, $2)", "Sho 'Nuff", 444556666)
	PanicOn(err)

	//_, err = db.Exec("INSERT INTO people(name, ssn) VALUES ($1, $2)", "Wilby Deleted", 777889999)
	//PanicOn(err)

	// Instead of the previous INSERT query, you can do this and get the id value back without a seperate query
	var id int
	err = db.QueryRow("INSERT INTO people(name, ssn) VALUES ($1, $2) RETURNING person_id", "Wilby Deleted", 777889999).Scan(&id)
	PanicOn(err)
	fmt.Println("Inserted person Wilby, and his ID is:", id)

	fmt.Println("The database now looks like this:")
	printDatabase(db)

	// delete a row
	_, err = db.Exec("DELETE FROM people WHERE name = ($1)", "Wilby Deleted")
	PanicOn(err)

	fmt.Println("After deleting a row, the database now looks like this:")
	printDatabase(db)

	// update an existing row
	_, err = db.Exec("UPDATE people SET ssn = $1 WHERE name = $2", 123456789, "Bruce Leroy")
	PanicOn(err)

	fmt.Println("After updating a row, the database now looks like this:")
	printDatabase(db)

}

func printDatabase(db *sql.DB) {
	rows, err := db.Query("SELECT person_id, name, ssn FROM people")
	PanicOn(err)

	for rows.Next() {
		var id int
		var name string
		var ssn int
		err := rows.Scan(&id, &name, &ssn)
		PanicOn(err)

		fmt.Printf("Person %5d %-15s %9d\n", id, name, ssn)
	}
	fmt.Println("")
}
