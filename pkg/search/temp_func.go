package search

// func GetCloseInBetweenIndex_temp(assetClass, internalSymbol, startDate, endDate string) (CloseRangeResult, error) {
// 	// Parse the start date string into a time.Time object.
// 	start, err := time.Parse(time.RFC3339, startDate)
// 	if err != nil {
// 		// Return an error if the start date format is invalid.
// 		return CloseRangeResult{}, errors.New("invalid start date format")
// 	}

// 	// Parse the end date string into a time.Time object.
// 	end, err := time.Parse(time.RFC3339, endDate)
// 	if err != nil {
// 		// Return an error if the end date format is invalid.
// 		return CloseRangeResult{}, errors.New("invalid end date format")
// 	}

// 	// Query the index for data around the start date.
// 	startResults, err := queryIndex(index, startDate)
// 	if err != nil {
// 		// Return an error if the index query fails for the start date.
// 		return CloseRangeResult{}, fmt.Errorf("query index error for start date: %v", err)
// 	}

// 	// Query the index for data around the end date.
// 	endResults, err := queryIndex(index, endDate)
// 	if err != nil {
// 		// Return an error if the index query fails for the end date.
// 		return CloseRangeResult{}, fmt.Errorf("query index error for end date: %v", err)
// 	}

// 	// Assume startResults[0] and endResults[0] are the closest matches for start and end dates.
// 	startData := startResults[0]
// 	endData := endResults[0]

// 	// Parse the start and end close prices from the string data retrieved.
// 	startClosePrice, _ := strconv.ParseFloat(startData["Close"], 64)
// 	endClosePrice, _ := strconv.ParseFloat(endData["Close"], 64)

// 	// Parse the closest start and end dates from the string data retrieved.
// 	startClosestDate, _ := time.Parse(time.RFC3339, startData["Date"])
// 	endClosestDate, _ := time.Parse(time.RFC3339, endData["Date"])

// 	// Extract the base currency from the symbol name.
// 	baseCurrency := extractBaseCurrency(internalSymbol)

// 	// Retrieve conversion rates for the start and end close prices.
// 	startConversionRate, startConversionRateDate, err := getConversionRateForCloseInBetween(baseCurrency, startClosestDate, startClosePrice)
// 	if err != nil {
// 		// Return an error if retrieving the start conversion rate fails.
// 		return CloseRangeResult{}, fmt.Errorf("failed to retrieve start conversion rate: %v", err)
// 	}
// 	endConversionRate, endConversionRateDate, err := getConversionRateForCloseInBetween(baseCurrency, endClosestDate, endClosePrice)
// 	if err != nil {
// 		// Return an error if retrieving the end conversion rate fails.
// 		return CloseRangeResult{}, fmt.Errorf("failed to retrieve end conversion rate: %v", err)
// 	}

// 	// Calculate the USD prices for start and end close prices.
// 	startClosePriceUSD := startClosePrice * startConversionRate
// 	endClosePriceUSD := endClosePrice * endConversionRate

// 	// Return a struct populated with the calculated data for the range.
// 	return CloseRangeResult{
// 		StartClosePriceUSD:      startClosePriceUSD,
// 		EndClosePriceUSD:        endClosePriceUSD,
// 		StartFetchedDate:        startClosestDate.Format(time.RFC3339),
// 		EndFetchedDate:          endClosestDate.Format(time.RFC3339),
// 		StartConversionRate:     startConversionRate,
// 		StartConversionRateDate: startConversionRateDate.Format(time.RFC3339),
// 		EndConversionRate:       endConversionRate,
// 		EndConversionRateDate:   endConversionRateDate.Format(time.RFC3339),
// 		Candle:                  "1d",
// 	}, nil
// }
