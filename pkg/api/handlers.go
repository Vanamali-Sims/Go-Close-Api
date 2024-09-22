package api

import (
	"encoding/json"
	"net/http"
	"pricing-api/pkg/search"
)

type GetCloseUSDRequest struct {
	AssetClass     string `json:"assetClass"`
	InternalSymbol string `json:"internalSymbol"`
	Date           string `json:"date"`
	Token          string `json:"token"`
}

type GetCloseInBetweenRequest struct {
	AssetClass     string `json:"assetClass"`
	InternalSymbol string `json:"internalSymbol"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
	Candle         string `json:"candle"`
	Token          string `json:"token"`
}

func GetCloseUSDHandler(w http.ResponseWriter, r *http.Request) {
	var req GetCloseUSDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Token != "ACTUAL_TOKEN" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//Searching for close price implementation.
	result, err := search.GetCloseUSD(req.AssetClass, req.InternalSymbol, req.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func GetCloseInBetweenHandler(w http.ResponseWriter, r *http.Request) {
	var req GetCloseInBetweenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//token validation
	if req.Token != "ACTUAL_TOKEN" {
		http.Error(w, "Unauthorised", http.StatusUnauthorized)
		return
	}

	result, err := search.GetCloseInBetween(req.AssetClass, req.InternalSymbol, req.StartDate, req.EndDate, req.Candle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
