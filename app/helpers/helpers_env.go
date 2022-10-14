package helpers

import (
	"os"

	"github.com/joho/godotenv"
)

// This function loads the .env file
func LoadEnv() {
	godotenv.Load()
}

// This function loads the .env file and returns the value of the key
func GetEnv(key string) string {
	return os.Getenv(key)
}
