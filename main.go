package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/anthonytb/account-auth-example/Services"
	"gopkg.in/validator.v2"

	"golang.org/x/crypto/bcrypt"
)

func signup(w http.ResponseWriter, r *http.Request) {

	var user Services.User

	if jsonErr := json.NewDecoder(r.Body).Decode(&user); jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.Validate(user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"message": "Malformed account creation request"})
		return
	}

	hashedPass, passErr := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	if passErr != nil {
		http.Error(w, passErr.Error(), http.StatusBadRequest)
	}

	user.Password = string(hashedPass)

	if creationErr := Services.CreateAccount(user); creationErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error creating account"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully created account"})
}

type Credentials struct {
	Email    string
	Password string
}

func login(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials

	jsonErr := json.NewDecoder(r.Body).Decode(&credentials)

	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
	}

	fmt.Println("email", credentials.Email)

	user := Services.GetUser("email", credentials.Email)

	fmt.Println("user", user)

	if user == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode([]byte(`{message: "Invalid Email/Password"}`))
		return
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))

	if compareErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid Email/Password"})
		return
	}

	jwt := Services.CreateJWT(*user)

	if jwt == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error creating JWT please try again if issue persist. Contact support."})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully logged in.", "authToken": *jwt})

	// Document
}

func handler() {
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func main() {
	handler()
}
