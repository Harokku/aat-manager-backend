package authenticator

import (
	"aat-manager/utils"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

// CreateAndSignJWT creates and signs a JSON Web Token (JWT) with the specified user and manager information.
// It returns the generated token as a string, along with any error encountered.
func CreateAndSignJWT(user string, manager bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Expire days from env
	expiredays, err := strconv.Atoi(utils.ReadEnvOrPanic(utils.JWTEXPIREINMONTH))
	if err != nil {
		return "", ErrNonNumericValue
	}

	// Token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user
	claims["manager"] = manager
	claims["exp"] = time.Now().AddDate(0, expiredays, 0).Unix()

	// Read secret from env and sig the token
	jwtSecret := utils.ReadEnvOrPanic(utils.JWTSECRET)
	// Check if blank secret is used and return error
	if jwtSecret == "" {
		return "", ErrBlankSecret
	}
	// Sign the token using provided secret
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
