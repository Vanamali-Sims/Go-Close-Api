package search

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type CloseResult struct {
	ClosePriceUSD      float64 `json:"closePriceUSD"`
	FetchedDate        string  `json:"fetchedDate"`
	ConversionRate     float64 `json:"conversionRate"`
	ConversionRateDate string  `json:"conversionRateDate"`
	Candle             string  `json:"candle"`
	RawClosePrice      float64 `json:"rawClosePrice"`
}

type CloseRangeResult struct {
	StartClosePriceUSD float64 `json:"startClosePriceUSD"`
	EndClosePriceUSD   float64 `json:"endClosePriceUSD"`
	StartFetchedDate   string  `json:"startFetchedDate"`
	EndFetchedDate     string  `json:"endFetchedDate"`
	ConversionRate     float64 `json:"conversionRate"`
	ConversionRateDate string  `json:"conversionRateDate"`
	Candle             string  `json:"candle"`
}

func GetCloseUSD(assetClass, internalSymbol, date string) (CloseResult, error) {
	baseCurrency := extractBaseCurrency(internalSymbol) // Extract base currency from symbol
	requestedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return CloseResult{}, errors.New("invalid date format")
	}

	// Fetch conversion rate
	conversionRate, conversionRateDate, err := getConversionRate(baseCurrency, requestedDate)
	if err != nil {
		return CloseResult{}, fmt.Errorf("conversion rate error: %v", err)
	}

	// Here would be the logic to fetch the raw close price from CSV data (not shown here)
	rawClosePrice := 0.522 // Example raw close price fetched from data
	closePriceUSD := rawClosePrice * conversionRate

	return CloseResult{
		ClosePriceUSD:      closePriceUSD,
		RawClosePrice:      rawClosePrice,
		FetchedDate:        requestedDate.Format(time.RFC3339),
		ConversionRate:     conversionRate,
		ConversionRateDate: conversionRateDate.Format(time.RFC3339),
		Candle:             "1d",
	}, nil
}

// GetCloseInBetween fetches the closing prices between two dates for an asset.
func GetCloseInBetween(assetClass, internalSymbol, startDate, endDate string) (CloseRangeResult, error) {
	// Parse the requested dates
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return CloseRangeResult{}, errors.New("invalid start date format")
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return CloseRangeResult{}, errors.New("invalid end date format")
	}

	// Define the CSV file paths for the asset class (crypto, stocks)
	csvFilePath := filepath.Join("data", assetClass, internalSymbol+".csv")

	// Forex file path to be used for conversion rates
	forexFilePath := filepath.Join("data", "forex", "conversion_rates.csv")

	// Read the CSV file for the asset (crypto/stocks)
	assetData, err := readCSV(csvFilePath)
	if err != nil {
		return CloseRangeResult{}, fmt.Errorf("failed to read asset CSV: %v", err)
	}

	// Find the closest matching dates and retrieve the closing prices
	startClosePrice, startClosestDate, err := findClosestDate(assetData, start)
	if err != nil {
		return CloseRangeResult{}, fmt.Errorf("error finding closest start date: %v", err)
	}

	endClosePrice, endClosestDate, err := findClosestDate(assetData, end)
	if err != nil {
		return CloseRangeResult{}, fmt.Errorf("error finding closest end date: %v", err)
	}

	// Assuming the asset is priced in local currency, perform conversion to USD if necessary
	conversionRate, conversionRateDate, err := getConversionRate(forexFilePath, start)
	if err != nil {
		return CloseRangeResult{}, fmt.Errorf("failed to retrieve conversion rate: %v", err)
	}

	// Calculate the USD price for both start and end
	startClosePriceUSD := startClosePrice * conversionRate
	endClosePriceUSD := endClosePrice * conversionRate

	return CloseRangeResult{
		StartClosePriceUSD: startClosePriceUSD,
		EndClosePriceUSD:   endClosePriceUSD,
		StartFetchedDate:   startClosestDate.Format(time.RFC3339),
		EndFetchedDate:     endClosestDate.Format(time.RFC3339),
		ConversionRate:     conversionRate,
		ConversionRateDate: conversionRateDate.Format(time.RFC3339),
		Candle:             "1d", // Assume daily candles for this example
	}, nil
}

// findClosestDate finds the row in the CSV with the closest date to the requested date.
func findClosestDate(data []map[string]string, requestedDate time.Time) (float64, time.Time, error) {
	var closestDate time.Time
	var closePrice float64
	smallestDiff := time.Duration(1<<63 - 1) // Large initial value

	for _, row := range data {
		recordDate, err := time.Parse(time.RFC3339, row["Date"])
		if err != nil {
			continue // skip rows with invalid dates
		}

		diff := requestedDate.Sub(recordDate).Abs()
		if diff < smallestDiff {
			smallestDiff = diff
			closestDate = recordDate
			closePrice, _ = strconv.ParseFloat(row["Close"], 64)
		}
	}

	if smallestDiff == time.Duration(1<<63-1) {
		return 0, time.Time{}, errors.New("no matching date found")
	}

	return closePrice, closestDate, nil
}

func extractBaseCurrency(symbol string) string {
	parts := strings.Split(symbol, "_")
	if len(parts) > 1 {
		if parts[1] == "USDT" { // Treat USDT as USD directly
			return "USD"
		}
		return parts[1] // Return the currency part
	}
	return "USD" // Default to USD if no currency part is found
}

// getConversionRate retrieves the conversion rate from a forex CSV file.
func getConversionRate(baseCurrency string, date time.Time) (float64, time.Time, error) {
	// Always targeting USD
	targetCurrency := "USD"

	// If the base currency is USD, skip file lookup and return 1.0 directly
	if baseCurrency == "USD" || baseCurrency == "USDT" {
		return 1.0, date, nil // No conversion needed, rate is 1:1, return the requested date as the rate date
	}

	// Construct the file path based on the currency pair
	filePath := filepath.Join("data", "forex", fmt.Sprintf("%s_%s.csv", baseCurrency, targetCurrency))

	data, err := readCSV(filePath)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to open forex file: %v", err)
	}

	var closestDate time.Time
	var conversionRate float64
	smallestDiff := time.Duration(1<<63 - 1) // Initialize with a large value

	for _, row := range data {
		recordDate, err := time.Parse(time.RFC3339, row["Date"])
		if err != nil {
			continue // Skip rows with invalid date formats
		}
		diff := date.Sub(recordDate).Abs()
		if diff < smallestDiff {
			smallestDiff = diff
			closestDate = recordDate
			conversionRate, err = strconv.ParseFloat(row["ConversionRate"], 64)
			if err != nil {
				continue // Skip rows with invalid conversion rates
			}
		}
	}

	if smallestDiff == time.Duration(1<<63-1) {
		return 0, time.Time{}, errors.New("no conversion rate found")
	}

	return conversionRate, closestDate, nil
}
