package main

import (
	"aat-manager/db"
	"aat-manager/gsuite"
	"aat-manager/handlers"
	"aat-manager/routing"
	"aat-manager/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func main() {

	// Initialize in memory db for OTP storage
	memoryDb := db.NewDB()

	// Initialize mail service
	mailService, err := gsuite.MailService{}.New()
	if err != nil {
		log.Fatalf("Error initializing mail service:\t%s\n", err)
	}

	// Create handler to setup routes
	handler := &handlers.Handler{
		Db:          memoryDb,
		MailService: mailService,
	}

	// Fiber app definition
	app := fiber.New()

	// Login files
	app.Static("/login", "./public/login")

	// Apply middleware to the app
	app.Use(logger.New())
	app.Use(cors.New())

	routing.SetupRoutes(app, *handler)

	port := utils.ReadEnvOrPanic("PORT")
	app.Listen(":" + port)
}
