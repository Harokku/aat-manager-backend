package gsuite

import (
	"aat-manager/utils"
	"context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
)

// SheetService represents a service for interacting with Google Sheets API.
// The SheetService struct contains the following fields:
// - Srv: A pointer to the sheets.Service object that represents the Google Sheets API service.
// - vehicleSheet: The name of the sheet that contains vehicle data.
// - stationSheet: The name of the sheet that contains station data.
type SheetService struct {
	Srv          *sheets.Service
	vehicleSheet string
	stationSheet string
}

// New initializes a new instance of the SheetService struct.
// It creates a new Sheet service by retrieving the credentials from the client file specified in the environment variable utils.GOOGLECREDENTIAL.
// It returns the new SheetService instance and an error if any.
func (ss SheetService) New() (SheetService, error) {
	ctx := context.Background()
	b := utils.ReadEnvOrPanic(utils.GOOGLECREDENTIAL)

	config, err := google.ConfigFromJSON([]byte(b), sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client file from config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheet client: %v", err)
	}

	// Read sheets id from env
	vehicleSheet := utils.ReadEnvOrPanic(utils.VEHICLESHEETID)
	stationSheet := utils.ReadEnvOrPanic(utils.STATIONSHEETID)

	return SheetService{
		Srv:          srv,
		vehicleSheet: vehicleSheet,
		stationSheet: stationSheet,
	}, nil
}
