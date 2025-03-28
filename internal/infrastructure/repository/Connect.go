package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func CheckDB() (*sql.DB, error) {
	dbHost := "0.0.0.0"
	dbUser := "hacker"
	dbPassword := "password"
	dbName := "1337bo4rd"
	dbPort := "5432"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
