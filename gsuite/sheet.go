package gsuite

import (
	"aat-manager/utils"
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"sync"
)

const (
	VehicleSheet = "vehicleSheet"
	StationSheet = "stationSheet"
)

// SheetService represents a service for interacting with Google Sheets API.
// The SheetService struct contains the following fields:
// - Srv: A pointer to the sheets.Service object that represents the Google Sheets API service.
// - vehicleSheet: The name of the sheet that contains vehicle data.
// - stationSheet: The name of the sheet that contains station data.
type SheetService struct {
	Srv      *sheets.Service
	sheets   map[string]string
	initOnce sync.Once
	initErr  error
}

// Initialize lazy initialize the sheet service when needed.
// It call initialize function and return error in fail
func (ss *SheetService) Initialize() error {
	ss.initOnce.Do(func() {
		ss.initErr = ss.initialize()
	})
	return ss.initErr
}

// initialize initializes the SheetService by setting up the Google Sheet client and reading the sheet IDs from the environment.
// It takes no parameters and returns an error if any occurred during initialization.
func (ss *SheetService) initialize() error {
	ctx := context.Background()
	b := utils.ReadEnvOrPanic(utils.GOOGLECREDENTIAL)

	config, err := google.ConfigFromJSON([]byte(b), gmail.GmailSendScope, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client file from config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheet client: %v", err)
	}

	ss.Srv = srv

	// Read sheets id from env
	vSheet := utils.ReadEnvOrPanic(utils.VEHICLESHEETID)
	sSheet := utils.ReadEnvOrPanic(utils.STATIONSHEETID)

	ss.sheets = map[string]string{
		VehicleSheet: vSheet,
		StationSheet: sSheet,
	}

	return nil
}

// Append appends data to a specific sheet in a Google Sheet.
// It takes the sheet identifier 's', the range 'r' where the data should be appended, and the data to be appended 'data' as input.
// The 'data' parameter should be a 2-dimensional slice of interface{} where each element of the slice represents a row of data and each element of a row represents a cell value.
// The function returns the HTTP status code of the request and an error if any.
func (ss *SheetService) Append(s string, r string, data [][]interface{}) (int, error) {
	if err := ss.Initialize(); err != nil {
		return 500, err
	}

	var values = sheets.ValueRange{
		Values: data,
	}

	res, err := ss.Srv.Spreadsheets.Values.Append(ss.sheets[s], r, &values).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return 500, err
	}

	return res.HTTPStatusCode, nil
}

// EnumerateSheets retrieves the list of sheet names in the specified spreadsheet.
// It takes a string parameter `s` which represents the spreadsheet ID.
// It returns a slice of strings containing the names of the sheets and an error if any.
// If an error occurs while retrieving the spreadsheet or its sheets, it will be returned.
func (ss *SheetService) EnumerateSheets(s string) ([]string, error) {
	if err := ss.Initialize(); err != nil {
		return nil, err
	}

	spreadsheet, err := ss.Srv.Spreadsheets.Get(s).Do()
	if err != nil {
		return nil, err
	}
	var sheetNames []string
	for _, sheet := range spreadsheet.Sheets {
		sheetNames = append(sheetNames, sheet.Properties.Title)
	}
	return sheetNames, nil
}

// GetAllRecords retrieves all records from the specified sheet in a spreadsheet.
// It takes in the spreadsheet ID (s), sheet name (sheet), and the column number (coln) to define the range from which to get the records.
// It returns a 2D slice of interfaces representing the retrieved records and an error if any.
// The records are retrieved from the range "sheet!A2:{colName}" where colName is the corresponding column name calculated from coln using the colNumToName function.
// If an error occurs while retrieving the records or if no values are found in the response, the function returns nil and the error.
// Otherwise, it returns the retrieved records.
// Example usage:
//
//	record, err := sheetService.GetAllRecords("<spreadsheetID>", "<sheetName>", 5)
func (ss *SheetService) GetAllRecords(s string, sheet string, coln int) ([][]interface{}, error) {
	if err := ss.Initialize(); err != nil {
		return nil, err
	}

	if coln <= 0 {
		return nil, fmt.Errorf("invalid column number: %d", coln)
	}

	cellRange := fmt.Sprintf("%s!A2:%s", sheet, colNumToName(coln))
	var result [][]interface{}

	resp, err := ss.Srv.Spreadsheets.Values.Get(s, cellRange).Do()
	if err != nil {
		return nil, err
	}

	result = resp.Values

	return result, nil
}
