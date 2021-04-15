package Services

import (
	"github.com/anthonytb/account-auth-example/Utils"
	"github.com/dgrijalva/jwt-go"
)

func CreateJWT(user User) *string {

	claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["user_id"] = user.id

	init := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, tokenErr := init.SignedString([]byte(Utils.GoDotEnvVariable("JWT_SECRET")))

	if tokenErr != nil {
		return nil
	}

	return &token

}
