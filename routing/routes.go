package routing

import (
	"aat-manager/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, handler handlers.Handler) {
	// Root landing page
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("App root")
	})

	// OAuth auth route
	app.Get("/oauth_callback", handlers.OauthCallback())

	// Login routes
	login := app.Group("/login")
	login.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("Login root")
	})
	login.Post("/", handler.GetMailAndSendBackOtp)
	login.Post("/checkotp", handler.GetOtpAndAuthenticate)

	//Api routes
	api := app.Group("/api/v1")

	protected := api.Group("/", handlers.JWTAuthenticationMiddleware)
	protected.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("Protected root")
	})
}
