package handlers

import (
	"aat-manager/authenticator"
	"aat-manager/db"
	"aat-manager/gsuite"
	"aat-manager/utils"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/sheets/v4"
	"net/mail"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	Db          *db.InMemoryDb     // In memory db interface
	MailService gsuite.MailService // Gmail service interface

	initialized bool // Indicate that the handler is initialized and safe for use
}
type AuthData struct {
	Mail string `json:"mail,omitempty" form:"mail"`
	Otp  string `json:"otp,omitempty" form:"otp"`
}

func OauthCallback() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		// Read state from auth request
		state := <-gsuite.StateCh

		// Check the state parameter
		responseState := ctx.Query("state")
		if responseState != state {
			return ctx.Status(fiber.StatusBadRequest).SendString("Invalid state parameter.")
		}

		// Get the authorization code from the response
		code := ctx.Query("code")

		// Recreate config from env
		b := utils.ReadEnvOrPanic(utils.GOOGLECREDENTIAL)
		config, err := google.ConfigFromJSON([]byte(b), gmail.GmailSendScope, sheets.SpreadsheetsScope)

		// Exchange the authorization code for an access token
		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			// If error, send back error to caller
			gsuite.TokenCh <- gsuite.Token{
				Err: err,
			}
			// If error, send error response
			return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token: " + err.Error())
		}

		// Send back token to caller
		gsuite.TokenCh <- gsuite.Token{
			Token: token,
			Err:   nil,
		}

		return ctx.Status(fiber.StatusOK).SendString("Authentication successful.")
	}
}

func (h *Handler) InitializeService(db *db.InMemoryDb, ms gsuite.MailService, init bool) {
	h.Db = db
	h.MailService = ms
	h.initialized = init
}

// GetMailAndSendBackOtp takes in a fiber.Ctx and retrieves the email from the request body.
// It then validates the email and generates an OTP for the user.
// The OTP is sent to the user
func (h *Handler) GetMailAndSendBackOtp(ctx *fiber.Ctx) error {
	if !h.initialized {
		return ctx.Status(fiber.StatusNotImplemented).SendString("This service is not enabled.")
	}

	formData := new(AuthData)

	// Read email field from request
	if err := ctx.BodyParser(formData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Parse email from request to validate it
	addr, err := mail.ParseAddress(formData.Mail)
	if err != nil {
		log.Errorf("Error parsing email from auth form:\t%s\n", err)
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Generate OTP for user
	otp, err := authenticator.GenOtpAndSave(*addr, h.Db)

	// Send OTP to user by e-mail
	err = h.MailService.SendMail("Codice di verifica", addr.Address, "Ecco il tuo codice di verifica:\t"+otp)
	if err != nil {
		log.Errorf("Error senting OTP mail:\t%s\n", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Set pending auth cookie
	ctx.Cookie(&fiber.Cookie{
		Name:        "pendingauth",
		Value:       addr.Address,
		Expires:     time.Now().Add(time.Minute * 3),
		Secure:      false,
		HTTPOnly:    true,
		SameSite:    "lax",
		SessionOnly: false,
	})

	redirectPath := getChainableRedirectPath(CheckOTPURL, ctx)
	ctx.Redirect(redirectPath, fiber.StatusSeeOther)
	return nil
}

// GetOtpAndAuthenticate takes in a fiber.Ctx and retrieves the user-entered OTP from the request body.
// It converts the OTP from string to int and reads the user's email from the pendingAuthCookie.
// It then checks if the user-entered OTP matches the stored OTP and returns the appropriate status and message.
func (h *Handler) GetOtpAndAuthenticate(ctx *fiber.Ctx) error {
	if !h.initialized {
		return ctx.Status(fiber.StatusNotImplemented).SendString("This service is not enabled.")
	}

	formData := new(AuthData)

	// Read user entered OTP from request
	if err := ctx.BodyParser(formData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Convert OTP from form to int
	otp, err := strconv.Atoi(formData.Otp)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Read user's email from pendingAuthCookie
	userEmail, err := mail.ParseAddress(ctx.Cookies(PendingAuthCookieName))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Check if user entered OTP match with stored one
	valid, err := authenticator.CheckOtpAndDelete(*userEmail, otp, h.Db)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// If valid create JWT, set it in a cookie and redirect to original url
	if valid {
		log.Info("Valid OTP")

		// Extract user from userEmail address
		// Search for last @ occurrence
		atIndex := strings.LastIndex(userEmail.Address, "@")
		if atIndex == -1 {
			return ctx.Status(fiber.StatusBadRequest).SendString("Malformed email address")
		}
		user := userEmail.Address[:atIndex]

		// Create and sign token
		token, err := authenticator.CreateAndSignJWT(user, false)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Read JWT expiration from env
		expiredays, err := strconv.Atoi(utils.ReadEnvOrPanic(utils.JWTEXPIREINMONTH))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Set auth cookie
		ctx.Cookie(&fiber.Cookie{
			Name:        "jwt",
			Value:       token,
			Expires:     time.Now().AddDate(0, expiredays, 0),
			Secure:      false,
			HTTPOnly:    true,
			SameSite:    "lax",
			SessionOnly: false,
		})

		// Clear auth pending cookie
		ctx.ClearCookie(PendingAuthCookieName)

		// Redirect to original url
		redirectPage := ctx.Query("redirect")
		ctx.Redirect(redirectPage, fiber.StatusSeeOther)
		return nil
	} else {
		// OTP il invalid, keep on same page till cookie expiration or valid code inserted
		log.Error("Invalid OTP")
		redirectPage := getChainableRedirectPath(CheckOTPURL, ctx)
		ctx.Redirect(redirectPage, fiber.StatusSeeOther)
		return nil
	}
}
