package search

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
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
