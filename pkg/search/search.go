package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
)

var index bleve.Index

type Metadata struct {
	FetchedDate        string  `json:"fetchedDate"`
	ConversionRate     float64 `json:"conversionRate"`
	ConversionRateDate string  `json:"conversionRateDate"`
	Candle             string  `json:"candle"`
}

type CloseUSDResponse struct {
	ClosePriceUSD float64  `json:"closePriceUSD"`
	Metadata      Metadata `json:"metadata"`
}

type ClosePriceDetail struct {
	Date          string   `json:"date"`
	ClosePriceUSD float64  `json:"closePriceUSD"`
	Metadata      Metadata `json:"metadata"`
}

type CloseInBetweenResponse struct {
	ClosePricesUSD []ClosePriceDetail `json:"closePricesUSD"`
}

type CloseResult struct {
	ClosePriceUSD      float64 `json:"closePriceUSD"`
	FetchedDate        string  `json:"fetchedDate"`
	ConversionRate     float64 `json:"conversionRate"`
	ConversionRateDate string  `json:"conversionRateDate"`
	Candle             string  `json:"candle"`
	RawClosePrice      float64 `json:"rawClosePrice"`
}

type CloseRangeResult struct {
	StartClosePriceUSD      float64 `json:"startClosePriceUSD"`
	EndClosePriceUSD        float64 `json:"endClosePriceUSD"`
	StartFetchedDate        string  `json:"startFetchedDate"`
	EndFetchedDate          string  `json:"endFetchedDate"`
	StartConversionRate     float64 `json:"startConversionRate"`
	StartConversionRateDate string  `json:"startConversionRateDate"`
	EndConversionRate       float64 `json:"endConversionRate"`
	EndConversionRateDate   string  `json:"endConversionRateDate"`
	Candle                  string  `json:"candle"`
}

func GetCloseUSDJSON(assetClass, internalSymbol string, date time.Time) (string, error) {
	result, err := GetCloseUSD(assetClass, internalSymbol, date)
	if err != nil {
		return "", fmt.Errorf("failed to get close USD data: %v", err)
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal results to JSON: %v", err)
	}

	return string(jsonData), nil
}

func GetCloseInBetweenJSON(assetClass, internalSymbol, startDate, endDate string) (string, error) {
	results, err := GetCloseInBetween(assetClass, internalSymbol, startDate, endDate)
	if err != nil {
		return "", fmt.Errorf("failed to get close price data in between: %v", err)
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("failed to marshal results to JSON: %v", err)
	}

	return string(jsonData), nil
}

func GetCloseUSD(assetClass, internalSymbol string, date time.Time) (CloseUSDResponse, error) {
	csvFilePath, err := findDataPath(assetClass, internalSymbol, date)
	if err != nil {
		return CloseUSDResponse{}, fmt.Errorf("failed to find CSV file path: %v", err)
	}

	assetData, err := readCSV(csvFilePath)
	if err != nil {
		return CloseUSDResponse{}, fmt.Errorf("failed to read asset CSV: %v", err)
	}

	rawClosePrice, closestDate, err := findClosestDate(assetData, date)
	if err != nil {
		return CloseUSDResponse{}, fmt.Errorf("error finding closest date: %v", err)
	}

	baseCurrency := extractBaseCurrency(internalSymbol)
	closePriceUSD, conversionRate, conversionRateDate, err := getConversionRate(baseCurrency, date, rawClosePrice)
	if err != nil {
		return CloseUSDResponse{}, fmt.Errorf("conversion rate error: %v", err)
	}

	metadata := Metadata{
		FetchedDate:        closestDate.Format(time.RFC3339),
		ConversionRate:     conversionRate,
		ConversionRateDate: conversionRateDate.Format(time.RFC3339),
		Candle:             "1d", // Adjust as necessary
	}

	return CloseUSDResponse{
		ClosePriceUSD: closePriceUSD,
		Metadata:      metadata,
	}, nil
}

