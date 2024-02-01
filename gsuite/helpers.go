package gsuite

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"time"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWebToConsole(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// create a random state string
	state := fmt.Sprintf("st%d", time.Now().UnixNano())

	// Generate the OAuth2.0 URL
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// make a channel to pass the token or error
	tokenCh := make(chan *oauth2.Token)
	errorCh := make(chan error)

	// create a temporary web server
	server := http.NewServeMux()
	s := &http.Server{Addr: ":8080", Handler: server}

	// handle OAuth2.0 responses
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// confirm the state matches
		if r.URL.Query().Get("state") != state {
			errorCh <- fmt.Errorf("state did not match")
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		// use the authorization code that is pushed to the redirect URL to fetch the access and refresh tokens
		tok, err := config.Exchange(context.TODO(), r.URL.Query().Get("code"))
		if err != nil {
			errorCh <- fmt.Errorf("failed to exchange token: %s", err.Error())
			http.Error(w, "failed to exchange token", http.StatusInternalServerError)
			return
		}

		tokenCh <- tok
	})

	// Start the web server
	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errorCh <- err
		}
	}()

	// print the url to authorize
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)

	var token *oauth2.Token
	select {
	case err := <-errorCh:
		log.Fatalf("Unable to retrieve token from web: %v", err)
	case token = <-tokenCh:
	}

	// shut down the server
	s.Shutdown(context.TODO())

	return token
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
