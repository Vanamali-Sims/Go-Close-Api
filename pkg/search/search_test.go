package search

import (
	"testing"
	"time"
)

// func TestGetCloseUSDForwSpecificDate(t *testing.T) {
// 	assetClass := "crypto"
// 	internalSymbol := "ADA_USDT"
// 	dateString := "2024-06-02T00:59:59Z"
// 	date, err := time.Parse(time.RFC3339, dateString)
// 	if err != nil {
// 		t.Fatalf("Error parsing date: %v", err)
// 	}

// 	// Run the GetCloseUSD function
// 	result, err := GetCloseUSD(assetClass, internalSymbol, date)
// 	if err != nil {
// 		t.Fatalf("Error when calling GetCloseUSD: %v", err)
// 	}

// 	// Print the result for manual verification
// 	fmt.Printf("Result: %+v\n", result)
// }

// func TestGetCloseInBetweenForSpecificDateRange(t *testing.T) {
// 	assetClass := "crypto"
// 	internalSymbol := "1000SATS_USDT"
// 	startDateString := "2024-07-02T15:59:59Z"
// 	endDateString := "2024-07-02T19:59:59Z"

// 	// Run the GetCloseInBetween function
// 	result, err := GetCloseInBetween(assetClass, internalSymbol, startDateString, endDateString)
// 	if err != nil {
// 		t.Fatalf("Error when calling GetCloseInBetween: %v", err)
// 	}

// 	// Print the result for manual verification
// 	fmt.Printf("Result: %+v\n", result)
// }

// func TestGetCloseUSDForSpecificDate(t *testing.T) {
// 	assetClass := "crypto"
// 	internalSymbol := "ADA_USDT"
// 	dateString := "2024-06-02T00:59:59Z"
// 	date, err := time.Parse(time.RFC3339, dateString)
// 	if err != nil {
// 		t.Fatalf("Error parsing date: %v", err)
// 	}

// 	// Run the GetCloseUSD function
// 	result, err := GetCloseUSD(assetClass, internalSymbol, date)
// 	if err != nil {
// 		t.Fatalf("Error when calling GetCloseUSD: %v", err)
// 	}

// 	// Always print the results for examination
// 	fmt.Println("Test Results:")
// 	fmt.Printf("Close Price USD: %v\n", result.ClosePriceUSD)
// 	fmt.Printf("Fetched Date: %s\n", result.Metadata.FetchedDate)
// 	fmt.Printf("Conversion Rate: %v\n", result.Metadata.ConversionRate)
// 	fmt.Printf("Conversion Rate Date: %s\n", result.Metadata.ConversionRateDate)
// 	fmt.Printf("Candle: %s\n", result.Metadata.Candle)
// }

// func TestGetCloseInBetween(t *testing.T) {
// 	assetClass := "crypto"
// 	internalSymbol := "ADA_USDT"
// 	startDateString := "2024-06-01T00:00:00Z"
// 	endDateString := "2024-06-03T00:00:00Z"

// 	// Run the GetCloseInBetween function
// 	result, err := GetCloseInBetween(assetClass, internalSymbol, startDateString, endDateString)
// 	if err != nil {
// 		t.Fatalf("Error when calling GetCloseInBetween: %v", err)
// 	}

// 	// Always print the results for examination
// 	fmt.Println("Test Results:")
// 	for _, closePrice := range result.ClosePricesUSD {
// 		fmt.Printf("Date: %s\n", closePrice.Date)
// 		fmt.Printf("Close Price USD: %v\n", closePrice.ClosePriceUSD)
// 		fmt.Printf("Fetched Date: %s\n", closePrice.Metadata.FetchedDate)
// 		fmt.Printf("Conversion Rate: %v\n", closePrice.Metadata.ConversionRate)
// 		fmt.Printf("Conversion Rate Date: %s\n", closePrice.Metadata.ConversionRateDate)
// 		fmt.Printf("Candle: %s\n", closePrice.Metadata.Candle)
// 	}
// }

// func TestGetCloseUSDJSON(t *testing.T) {
// 	assetClass := "crypto"
// 	internalSymbol := "ADA_USDT"
// 	date, _ := time.Parse(time.RFC3339, "2024-06-02T00:59:59Z")

// 	// Mock the response from GetCloseUSD (you would normally use a mocking framework)
// 	expectedJSON := `{"closePriceUSD":0.4521,"metadata":{"fetchedDate":"2024-06-02T00:59:59Z","conversionRate":1,"conversionRateDate":"2024-06-02T00:59:59Z","candle":"1d"}}`

// 	// Assuming GetCloseUSD would be replaced with a mock if using a more complex environment
// 	jsonResult, err := GetCloseUSDJSON(assetClass, internalSymbol, date)
// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}
// 	if jsonResult != expectedJSON {
// 		t.Errorf("Expected JSON response did not match.\nExpected: %s\nGot: %s", expectedJSON, jsonResult)
// 	}
// }

func TestGetCloseInBetweenJSON(t *testing.T) {
	assetClass := "crypto"
	internalSymbol := "ADA_USDT"
	startDate, _ := time.Parse(time.RFC3339, "2024-06-01T00:00:00Z")
	endDate, _ := time.Parse(time.RFC3339, "2024-06-03T00:00:00Z")

	// Mock the response from GetCloseInBetween
	expectedJSON := `{"closePricesUSD":[{"date":"2024-06-01T00:00:00Z","closePriceUSD":0.45,"metadata":{"fetchedDate":"2024-06-01T00:05:00Z","conversionRate":1,"conversionRateDate":"2024-06-01T00:00:00Z","candle":"1h"}},{"date":"2024-06-03T00:00:00Z","closePriceUSD":0.47,"metadata":{"fetchedDate":"2024-06-03T00:05:00Z","conversionRate":1,"conversionRateDate":"2024-06-03T00:00:00Z","candle":"1h"}}]}`

	jsonResult, err := GetCloseInBetweenJSON(assetClass, internalSymbol, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if jsonResult != expectedJSON {
		t.Errorf("Expected JSON response did not match.\nExpected: %s\nGot: %s", expectedJSON, jsonResult)
	}
}