func GetCloseInBetween(assetClass, internalSymbol, startDate, endDate string) (CloseInBetweenResponse, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return CloseInBetweenResponse{}, fmt.Errorf("invalid start date format: %v", err)
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return CloseInBetweenResponse{}, fmt.Errorf("invalid end date format: %v", err)
	}

	csvFilePath, err := findDataPath(assetClass, internalSymbol, start)
	if err != nil {
		return CloseInBetweenResponse{}, fmt.Errorf("failed to find CSV file path: %v", err)
	}

	assetData, err := readCSV(csvFilePath)
	if err != nil {
		return CloseInBetweenResponse{}, fmt.Errorf("failed to read asset CSV: %v", err)
	}

	startClosePrice, startClosestDate, err := findClosestDate(assetData, start)
	if err != nil {
		return CloseInBetweenResponse{}, fmt.Errorf("error finding closest start date: %v", err)
	}

	endClosePrice, endClosestDate, err := findClosestDate(assetData, end)
	if err != nil {
		return CloseInBetweenResponse{}, fmt.Errorf("error finding closest end date: %v", err)
	}

	baseCurrency := extractBaseCurrency(internalSymbol)

	startConversionRate, startConversionRateDate, err := getConversionRateForCloseInBetween(baseCurrency, startClosestDate, startClosePrice)
	if err != nil {
		return CloseInBetweenResponse{}, fmt.Errorf("failed to retrieve start conversion rate: %v", err)
	}

	endConversionRate, endConversionRateDate, err := getConversionRateForCloseInBetween(baseCurrency, endClosestDate, endClosePrice)
	if err != nil {
		return CloseInBetweenResponse{}, fmt.Errorf("failed to retrieve end conversion rate: %v", err)
	}

	startDetail := ClosePriceDetail{
		Date:          startClosestDate.Format(time.RFC3339),
		ClosePriceUSD: startClosePrice * startConversionRate,
		Metadata: Metadata{
			FetchedDate:        startClosestDate.Format(time.RFC3339),
			ConversionRate:     startConversionRate,
			ConversionRateDate: startConversionRateDate.Format(time.RFC3339),
			Candle:             "1d", // Adjust as necessary
		},
	}

	endDetail := ClosePriceDetail{
		Date:          endClosestDate.Format(time.RFC3339),
		ClosePriceUSD: endClosePrice * endConversionRate,
		Metadata: Metadata{
			FetchedDate:        endClosestDate.Format(time.RFC3339),
			ConversionRate:     endConversionRate,
			ConversionRateDate: endConversionRateDate.Format(time.RFC3339),
			Candle:             "1d", // Adjust as necessary
		},
	}

	return CloseInBetweenResponse{
		ClosePricesUSD: []ClosePriceDetail{startDetail, endDetail},
	}, nil
}

func jsonResponse(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error Marshalling data: %v", err)
		return "{}"
	}
	return string(jsonData)
}

func extractBaseCurrency(filename string) string {
	// Remove the extension first
	withoutExtension := strings.TrimSuffix(filename, ".csv")

	// Split on the last underscore to hopefully get the currency as the last part
	parts := strings.Split(withoutExtension, "_")
	if len(parts) > 1 {
		baseCurrency := parts[len(parts)-1] // Get the last part after the last underscore
		if baseCurrency == "USDT" {
			baseCurrency = "USD" // Normalize USDT to USD
		}
		return baseCurrency
	}
	return "USD" // Default to USD if parsing fails or no underscore is found
}

func getConversionRateForCloseInBetween(baseCurrency string, date time.Time, closePrice float64) (float64, time.Time, error) {

	if baseCurrency == "USD" || baseCurrency == "USDT" {
		return 1.0, date, nil
	}

	conversionRate, conversionRateDateFloat, _, err := getConversionRate(baseCurrency, date, closePrice)
	if err != nil {
		return 0, time.Time{}, err
	}

	// Convert conversionRateDateFloat from float64 to time.Time
	conversionRateDate := time.Unix(int64(conversionRateDateFloat), 0)

	return conversionRate, conversionRateDate, nil
}

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

func getConversionRate(baseCurrency string, date time.Time, rawClosePrice float64) (float64, float64, time.Time, error) {
	if baseCurrency == "USD" || baseCurrency == "USDT" {
		return rawClosePrice, 1.0, date, nil // Directly return for USD as no conversion is needed
	}

	// Find the path to the forex data file
	forexFilePath, err := findForexPath(baseCurrency, date)
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("failed to find forex file path: %v", err)
	}

	// Read the forex data from the CSV file found
	data, err := readCSV(forexFilePath)
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("failed to open forex file: %v", err)
	}

	// Find the closest conversion rate in the forex data
	closestDate, conversionRate, err := findClosestConversionRate(data, date)
	if err != nil {
		return 0, 0, time.Time{}, err
	}

	// Calculate the close price in USD
	closePriceUSD := rawClosePrice * conversionRate
	return closePriceUSD, conversionRate, closestDate, nil
}

