package DB

import (
	"database/sql"
	"fmt"

	"github.com/anthonytb/account-auth-example/Utils"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	// Leaving these values in but in a prod app you should utilize env varibles
	host := Utils.GoDotEnvVariable("DB_HOST")
	user := Utils.GoDotEnvVariable("DB_USER")
	password := Utils.GoDotEnvVariable("DB_PASSWORD")
	db_name := Utils.GoDotEnvVariable("DB_NAME")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s", user, password, host, db_name)

	db, dbErr := sql.Open("postgres", connectionString)

	if dbErr != nil {
		return nil, dbErr
	}

	return db, nil
}
