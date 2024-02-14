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
	"strconv"
)

func main() {
	// Check that all env variable are set
	utils.CheckEnvCompliance()

	// Initialize in memory db for OTP storage
	memoryDb := db.NewDB()

	// Create postgres db conn pool and ping DB

	// Read google service enable flag
	googleServiceEnable, err := strconv.ParseBool(utils.ReadEnvOrPanic(utils.WITHGOOGLESERVICE))
	if err != nil {
		googleServiceEnable = false
	}

	var handler handlers.Handler
	// Enable Google service according to env flag
	if googleServiceEnable {
		mailService, err := gsuite.MailService{}.New()
		if err != nil {
			log.Fatalf("Error initializing mail service:\t%s\n", err)
		}

		// Create handler to setup routes
		handler.InitializeService(memoryDb, mailService, true)

	} else {
		handler.InitializeService(nil, gsuite.MailService{}, false)
	}

	// Fiber app definition
	app := fiber.New()

	// Google verify file
	app.Static("/", "./public/googleverify")

	// Login files
	app.Static("/login", "./public/login")

	// Apply middleware to the app
	app.Use(logger.New())
	app.Use(cors.New())

	routing.SetupRoutes(app, handler)

	port := utils.ReadEnvOrPanic(utils.PORT)
	app.Listen(":" + port)
}
