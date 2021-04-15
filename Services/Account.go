package Services

import (
	"fmt"
	"log"

	"github.com/anthonytb/account-auth-example/DB"
)

type User struct {
	id         int
	Name       string `json:"Name" validate:"min=1"`
	Email      string `json:"Email" validate:"min=1"`
	Password   string `json:"Password" validate:"min=1"`
	created_at string
}

func CreateAccount(user User) error {
	database, dbErr := DB.InitDB()

	if dbErr != nil {
		return dbErr
	}

	_, insertErr := database.Exec("INSERT INTO users(Name,Email,Password) VALUES($1, $2, $3)", user.Name, user.Email, user.Password)

	if insertErr != nil {
		return insertErr
	}

	return nil
}

func GetUser(constraint string, constraintVal string) *User {
	var user User
	var userPointer *User = &user

	database, dbErr := DB.InitDB()

	if dbErr != nil {
		log.Fatal("Database connection failed")
	}

	queryErr := database.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE %s = $1;", constraint), constraintVal).Scan(&user.id, &user.Name, &user.Email, &user.Password, &user.created_at)

	if queryErr != nil {
		return nil
	}

	return userPointer
}
