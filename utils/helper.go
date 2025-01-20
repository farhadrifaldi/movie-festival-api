package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

/*
*
Function to get environment variable from .env
*/
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot load env")
		os.Exit(1)
	}

	return os.Getenv(key)
}
