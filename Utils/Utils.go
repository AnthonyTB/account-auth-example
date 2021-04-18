package Utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Function for returning enviroment varibles
// Params: varibleKey <string>
// Returns: enviromentVariableValue<string> or logs error if issue loading env
func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
