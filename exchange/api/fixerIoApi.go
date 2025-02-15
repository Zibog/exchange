package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Endpoint string

const (
	Latest  Endpoint = "/latest"
	Symbols Endpoint = "/symbols"
)

var apiKey = os.Getenv("FIXER_IO_API_KEY")
var address = os.Getenv("FIXER_IO_ADDRESS")

// TODO: pass symbols as argument. And a callback? Maybe in smth like context
func CallFixerIo(context Context) FixerResponse {
	url := toUrl(address, context.Endpoint, apiKey)
	url = decorateUrl(url, context)
	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err.Error())
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

func (fe FixerError) String() string {
	return fmt.Sprintf("Code: %d\nType: %s\nInfo: %s\n",
		fe.Code, fe.Type, fe.Info)
}

func toUrl(address string, endpoint Endpoint, apiKey string) string {
	// TODO: add validation endpoint<=>symbols
	return address + string(endpoint) + "?access_key=" + apiKey
}

func withSymbols(url string, symbols []string) string {
	if len(symbols) != 0 {
		url += "&symbols=" + strings.Join(symbols, ",")
	}
	return url
}

func withBase(url string, base string) string {
	if base != "" {
		url += "&base=" + base
	}
	return url
}

func decorateUrl(url string, context Context) string {
	url = withSymbols(url, strings.Split(context.Symbols, ","))
	url = withBase(url, context.Base)
	return url
}

type Context struct {
	Endpoint Endpoint
	Symbols  string
	Base     string
}
