package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/anthonytb/account-auth-example/Services"

	"golang.org/x/crypto/bcrypt"
)

func signup(w http.ResponseWriter, r *http.Request) {

	var user Services.User

	jsonErr := json.NewDecoder(r.Body).Decode(&user)

	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
		return
	}

	hashedPass, passErr := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	if passErr != nil {
		http.Error(w, passErr.Error(), http.StatusBadRequest)
	}

	user.Password = string(hashedPass)

	Services.CreateAccount(user)

	fmt.Fprint(w, "User:", user)
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

	Services.GetUser("email", credentials.Email)

	//TODO Compare Passwords
	//
}

func handler() {
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func main() {
	handler()
}
