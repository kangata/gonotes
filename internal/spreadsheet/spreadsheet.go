package spreadsheet

import (
	"context"
	"io/ioutil"

	"github.com/kangata/gonotes/internal/client"
	"github.com/kangata/gonotes/internal/env"
	"github.com/kangata/gonotes/models"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Spreadsheet struct {
	SpreadsheetID string
	WriteRange    string
	ValueRange    sheets.ValueRange
	Service       sheets.Service
}

func (s *Spreadsheet) NewService() (*sheets.Service, error) {
	s.SpreadsheetID = env.Get("SPREADSHEET_ID")
	s.WriteRange = "A1"

	ctx := context.Background()

	contentBytes, err := ioutil.ReadFile("storage/credentials.json")

	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(contentBytes, "https://www.googleapis.com/auth/spreadsheets")

	if err != nil {
		return nil, err
	}

	client := client.GetClient(config)

	Service, err := sheets.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		return nil, err
	}

	s.Service = *Service

	return &s.Service, nil
}

func (s *Spreadsheet) AddItem(m models.Message) (*sheets.AppendValuesResponse, error) {
	valueRange := sheets.ValueRange{}

	valueRange.Values = append(valueRange.Values, m.ToSheetRow())

	return s.Service.Spreadsheets.Values.Append(s.SpreadsheetID, s.WriteRange, &valueRange).ValueInputOption("RAW").Do()
}
