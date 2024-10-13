// db_config.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

// InitDB initializes and returns a singleton MySQL database connection
func InitDB() (*sql.DB, error) {
	var err error

	// Ensure this block is executed only once
	once.Do(func() {
		// Load environment variables for MySQL
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USERNAME")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_DATABASE")

		// Construct MySQL connection string
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

		// Open the MySQL connection
		dbInstance, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Println("Error connecting to database: ", err)
			return
		}

		// Ping the database to ensure connection is successful
		if err = dbInstance.Ping(); err != nil {
			log.Println("Error pinging database: ", err)
			return
		}

		log.Println("Database connection successfully established")
	})

	if err != nil {
		return nil, err
	}

	return dbInstance, nil
}
