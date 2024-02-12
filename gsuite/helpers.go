package gsuite

import (
	"aat-manager/utils"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type Token struct {
	Token *oauth2.Token // oauth2 signed Token
	Err   error         // error
}

type SharedState struct {
	state string
	mux   sync.Mutex
}

var sharedState = &SharedState{}

var TokenCh = make(chan Token) // TokenCh to xchange oauth2 token

// GetState Getter for state
func GetState() string {
	sharedState.mux.Lock()
	defer sharedState.mux.Unlock()

	return sharedState.state
}

// SetState Setter for state
func SetState(state string) {
	sharedState.mux.Lock()
	defer sharedState.mux.Unlock()

	sharedState.state = state
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		useWebAuth, _ := strconv.ParseBool(utils.ReadEnvOrPanic(utils.WEBAUTH))
		if useWebAuth {
			getTokenFromWeb(config)
			go func() {
				// Listen to channel for signed token or error
				tok := <-TokenCh
				if tok.Err != nil {
					log.Printf("Unable to retrieve token from web: %v", tok.Err)
				}
				saveToken(tokFile, tok.Token)
			}()
		} else {
			tok = getTokenFromWebToConsole(config)
			saveToken(tokFile, tok)
		}
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

// getTokenFromWeb generates an OAuth2.0 URL, sets a random state string, and prints the authorization URL for the user to follow in their browser to authorize the application.
// The generated state is set as a shared state for later processing.
// This function does not return any value.
func getTokenFromWeb(config *oauth2.Config) {
	// create a random state string
	state := fmt.Sprintf("st%d", time.Now().UnixNano())

	// Set state for handler to process
	SetState(state)

	// Generate the OAuth2.0 URL
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// print the url to authorize
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)

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

// -------------------------
// Gsheet helpers
// -------------------------

// checkA1Validity checks the validity of a string in the form "Sheet1!A1:B2".
// It uses a regular expression to match the format.
// Returns true if the string matches the format, false otherwise.
func checkA1Validity(s string) bool {
	// This regex will match strings of the form "Sheet1!A1:B2". Note it doesn't work with Named Ranges.
	rx := regexp.MustCompile(`^[^\s!]+![A-Z]+\d+:[A-Z]+\d+$`)

	return rx.MatchString(s)
}

// colNumToName converts a column number to its corresponding column name in Excel spreadsheet.
// The function takes an integer column number as input and returns the corresponding column name as a string.
// If the input column number is 0, an empty string is returned.
// The conversion is based on a 26-letter system, where A = 1, B = 2, ..., Z = 26.
// For column numbers greater than 26, multiple characters are used to represent the column name.
// For example, 27 is represented as "AA", 52 is represented as "AZ", and 700 is represented as "ZX".
// If a non-positive column number is provided, an empty string is returned.
func colNumToName(colNum int) string {
	colName := "" // Start with empty string
	if colNum == 0 {
		return colName
	}
	for colNum > 0 {
		rem := (colNum - 1) % 26
		colNum = (colNum - rem) / 26
		colName = string('A'+rune(rem)) + colName
	}
	return colName
}
