package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	DB *sql.DB
}

var wg sync.WaitGroup

func insert(conn *sql.DB, st, en int, user []User) {

	defer wg.Done()
	for idx := st; idx <= en; idx++ {

		query := `
		INSERT INTO users (id, first_name, last_name, email, gender, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6);
		`

		_, err := conn.Exec(query, user[idx].ID, user[idx].FirstName, user[idx].LastName, user[idx].Email, user[idx].Gender, user[idx].IPAddress)
		if err != nil {
			log.Println("Error inserting data: ", err)
			return
		}

		log.Println("Data inserted successfully ", idx)
	}
}

func ConnectToDB() *DB {
	// Connection string for PostgreSQL
	connStr := "user=user dbname=mydb sslmode=disable password=password host=localhost port=5432"

	// Open the database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
		return &DB{
			DB: nil,
		}
	}

	// Configure connection pooling
	db.SetMaxOpenConns(25)                 // Maximum number of open connections to the database
	db.SetMaxIdleConns(25)                 // Maximum number of idle connections in the pool
	db.SetConnMaxIdleTime(5 * time.Minute) // Maximum amount of time a connection may be idle

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
		return &DB{
			DB: nil,
		}
	}

	fmt.Println("Successfully connected to the database")
	return &DB{
		DB: db,
	}
}

type Database interface {
	Insert()
}

func (db *DB) Insert(user []User) {
	log.Println("Insert is called")

	start := time.Now()
	wg.Add(4)

	go insert(db.DB, 0, 2500, user)
	go insert(db.DB, 2501, 5000, user)
	go insert(db.DB, 5001, 7500, user)
	go insert(db.DB, 7501, 10000, user)

	wg.Wait()

	end := time.Since(start)

	log.Println("All Data inserted successfully")
	log.Println("Total time taken is", end)
}
