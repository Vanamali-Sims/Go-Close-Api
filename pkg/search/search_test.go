package search

import (
	"fmt"
	"testing"
)

// func TestGetCloseUSDForwSpecificDate(t *testing.T) {
// 	assetClass := "stocks"
// 	internalSymbol := "AAJ__XASX_AUD"
// 	dateString := "2024-07-11T23:59:59Z"
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

func TestGetCloseInBetweenForSpecificDateRange(t *testing.T) {
	assetClass := "crypto"
	internalSymbol := "1000SATS_USDT"
	startDateString := "2024-07-02T15:59:59Z"
	endDateString := "2024-07-02T19:59:59Z"

	// Run the GetCloseInBetween function
	result, err := GetCloseInBetween(assetClass, internalSymbol, startDateString, endDateString)
	if err != nil {
		t.Fatalf("Error when calling GetCloseInBetween: %v", err)
	}

	// Print the result for manual verification
	fmt.Printf("Result: %+v\n", result)
}
