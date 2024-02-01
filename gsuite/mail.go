package gsuite

import (
	"aat-manager/utils"
	"bytes"
	"context"
	"encoding/base64"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
)

type MailService struct {
	srv *gmail.Service
}

// New initializes and returns a new instance of the MailService struct.
// It requires a credentials.json file to be present in the working directory.
// It returns the new MailService struct with a Gmail service and any error encountered while initializing.
func (ms MailService) New() (MailService, error) {
	ctx := context.Background()
	b := utils.ReadEnvOrPanic(utils.GOOGLECREDENTIAL)

	config, err := google.ConfigFromJSON([]byte(b), gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	return MailService{srv: srv}, nil
}

// SendMail sends an email with the given subject, recipient, and message.
// It returns an error if the email sending fails.
func (ms MailService) SendMail(subject string, to string, message string) error {
	var email bytes.Buffer
	//email.WriteString("From: AAT One Time Password Provider\r\n")
	email.WriteString("To: " + to + "\r\n")
	email.WriteString("Subject: " + subject + "\r\n")
	email.WriteString("\r\n" + message)

	raw := base64.URLEncoding.EncodeToString(email.Bytes())

	var msg gmail.Message
	msg.Raw = raw

	_, err := ms.srv.Users.Messages.Send("me", &msg).Do()
	if err != nil {
		return err
	}

	return nil
}
