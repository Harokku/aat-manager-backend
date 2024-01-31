package handlers

import (
	"aat-manager/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
)

func JWTAuthenticationMiddleware(ctx *fiber.Ctx) error {
	jwtCookie := ctx.Cookies(JWTTokenCookieName)            // JWT auth token
	pendingAuthCookie := ctx.Cookies(PendingAuthCookieName) // Pending auth status

	// Check if JWT is present in cookie:
	// Ok -> Proceed with validation
	// Ko -> Check if auth is pending waiting OTP (pending auth cookie is present)
	//
	// Check if auth is pending waiting OTP
	// OK -> Redirect to OTP check page
	// KO -> Redirect to login page
	if jwtCookie == "" {
		redirectPage := LoginURL

		if pendingAuthCookie != "" {
			redirectPage = CheckOTPURL
		}

		redirectPath := getRedirectPath(redirectPage, ctx)
		ctx.Redirect(redirectPath, fiber.StatusTemporaryRedirect)
		return nil
	}

	jwtSecret := utils.ReadEnvOrPanic("JWTSECRET")
	if jwtSecret == "" {
		log.Println("JWT Secret not set")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error: JWT Secret not set.")
	}

	// Parse token after signing method verification
	token, err := jwt.Parse(jwtCookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	// If parsing error (apart from blank secret), redirect to login
	if err != nil {
		log.Printf("Failed to parse JWT: %v", err)
		redirectPath := getRedirectPath(LoginURL, ctx)
		ctx.Redirect(redirectPath, fiber.StatusTemporaryRedirect)
		return nil
	}
	if !token.Valid {
		log.Println("Provided JWT is not valid")
		redirectPath := getRedirectPath(LoginURL, ctx)
		ctx.Redirect(redirectPath, fiber.StatusTemporaryRedirect)
		return nil
	}

	return ctx.Next()
}

func AuthPendingMiddleware(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("pendingauth")
	if cookie == "" {
		// No intermediate auth cookie found -> redirect to login
		redirectPath := getRedirectPath(LoginURL, ctx)
		ctx.Redirect(redirectPath, fiber.StatusTemporaryRedirect)
		return nil
	}

	// If present, continue to next handler
	return ctx.Next()
}