func findClosestConversionRate(data []map[string]string, targetDate time.Time) (time.Time, float64, error) {
	var closestDate time.Time
	var conversionRate float64
	smallestDiff := time.Duration(1<<63 - 1)

	for _, row := range data {
		recordDate, err := time.Parse(time.RFC3339, row["Date"])
		if err != nil {
			continue
		}
		rate, err := strconv.ParseFloat(row["Close"], 64) // Changed from "ConversionRate" to "Close"
		if err != nil {
			fmt.Printf("Skipping row with invalid rate: '%v' at %v\n", row["Close"], row["Date"])
			continue
		}

		diff := targetDate.Sub(recordDate).Abs()
		if diff < smallestDiff {
			smallestDiff = diff
			closestDate = recordDate
			conversionRate = rate
		}
	}

	if smallestDiff == time.Duration(1<<63-1) {
		return time.Time{}, 0, errors.New("no conversion rate found")
	}

	return closestDate, conversionRate, nil
}

func findDataPath(assetClass, internalSymbol string, date time.Time) (string, error) {
	basePath := filepath.Join("C:\\Users\\isvan\\OneDrive\\Documents\\work\\GoApi\\data", assetClass)
	year := date.Format("2006") // Ensure four-digit year
	month := date.Format("01")
	day := date.Format("02")

	intervals := []string{"1m", "2m", "5m", "15m", "1h", "1w", "1d"}
	for _, interval := range intervals {
		path := filepath.Join(basePath, year, month, day, interval, fmt.Sprintf("%s.csv", internalSymbol))
		fmt.Println("Checking path:", path)
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Found file at:", path)
			return path, nil
		} else {
			fmt.Println("Could not find file at:", path, "; Error:", err)
		}
	}

	allPath := filepath.Join(basePath, "all", fmt.Sprintf("%s.csv", internalSymbol))
	fmt.Println("Checking all directory path:", allPath)
	if _, err := os.Stat(allPath); err == nil {
		fmt.Println("Found file in all directory at:", allPath)
		return allPath, nil
	} else {
		fmt.Println("Could not find file in all directory;", "Error:", err)
	}

	return "", fmt.Errorf("no valid data path found for the date: %s", date)
}

func findForexPath(baseCurrency string, date time.Time) (string, error) {
	basePath := "C:\\Users\\isvan\\OneDrive\\Documents\\work\\GoApi\\data\\forex"
	year := date.Format("2006")
	month := date.Format("01")
	day := date.Format("02")

	// Target currency is always USD
	targetCurrency := "USD"
	fileName := fmt.Sprintf("%s_%s.csv", baseCurrency, targetCurrency)

	// Intervals to check in order of priority
	intervals := []string{"1m", "2m", "5m", "15m", "1h", "1w", "1d"}
	for _, interval := range intervals {
		path := filepath.Join(basePath, year, month, day, interval, fileName)
		fmt.Println("Checking forex path:", path)
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Found forex file at:", path)
			return path, nil
		} else {
			fmt.Println("Could not find forex file at:", path, "; Error:", err)
		}
	}

	// Fallback to the 'all' directory at the year level if no specific interval file is found
	allPath := filepath.Join(basePath, "all", fileName)
	fmt.Println("Checking forex all directory path:", allPath)
	if _, err := os.Stat(allPath); err == nil {
		fmt.Println("Found forex file in all directory at:", allPath)
		return allPath, nil
	} else {
		fmt.Println("Could not find forex file in all directory;", "Error:", err)
	}

	return "", fmt.Errorf("no valid forex data path found for the date: %s", date)
}

