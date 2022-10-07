package helpers

import (
	"os"

	"github.com/joho/godotenv"
)

// This function loads the .env file and returns the value of the key
func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	return os.Getenv(key)
}
