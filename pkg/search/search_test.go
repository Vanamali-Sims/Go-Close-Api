package search

import (
	"fmt"
	"testing"
)

func TestGetCloseUSD(t *testing.T) {
	assetClass := "crypto"
	internalSymbol := "1000SATS_USDT"
	date := "2024-01-22T00:00:59Z"
	result, err := GetCloseUSD(assetClass, internalSymbol, date)
	if err != nil {
		t.Fatalf("Error when calling GetCloseUSD: %v", err)
	}

	fmt.Printf("Raw Close Price: %f\n", result.RawClosePrice)
	fmt.Printf("Converted Close Price USD: %f\n", result.ClosePriceUSD)
	fmt.Printf("Conversion Rate: %f on %s\n", result.ConversionRate, result.ConversionRateDate)
	fmt.Printf("Data Fetched on: %s\n", result.FetchedDate)

	if result.ClosePriceUSD == 0 {
		t.Errorf("Expected non-zero ClosePriceUSD, got %f", result.ClosePriceUSD)
	}
}
