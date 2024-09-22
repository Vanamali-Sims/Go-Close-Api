package search

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"time"
)

type CloseResult struct {
	ClosePriceUSD      float64 `json:"closePriceUSD"`
	FetchedDate        string  `json:"fetchedDate"`
	ConversionRate     float64 `json:"conversionRate"`
	ConversionRateDate string  `json:"conversionRateDate"`
	Candle             string  `json:"candle"`
}

// GetCloseUSD searches for the closing price of an asset, converts it to USD if necessary.
func GetCloseUSD(assetClass, internalSymbol, date string) (CloseResult, error) {
	// Parse the requested date
	requestedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return CloseResult{}, errors.New("invalid date format")
	}

	// Define file paths
	csvFilePath := filepath.Join("data", assetClass, internalSymbol+".csv")
	forexFilePath := filepath.Join("data", "forex", "conversion_rates.csv")

	// Read the CSV file for the asset
	assetData, err := readCSV(csvFilePath)
	if err != nil {
		return CloseResult{}, fmt.Errorf("failed to read asset CSV: %v", err)
	}

	// Find the closest matching date and retrieve the closing price
	closePrice, closestDate, err := findClosestDate(assetData, requestedDate)
	if err != nil {
		return CloseResult{}, fmt.Errorf("error finding closest date: %v", err)
	}

	// Assuming the asset is priced in local currency, perform conversion to USD if necessary
	conversionRate, conversionRateDate, err := getConversionRate(forexFilePath, requestedDate)
	if err != nil {
		return CloseResult{}, fmt.Errorf("failed to retrieve conversion rate: %v", err)
	}

	// Calculate the USD price
	closePriceUSD := closePrice * conversionRate

	return CloseResult{
		ClosePriceUSD:      closePriceUSD,
		FetchedDate:        closestDate.Format(time.RFC3339),
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

// getConversionRate retrieves the conversion rate from a forex CSV file.
func getConversionRate(filePath string, date time.Time) (float64, time.Time, error) {
	data, err := readCSV(filePath)
	if err != nil {
		return 0, time.Time{}, err
	}

	var closestDate time.Time
	var conversionRate float64
	smallestDiff := time.Duration(1<<63 - 1) // Large initial value

	for _, row := range data {
		recordDate, err := time.Parse(time.RFC3339, row["Date"])
		if err != nil {
			continue
		}

		diff := date.Sub(recordDate).Abs()
		if diff < smallestDiff {
			smallestDiff = diff
			closestDate = recordDate
			conversionRate, _ = strconv.ParseFloat(row["USDConversionRate"], 64)
		}
	}

	if smallestDiff == time.Duration(1<<63-1) {
		return 0, time.Time{}, errors.New("no conversion rate found")
	}

	return conversionRate, closestDate, nil
}
