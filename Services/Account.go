package Services

import (
	"fmt"
	"log"

	"github.com/anthonytb/account-auth-example/db"
)

type User struct {
	id         int
	Name       string
	Email      string
	Password   string
	created_at string
}

func CreateAccount(user User) {
	database, dbErr := db.InitDB()

	fmt.Println("User Creation:", user)

	if dbErr != nil {
		log.Fatal("Database connection failed")
	}

	_, insertErr := database.Exec("INSERT INTO users(Name,Email,Password) VALUES($1, $2, $3)", user.Name, user.Email, user.Password)

	if insertErr != nil {
		log.Fatal("Error inserting in database table", insertErr)
	}
}

func GetUser(constraint string, constraintVal string) *User {
	var user User
	var userPointer *User

	database, dbErr := db.InitDB()

	if dbErr != nil {
		log.Fatal("Database connection failed")
	}

	queryErr := database.QueryRow("SELECT * FROM users WHERE $1=$2", constraint, constraintVal).Scan(&user.id, &user.Name, &user.Email, &user.Password, &user.created_at)

	if queryErr != nil {
		user = User{}
	}

	fmt.Println("user found", user)

	return userPointer
}
