package Services

import (
	"fmt"

	"github.com/anthonytb/account-auth-example/DB"
)

type User struct {
	id         int
	Name       string `validate:"required"`
	Email      string `validate:"required,email"`
	Password   string `validate:"required"`
	created_at string
}

// Function for inserting users into DB
// Params: user<User>
// Returns: nil or error
func CreateAccount(user User) error {
	// Initializes DB connection
	database, dbErr := DB.InitDB()

	// If DB connection encountered an error then returns error
	if dbErr != nil {
		return dbErr
	}

	// Inserts user into DB
	_, insertErr := database.Exec("INSERT INTO users(Name,Email,Password) VALUES($1, $2, $3)", user.Name, user.Email, user.Password)

	// If insertion into DB encountered an error then returns error
	if insertErr != nil {
		return insertErr
	}

	// Returns nil if encountered no errors
	return nil
}

// Function for finsing a user by a unique value
// Params: constraint<UNIQUE string>, constraintVal<string>
// Returns: user pointer or nil
func GetUserByUniqueConstraint(constraint string, constraintVal string) *User {
	var user User
	var userPointer *User = &user

	// Initializes DB connection
	database, dbErr := DB.InitDB()

	// If DB connection encountered an error then returns nil
	if dbErr != nil {
		return nil
	}

	// Checks DB for user with the following constraintVal and if user is found then assigns DB user values to user varible
	queryErr := database.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE %s = $1;", constraint), constraintVal).Scan(&user.id, &user.Name, &user.Email, &user.Password, &user.created_at)

	// If no users found or encountered error returns nil
	if queryErr != nil {
		return nil
	}

	// Returns the pointer that way you can return either the user data or nil
	return userPointer
}
