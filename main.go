package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/anthonytb/account-auth-example/Services"
	"github.com/go-playground/validator"

	"golang.org/x/crypto/bcrypt"
)

// Route for users attempting to create an account
// Request Example: {Name: "Jane Doe", Email: "example@gmail.com", Password: "password123"}
func signup(w http.ResponseWriter, r *http.Request) {

	// Creates initial varible w/ User struct
	var user Services.User

	// Checks req body for JSON and parses it and applies it to user
	if jsonErr := json.NewDecoder(r.Body).Decode(&user); jsonErr != nil {
		// If no json body sends bad request
		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
		return
	}

	// Checks if user was updated w/ valid values
	if err := validator.New().Struct(user); err != nil {
		// Sends back json error message if check fails
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"message": "Malformed account creation request"})
		return
	}

	// Lowercases users email before being stored in DB to lowercase for more efficent searching later on
	user.Email = strings.ToLower(user.Email)

	// Takes user submitted password from the req body and returns the hashed password or an error
	hashedPass, passErr := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	// If hashing the password encountered a error then sends back json error message
	if passErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"message": "Issue creating account please try again. If issue persist contact support."})
	}

	// Assigns the newly hashed password as the user's password
	user.Password = string(hashedPass)

	// Submits users data from the body req to db to be stored
	if creationErr := Services.CreateAccount(user); creationErr != nil {

		// If function encountered a error then returns a json error message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)

		// Checks if the creation was caused buy duplicate unique values then if so then returns a more specific error message informing them about the duplicate key
		if creationErr.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			json.NewEncoder(w).Encode(map[string]string{"message": "Email is already assocaited with an account"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Error creating account"})
		fmt.Println("error", creationErr.Error())
		return
	}

	// Sends back 200 and success json message informing user their account has been successfully created
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully created account"})
}

type Credentials struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

// Route for users attempting to log into existing account
// Request Example: {Email: "example@gmail.com", Password: "password123"}
func login(w http.ResponseWriter, r *http.Request) {

	// Creates varible w/ the credentials struct
	var credentials Credentials

	// Checks req body for json body and parses it then assigns values to credentials
	jsonErr := json.NewDecoder(r.Body).Decode(&credentials)

	// If no json was in the req body then returns bad req
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
	}

	// Checks that credentials has valid values
	if err := validator.New().Struct(credentials); err != nil {
		// Sends error json message informing user of the bad request values
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"message": "Malformed login request"})
		return
	}

	// Lowercases users email before being comparing in DB
	credentials.Email = strings.ToLower(credentials.Email)

	// Checks DB for users that have the email that the req body has if so returns DB user
	user := Services.GetUserByUniqueConstraint("email", credentials.Email)

	// If no user was found in the DB then returns json error message informing abount the invalid data
	if user == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode([]byte(`{message: "Invalid Email/Password"}`))
		return
	}

	// Compares the req body password to the password from the returned DB user
	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))

	// If comparing the password fails then returns json error message informing user of invalid data
	if compareErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid Email/Password"})
		return
	}

	// Creates JWT for the user
	jwt := Services.CreateJWT(*user)

	// If creating the JWT encountered a error then returns json error message informing user of the issue and gives helpful solution
	if jwt == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error creating JWT please try again if issue persist. Contact support."})
		return
	}

	// Sends 200 request w/ json message containing the users authToken for them to utilize
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully logged in.", "authToken": *jwt})
}

func handler() {
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func main() {
	handler()
}
