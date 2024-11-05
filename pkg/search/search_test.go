package search

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/blevesearch/bleve"
)

func TestMain(m *testing.M) {
	// Set up the index before running tests
	basePath := "C:\\Users\\isvan\\OneDrive\\Documents\\work\\GoApi\\data"
	indexPath := filepath.Join(basePath, "index", "search.bleve")
	index, err := initializeIndex(indexPath)
	if err != nil {
		log.Fatalf("Failed to initialize or open index: %v", err)
	}

	// Check if indexing is needed or not
	if needsIndexing(index) {
		if err := indexAllCSVs(basePath); err != nil {
			log.Fatalf("Failed to index CSV files: %v", err)
		}
	}

	// Run tests
	exitVal := m.Run()

	// Optionally clean up after tests
	// os.RemoveAll(indexPath) // Be careful with this in real environments!

	os.Exit(exitVal)
}

func needsIndexing(index bleve.Index) bool {
	// Implement logic to decide if indexing is needed
	return true // For now, always true for simplicity
}

func TestGetCloseUSDIndex(t *testing.T) {
	date, _ := time.Parse(time.RFC3339, "2024-06-01T00:59:59Z")
	result, err := GetCloseUSDIndex("crypto", "BTC_USD", date)
	if err != nil {
		t.Errorf("Error retrieving close USD index: %v", err)
		return
	}

	if result.ClosePriceUSD <= 0 {
		t.Errorf("Expected positive close price USD, got %v", result.ClosePriceUSD)
	}
}

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
