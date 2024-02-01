package authenticator

import (
	"aat-manager/db"
	"aat-manager/utils"
	"net/mail"
	"strconv"
	"strings"
)

// GenOtpAndSave generates a one-time password (OTP) and saves it in an in-memory database.
// The function takes a mail address and an in-memory database as input parameters.
// It returns the generated OTP on success or an error if any occurs.
//
// The function first searches for the last occurrence of "@" in the mail address. If "@" doesn't exist,
// it returns an error of ErrMalformedMail.
//
// It then checks if the domain of the mail address is authorized by comparing it with the value of the environment variable utils.AUTHORIZEDDOMAIN.
// If the domain is not authorized, it returns an error of ErrUnauthorizedDomain.
//
// Next, it generates a new OTP based on the value of the environment variable utils.OTPLENGTH.
// If the value is not a numeric value, it returns an error of ErrNonNumericValue.
// The OTP length is determined by converting the environment variable value to an integer.
// The OTP is generated using the randSecret function.
//
// After generating the OTP, the function extracts the user part of the mail address (before the "@" character)
// and saves the user and OTP pair in the in-memory database using the db.Set method.
// The user is used as the key and the OTP as the value.
//
// Finally, the function returns the generated OTP and a nil error.
func GenOtpAndSave(mail mail.Address, db *db.InMemoryDb) (string, error) {
	// Search for last @ occurrence
	atIndex := strings.LastIndex(mail.Address, "@")
	if atIndex == -1 {
		return "", ErrMalformedMail
	}

	// Check if the domain is authorized
	authDomain := utils.ReadEnvOrPanic(utils.AUTHORIZEDDOMAIN)
	domain := mail.Address[atIndex+1:]
	if domain != authDomain {
		return "", ErrUnauthorizedDomain
	}

	// Generate new otp
	otpLength, err := strconv.Atoi(utils.ReadEnvOrPanic(utils.OTPLENGTH))
	if err != nil {
		return "", ErrNonNumericValue
	}
	otp := randSecret(otpLength)

	// Save user (mailbox from mail) and OTP to in memory db
	user := mail.Address[:atIndex]
	db.Set(user, otp)

	// Return OTP
	return otp, nil
}

// CheckOtpAndDelete checks if the passed OTP is equal to the stored OTP for a given user in the in-memory database
// deletes the user from the database if the OTP check is passed
func CheckOtpAndDelete(mail mail.Address, otp int, db *db.InMemoryDb) (bool, error) {
	// Extract user from mail address
	// Search for last @ occurrence
	atIndex := strings.LastIndex(mail.Address, "@")
	if atIndex == -1 {
		return false, ErrMalformedMail
	}
	user := mail.Address[:atIndex]

	// Retrieve user/OTP pair from in memory db and convert to int
	storedOtp, exist := db.Get(user)
	if !exist {
		return false, ErrUserNotFound
	}
	storedIntOtp, err := strconv.Atoi(storedOtp)
	if err != nil {
		return false, ErrNonNumericValue
	}

	// Check if passed OTP is equal to stored one
	// Return if check is passed
	if otp == storedIntOtp {
		return true, nil
	}
	return false, nil
}
