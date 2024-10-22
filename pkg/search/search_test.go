package search

import (
	"fmt"
	"testing"
	"time"
)

func TestGetCloseUSDForSpecificDate(t *testing.T) {
	assetClass := "stocks"
	internalSymbol := "1120__XSAU_SAR"
	dateString := "2024-07-09T09:30:00Z"
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		t.Fatalf("Error parsing date: %v", err)
	}

	// Run the GetCloseUSD function
	result, err := GetCloseUSD(assetClass, internalSymbol, date)
	if err != nil {
		t.Fatalf("Error when calling GetCloseUSD: %v", err)
	}

	// Print the result for manual verification
	fmt.Printf("Result: %+v\n", result)
}
