package main

import (
	"log"
	"net/http"
	"pricing-api/pkg/api"
)

func main() {
	router := api.SetupRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

	// basePath := "C:\\Users\\isvan\\OneDrive\\Documents\\work\\GoApi\\data"

	// // Call to index all CSVs in the data directory
	// err := indexAllCSVs(basePath)
	// if err != nil {
	// 	log.Fatalf("Failed to index CSV files: %v", err)
	// }

}
