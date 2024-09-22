package api

import (
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/getCloseUSD", GetCloseUSDHandler).Methods("GET")
	router.HandleFunc("/getCloseInBetween", GetCloseInBetweenHandler).Methods("GET")
	return router
}
