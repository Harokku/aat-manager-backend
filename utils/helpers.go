package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// ReadEnvOrPanic reads the value of the specified environment variable with the given name.
// If the variable is not set in the environment, it attempts to load the values from a .env file using godotenv.Load().
// If the .env file cannot be loaded or the variable is still not set, it panics with an appropriate error message.
func ReadEnvOrPanic(name string) string {
	var res string

	res, ok := os.LookupEnv(name)
	if !ok {
		err := godotenv.Load()
		if err != nil {
			log.Panicf("Error loading .env file:\t%s", err)
		}

		res, ok = os.LookupEnv(name)
		if !ok {
			log.Panicf("%s is not set in environment or .env file", name)
		}
	}

	return res
}
