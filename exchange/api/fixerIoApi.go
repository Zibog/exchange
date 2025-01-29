package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var apiKey = os.Getenv("FIXER_IO_API_KEY")
var fixerIoURL = "http://data.fixer.io/api/latest?access_key=" + apiKey

func CallFixerIo() FixerResponse {
	response, err := http.Get(fixerIoURL)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return UnmarshalFixerResponse(responseData)
}

func UnmarshalFixerResponse(data []byte) FixerResponse {
	var responseObject FixerResponse
	json.Unmarshal(data, &responseObject)

	return responseObject
}

// TODO: unite both responses
type FixerResponse struct {
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp,omitempty"`
	Base      string             `json:"base,omitempty"`
	Date      string             `json:"date,omitempty"`
	Rates     map[string]float64 `json:"rates,omitempty"`
	Error     FixerError         `json:"error,omitempty"`
	Symbols   map[string]string  `json:"symbols,omitempty"`
}

type FixerError struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

func (fr FixerResponse) String() string {
	return fmt.Sprintf("Success: %t\nTimestamp: %v\nBase: %s\nDate: %s\nRates: %v\n",
		fr.Success, fr.Timestamp, fr.Base, fr.Date, fr.Rates)
}
