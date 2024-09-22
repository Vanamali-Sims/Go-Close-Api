package main

import (
	"log"
	"net/http"
	"pricing-api/pkg/api"
)

func main() {
	router := api.SetupRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
