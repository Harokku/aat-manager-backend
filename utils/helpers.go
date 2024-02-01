package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	PORT             = "PORT"           // Serve port
	JWTSECRET        = "JWTSECRET"      // Secret for JWT signing
	JWTEXPIREINMONTH = "JWTEXPIREM"     // JWT expiration in month
	AUTHORIZEDDOMAIN = "AUTHDOMAIN"     // Authorized e-mail domain for login
	OTPLENGTH        = "OTPLENGTH"      // Length of the generated numerical OTP in character
	GOOGLECREDENTIAL = "GSECRET"        // Google API credential JSON
	VEHICLESHEETID   = "VEHICLESHEETID" // Sheet ID for vehicle issue report
	STATIONSHEETID   = "STATIONSHEETID" // Sheet ID for station issue report
)

// CheckEnvCompliance verifies that all required environment variables are set.
// It checks each variable in the envList and panics if any of them is not set.
// It uses the ReadEnvOrPanic function to read the value of each variable.
func CheckEnvCompliance() {
	envList := []string{
		PORT,
		JWTSECRET,
		JWTEXPIREINMONTH,
		AUTHORIZEDDOMAIN,
		OTPLENGTH,
		GOOGLECREDENTIAL,
		VEHICLESHEETID,
		STATIONSHEETID,
	}

	// Cycle environment variable list and panic if anyone is not set
	for _, env := range envList {
		_ = ReadEnvOrPanic(env)
	}
}

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