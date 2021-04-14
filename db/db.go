package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	// Leaving these values in but in a prod app you should utilize env varibles
	host := "localhost:5432"
	user := "postgres"
	password := "password"
	db_name := "example"

	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s", user, password, host, db_name)

	fmt.Println("Connection:", connectionString)

	db, dbErr := sql.Open("postgres", connectionString)

	if dbErr != nil {
		fmt.Println("Connection Err:", dbErr)
		return nil, dbErr
	}

	return db, nil
}
