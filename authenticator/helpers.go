package authenticator

import (
	"errors"
	"math/rand"
)

// Error definition
var (
	ErrMalformedMail      = errors.New("malformed mail")
	ErrUnauthorizedDomain = errors.New("domain not authorized")
	ErrNonNumericValue    = errors.New("value is not a number")
	ErrUserNotFound       = errors.New("user not found")
	ErrBlankSecret        = errors.New("blank secret key used")
)

// Secret generator dictionary
const secretBytes = "0123456789"

// randSecret generates a random secret string of the given length.
// It uses the characters from the secretBytes constant to generate the secret.
// The length of the generated secret will be equal to the given length.
// If the given length is 0, an empty string will be returned.
// The generated secret string will be a combination of random characters from the secretBytes constant.
//
// Example usage:
//
//	secret := randSecret(8) // generates a random secret string with length 8
//	fmt.Println(secret)    // output: "92638475"
//
// Note:
// The constant secretBytes should contain all the characters that can be used in the generated secret string.
func randSecret(n int) string {
	// Make a byte array of len = secret length
	b := make([]byte, n)

	// For every byte in b generate a random digit from secretByte const
	for i := range b {
		b[i] = secretBytes[rand.Intn(len(secretBytes))]
	}

	return string(b)
}
