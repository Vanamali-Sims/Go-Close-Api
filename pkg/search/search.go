package search

import (
	"errors"
	"fmt"
	"os"
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

// Corrected GetCloseUSD function
func GetCloseUSD(assetClass, internalSymbol string, date time.Time) (CloseResult, error) {
	// Path finding using the adjusted findDataPath function
	csvFilePath, err := findDataPath(assetClass, internalSymbol, date)
	if err != nil {
		return CloseResult{}, fmt.Errorf("failed to find CSV file path: %v", err)
	}
	fmt.Println("CSV data path:", csvFilePath) // Debugging output to verify path correctness

	// Reading the CSV file using the found path
	assetData, err := readCSV(csvFilePath)
	if err != nil {
		return CloseResult{}, fmt.Errorf("failed to read asset CSV: %v", err)
	}

	// Finding the closest matching date and retrieving the raw close price
	rawClosePrice, closestDate, err := findClosestDate(assetData, date)
	if err != nil {
		return CloseResult{}, fmt.Errorf("error finding closest date: %v", err)
	}

	// Fetching the conversion rate using a dynamically constructed forex file path
	forexFilePath := filepath.Join("data", "forex", "all", fmt.Sprintf("%s_%s.csv", extractBaseCurrency(internalSymbol), "USD"))
	conversionRate, conversionRateDate, err := getConversionRate(forexFilePath, date)
	if err != nil {
		return CloseResult{}, fmt.Errorf("conversion rate error: %v", err)
	}

	// Calculating the USD price
	closePriceUSD := rawClosePrice * conversionRate

	// Returning the result struct with all relevant data
	return CloseResult{
		ClosePriceUSD:      closePriceUSD,
		RawClosePrice:      rawClosePrice,
		FetchedDate:        closestDate.Format(time.RFC3339),
		ConversionRate:     conversionRate,
		ConversionRateDate: conversionRateDate.Format(time.RFC3339),
		Candle:             "1d",
	}, nil
}

// GetCloseInBetween fetches the closing prices between two dates for an asset.
func GetCloseInBetween(assetClass, internalSymbol, startDate, endDate string) (CloseRangeResult, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return CloseRangeResult{}, errors.New("invalid start date format")
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return CloseRangeResult{}, errors.New("invalid end date format")
	}

	// Define the CSV file path
	csvFilePath := filepath.Join("data", assetClass, internalSymbol+".csv")

	// Read the CSV file for the asset
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
	baseCurrency := extractBaseCurrency(internalSymbol)
	conversionRate, conversionRateDate, err := getConversionRate(baseCurrency, start)
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

// func getConversionRate(forexFilePath string, date time.Time) (float64, time.Time, error) {
// 	// If the base currency is USD, no need to fetch conversion rates as it is 1:1
// 	if strings.Contains(forexFilePath, "USD_USD.csv") {
// 		return 1.0, date, nil
// 	}

// 	// Load the CSV data from the forex file path
// 	data, err := readCSV(forexFilePath)
// 	if err != nil {
// 		return 0, time.Time{}, fmt.Errorf("failed to open forex file: %v", err)
// 	}

// 	var closestDate time.Time
// 	var conversionRate float64
// 	smallestDiff := time.Duration(1<<63 - 1) // Initialize with a large value

// 	// Iterate through the data to find the closest date and its conversion rate
// 	for _, row := range data {
// 		recordDate, err := time.Parse(time.RFC3339, row["Date"])
// 		if err != nil {
// 			continue // Skip rows with invalid date formats
// 		}
// 		diff := date.Sub(recordDate).Abs()
// 		if diff < smallestDiff {
// 			smallestDiff = diff
// 			closestDate = recordDate
// 			conversionRate, err = strconv.ParseFloat(row["ConversionRate"], 64)
// 			if err != nil {
// 				continue // Skip rows with invalid conversion rates
// 			}
// 		}
// 	}

// 	if smallestDiff == time.Duration(1<<63-1) {
// 		return 0, time.Time{}, errors.New("no conversion rate found")
// 	}

// 	return conversionRate, closestDate, nil
// }

func getConversionRate(baseCurrency string, date time.Time) (float64, time.Time, error) {
	if baseCurrency == "USD" {
		return 1.0, date, nil // No conversion needed for USD to USD
	}

	forexFilePath, err := findForexPath(baseCurrency, date)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to find forex file path: %v", err)
	}

	fmt.Println("Using forex data path:", forexFilePath)

	data, err := readCSV(forexFilePath)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to open forex file: %v", err)
	}

	var closestDate time.Time
	var conversionRate float64
	smallestDiff := time.Duration(1<<63 - 1)

	for _, row := range data {
		recordDate, err := time.Parse(time.RFC3339, row["Date"])
		if err != nil {
			continue
		}
		diff := date.Sub(recordDate).Abs()
		if diff < smallestDiff {
			smallestDiff = diff
			closestDate = recordDate
			conversionRate, err = strconv.ParseFloat(row["ConversionRate"], 64)
			if err != nil {
				continue
			}
		}
	}

	if smallestDiff == time.Duration(1<<63-1) {
		return 0, time.Time{}, errors.New("no conversion rate found")
	}

	return conversionRate, closestDate, nil
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

	allPath := filepath.Join(basePath, year, "all", fmt.Sprintf("%s.csv", internalSymbol))
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
	allPath := filepath.Join(basePath, year, "all", fileName)
	fmt.Println("Checking forex all directory path:", allPath)
	if _, err := os.Stat(allPath); err == nil {
		fmt.Println("Found forex file in all directory at:", allPath)
		return allPath, nil
	} else {
		fmt.Println("Could not find forex file in all directory;", "Error:", err)
	}

	return "", fmt.Errorf("no valid forex data path found for the date: %s", date)
}
