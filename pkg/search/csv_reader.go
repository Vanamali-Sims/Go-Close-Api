package search

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/blevesearch/bleve"
)

func readCSV(filePath string) ([]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	var data []map[string]string
	var headers []string
	firstLine := true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err // Handle non-EOF errors
		}

		if firstLine {
			headers = record
			firstLine = false
			continue
		}

		row := make(map[string]string)
		for i, value := range record {
			if i < len(headers) { // Check to prevent index out of range errors
				row[headers[i]] = value
			}
		}
		data = append(data, row)
	}
	return data, nil
}

func queryIndex(index bleve.Index, queryString string) ([]map[string]string, error) {
	// Construct a query for the Bleve search index. In this case, using a simple query syntax.
	query := bleve.NewQueryStringQuery(queryString)

	// Create a search request based on the query.
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"Date", "Close"} // Specify fields you want to retrieve.

	// Execute the search query.
	searchResults, err := index.Search(search)
	if err != nil {
		return nil, err
	}

	// Process results to fit your expected output.
	var results []map[string]string
	for _, hit := range searchResults.Hits {
		result := make(map[string]string)
		for field, value := range hit.Fields {
			result[field] = fmt.Sprintf("%v", value)
		}
		results = append(results, result)
	}

	return results, nil
}
