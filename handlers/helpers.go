package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

const (
	JWTTokenCookieName    = "jwt"
	PendingAuthCookieName = "pendingauth"
	LoginURL              = "/login/login.html"
	CheckOTPURL           = "/login/checkotp.html"
)

// getRedirectPath constructs the redirection URL with the current path.
func getRedirectPath(page string, ctx *fiber.Ctx) string {
	return fmt.Sprintf("%s?redirect=%s", page, ctx.Path())
}

// getChainableRedirectPath constructs a redirection URL with an additional redirect query parameter.
// The redirectPage parameter is extracted from the "redirect" query string of the given fiber.Ctx instance.
// The constructed URL is formed by appending the redirectPage to the provided page parameter.
func getChainableRedirectPath(page string, ctx *fiber.Ctx) string {
	redirectPage := ctx.Query("redirect")
	return fmt.Sprintf("%s?redirect=%s", page, redirectPage)
}
