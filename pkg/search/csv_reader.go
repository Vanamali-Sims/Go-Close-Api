package search

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/blevesearch/bleve"
)

type Record struct {
	Date   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

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

func queryIndex(index bleve.Index, queryStr string) ([]map[string]string, error) {
	query := bleve.NewQueryStringQuery(queryStr) // Use QueryStringQuery for more flexible queries
	searchRequest := bleve.NewSearchRequest(query)

	log.Printf("Executing search for query: %s", queryStr)

	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("error executing search: %v", err)
	}

	log.Printf("Found %d hits for query: %s", searchResult.Total, queryStr)

	var results []map[string]string
	for _, hit := range searchResult.Hits {
		doc, err := index.Document(hit.ID)
		if err != nil {
			continue // Optionally handle error
		}

		data := make(map[string]string)
		for _, field := range doc.Fields {
			fieldName := field.Name()
			fieldValue := string(field.Value())
			data[fieldName] = fieldValue
		}
		results = append(results, data)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for the query: %s", queryStr)
	}

	return results, nil
}

func indexCSV(filePath string, index bleve.Index) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	headers, err := reader.Read() // Assuming the first row is headers
	if err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err // Handle non-EOF errors
		}

		data := make(map[string]string)
		for i, value := range record {
			if i < len(headers) {
				data[headers[i]] = value
			}
		}

		// Create a Record struct to be indexed
		rec := Record{
			Date:   data["Date"],
			Open:   parseFloat(data["Open"]),
			High:   parseFloat(data["High"]),
			Low:    parseFloat(data["Low"]),
			Close:  parseFloat(data["Close"]),
			Volume: parseInt(data["Volume"]),
		}

		// Indexing the record using Bleve
		if err := index.Index(data["Date"], rec); err != nil {
			return err
		}
	}
	return nil
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func parseInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
