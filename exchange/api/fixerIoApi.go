package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/augurysys/timestamp"
)

var apiKey = os.Getenv("FIXER_IO_API_KEY")
var fixerIoURL = "http://data.fixer.io/api/latest?access_key=" + apiKey

func CallFixerIo() FixerResponse {
	// response, err := http.Get(fixerIoURL)
	//
	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	os.Exit(1)
	// }
	//
	// responseData, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var responseObject FixerResponse
	json.Unmarshal([]byte(input), &responseObject)

	fmt.Println(responseObject)

	return responseObject
}

// TODO: unite both responses
type FixerResponse struct {
	Success   bool                `json:"success"`
	Timestamp timestamp.Timestamp `json:"timestamp"`
	Base      string              `json:"base"`
	Date      string              `json:"date"`
	Rates     map[string]float64  `json:"rates"`
}

func (fr FixerResponse) String() string {
	return fmt.Sprintf("Success: %t\nTimestamp: %s\nBase: %s\nDate: %s\nRates: %v\n",
		fr.Success, fr.Timestamp, fr.Base, fr.Date, fr.Rates)
}

const input = `{
    "success": true,
    "timestamp": 1519296206,
    "base": "EUR",
    "date": "2025-01-28",
    "rates": {
        "AUD": 1.566015,
        "CAD": 1.560132,
        "CHF": 1.154727,
        "CNY": 7.827874,
        "GBP": 0.882047,
        "JPY": 132.360679,
        "USD": 1.23396
    }
}`
