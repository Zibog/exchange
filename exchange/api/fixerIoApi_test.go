package api

import (
	"testing"
)

func TestParseRequestLatestRates(t *testing.T) {
	// http://data.fixer.io/api/latest
	const latest = `{
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

	response := UnmarshalFixerResponse([]byte(latest))
	AssertEquals(t, true, response.Success)
	AssertEquals(t, int64(1519296206), response.Timestamp)
	AssertEquals(t, "EUR", response.Base)
	AssertEquals(t, "2025-01-28", response.Date)
	AssertEquals(t, 7, len(response.Rates))

	AssertEquals(t, FixerError{}, response.Error)
	AssertEquals(t, 0, len(response.Symbols))
}

func TestParseErrorResponse(t *testing.T) {
	const error = `{
    "success": false,
    "error": {
      "code": 104,
      "info": "Your monthly API request volume has been reached. Please upgrade your plan."
    }
  }`

	response := UnmarshalFixerResponse([]byte(error))
	AssertEquals(t, false, response.Success)
	AssertEquals(t, 104, response.Error.Code)
	AssertEquals(t, "Your monthly API request volume has been reached. Please upgrade your plan.", response.Error.Info)

	AssertEquals(t, 0, len(response.Rates))
	AssertEquals(t, 0, len(response.Symbols))
}

func TestParseRequestSymbols(t *testing.T) {
	// http://data.fixer.io/api/symbols
	const symbols = `{
    "success": true,
    "symbols": {
      "AED": "United Arab Emirates Dirham",
      "AFN": "Afghan Afghani",
      "ALL": "Albanian Lek",
      "AMD": "Armenian Dram"
      }
  }`

	response := UnmarshalFixerResponse([]byte(symbols))
	AssertEquals(t, true, response.Success)
	AssertEquals(t, 4, len(response.Symbols))

	AssertEquals(t, 0, len(response.Rates))
	AssertEquals(t, FixerError{}, response.Error)
}

func TestToUrl(t *testing.T) {
	url := toUrl("http://data.fixer.io", Latest, "MY_API_KEY")
	AssertEquals(t, "http://data.fixer.io/latest?access_key=MY_API_KEY", url)
}

func TestAppendSymbols(t *testing.T) {
	// TODO: add tests for symbols
}

func TestParseRequestLatestRatesWithSymbolsAndBase(t *testing.T) {
	// http://data.fixer.io/api/latest?symbols=GBP,JPY,EUR&base=USD
	const latestWithSymbolsAndBase = `{
      "success": true,
      "timestamp": 1519296206,
      "base": "USD",
      "date": "2025-01-28",
      "rates": {
          "GBP": 0.72007,
          "JPY": 107.346001,
          "EUR": 0.813399
      }
  }`

	response := UnmarshalFixerResponse([]byte(latestWithSymbolsAndBase))
	AssertEquals(t, true, response.Success)
	AssertEquals(t, int64(1519296206), response.Timestamp)
	AssertEquals(t, "USD", response.Base)
	AssertEquals(t, "2025-01-28", response.Date)
	AssertEquals(t, 3, len(response.Rates))

	AssertEquals(t, FixerError{}, response.Error)
	AssertEquals(t, 0, len(response.Symbols))
}

func AssertEquals(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

// http://data.fixer.io/api/latest?callback=MY_FUNCTION
const latestWithCallback = `({
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
})`
