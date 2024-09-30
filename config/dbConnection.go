package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to open the database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	fmt.Println("Successfully connected to database")
}
