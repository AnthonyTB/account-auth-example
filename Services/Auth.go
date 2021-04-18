package Services

import (
	"github.com/anthonytb/account-auth-example/Utils"
	"github.com/dgrijalva/jwt-go"
)

// Function for creating JWT for user
// Params: user<user>
// Returns: string<pointer> or nil
func CreateJWT(user User) *string {

	// Creates variable w/ empty object as value
	claims := jwt.MapClaims{}

	// Assigns keys and values to the claims object
	claims["authorized"] = true
	claims["user_id"] = user.id

	// Creates base JWT
	init := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Hashes JWT with JWT secret then returns the token or an error if it encountered one
	token, tokenErr := init.SignedString([]byte(Utils.GoDotEnvVariable("JWT_SECRET")))

	// If hashing the JWT encountered an error then returns nil
	if tokenErr != nil {
		return nil
	}

	// Returns token if no error occured
	return &token
}
