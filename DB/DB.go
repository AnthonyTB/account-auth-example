package DB

import (
	"database/sql"
	"fmt"

	"github.com/anthonytb/account-auth-example/Utils"
	_ "github.com/lib/pq"
)

// Function for initializing DB connection
// Returns DB<pointer> or error
func InitDB() (*sql.DB, error) {
	// DB info from .env
	host := Utils.GoDotEnvVariable("DB_HOST")
	user := Utils.GoDotEnvVariable("DB_USER")
	password := Utils.GoDotEnvVariable("DB_PASSWORD")
	db_name := Utils.GoDotEnvVariable("DB_NAME")

	// Generates connection string for a postgres DB
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s", user, password, host, db_name)

	// Attempts to connect to DB
	db, dbErr := sql.Open("postgres", connectionString)

	// If connecting to DB was unsuccessful returns error
	if dbErr != nil {
		return nil, dbErr
	}

	// If connecting was successful then returns DB instance
	return db, nil
}