func GetCloseUSDIndex(assetClass, internalSymbol string, date time.Time) (CloseResult, error) {
	// Format the date to a string as expected by the query function (ISO8601/RFC3339 format).
	dateQuery := date.Format(time.RFC3339)
	// Query the index for data that matches the specified date.
	results, err := queryIndex(index, dateQuery)
	if err != nil {
		// Return an error if the index query fails.
		return CloseResult{}, fmt.Errorf("query index error: %v", err)
	}

	if len(results) == 0 {
		return CloseResult{}, fmt.Errorf("no results found for the date: %s", dateQuery)
	}

	// Assume results[0] is the closest match. This is a simplification and may need better handling.
	data := results[0]
	// Parse the close price from the string data retrieved.
	rawClosePrice, _ := strconv.ParseFloat(data["Close"], 64)
	// Parse the date from the string data retrieved.
	closestDate, _ := time.Parse(time.RFC3339, data["Date"])
	// Extract the base currency from the symbol name.
	baseCurrency := extractBaseCurrency(internalSymbol)

	// Calculate the USD close price using the conversion rate.
	closePriceUSD, conversionRate, conversionRateDate, err := getConversionRate(baseCurrency, closestDate, rawClosePrice)
	if err != nil {
		// Return an error if there is a problem fetching the conversion rate.
		return CloseResult{}, fmt.Errorf("conversion rate error: %v", err)
	}

	// Return a struct populated with the calculated data.
	return CloseResult{
		ClosePriceUSD:      closePriceUSD,
		RawClosePrice:      rawClosePrice,
		FetchedDate:        closestDate.Format(time.RFC3339),
		ConversionRate:     conversionRate,
		ConversionRateDate: conversionRateDate.Format(time.RFC3339),
		Candle:             "1d",
	}, nil
}

func GetCloseInBetweenIndex(assetClass, internalSymbol, startDate, endDate string) (CloseRangeResult, error) {
	// Parse the start and end date strings into time.Time objects.
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return CloseRangeResult{}, fmt.Errorf("invalid start date format: %v", err)
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return CloseRangeResult{}, fmt.Errorf("invalid end date format: %v", err)
	}

	// Query the index for data matching the start date to the end date.
	queryStr := fmt.Sprintf("+Date:[%s TO %s]", start.Format(time.RFC3339), end.Format(time.RFC3339))
	results, err := queryIndex(index, queryStr)
	if err != nil {
		return CloseRangeResult{}, err // Pass the error up
	}

	if len(results) == 0 {
		return CloseRangeResult{}, errors.New("no results found for the provided date range")
	}

	// Assume startResults[0] and endResults[len(results)-1] are the closest matches for start and end dates.
	startData := results[0]
	endData := results[len(results)-1]

	// Parse the start and end close prices.
	startClosePrice, err := strconv.ParseFloat(startData["Close"], 64)
	if err != nil {
		return CloseRangeResult{}, fmt.Errorf("error parsing start close price: %v", err)
	}
	endClosePrice, err := strconv.ParseFloat(endData["Close"], 64)
	if err != nil {
		return CloseRangeResult{}, fmt.Errorf("error parsing end close price: %v", err)
	}

	// Parse the closest start and end dates.
	startClosestDate, _ := time.Parse(time.RFC3339, startData["Date"])
	endClosestDate, _ := time.Parse(time.RFC3339, endData["Date"])

	// Extract the base currency from the symbol name.
	baseCurrency := extractBaseCurrency(internalSymbol)

	// Retrieve conversion rates for the start and end close prices.
	startConversionRate, startConversionRateDate, err := getConversionRateForCloseInBetween(baseCurrency, startClosestDate, startClosePrice)
	if err != nil {
		return CloseRangeResult{}, fmt.Errorf("failed to retrieve start conversion rate: %v", err)
	}

	endConversionRate, endConversionRateDate, err := getConversionRateForCloseInBetween(baseCurrency, endClosestDate, endClosePrice)
	if err != nil {
		return CloseRangeResult{}, fmt.Errorf("failed to retrieve end conversion rate: %v", err)
	}

	// Calculate the USD prices for start and end close prices.
	startClosePriceUSD := startClosePrice * startConversionRate
	endClosePriceUSD := endClosePrice * endConversionRate

	// Return a struct populated with the calculated data for the range.
	return CloseRangeResult{
		StartClosePriceUSD:      startClosePriceUSD,
		EndClosePriceUSD:        endClosePriceUSD,
		StartFetchedDate:        startClosestDate.Format(time.RFC3339),
		EndFetchedDate:          endClosestDate.Format(time.RFC3339),
		StartConversionRate:     startConversionRate,
		StartConversionRateDate: startConversionRateDate.Format(time.RFC3339),
		EndConversionRate:       endConversionRate,
		EndConversionRateDate:   endConversionRateDate.Format(time.RFC3339),
		Candle:                  "1d",
	}, nil
}
