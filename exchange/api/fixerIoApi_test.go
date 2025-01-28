package api

import "testing"

func TestRequestLatestRates(t *testing.T) {
	response := CallFixerIo()
	if response.Success != true {
		t.Errorf("Expected success to be true, but got %t", response.Success)
	}
}

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

const error = `{
  "success": false,
  "error": {
    "code": 104,
    "info": "Your monthly API request volume has been reached. Please upgrade your plan."
  }
}`

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
