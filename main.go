package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/*
 * A handler function for handling what functions are called when the different urls are visited
 */
func handler() {
	r := mux.NewRouter()
	r.HandleFunc("/exchange/v1/exchangehistory/{country_name}/{begin_date-end_date}", exchangehistory).Queries("limit", "{limit}")
	r.HandleFunc("/exchange/v1/exchangehistory/{country_name}/{begin_date-end_date}", exchangehistory)
	r.HandleFunc("/exchange/v1/exchangeborder/", exchangeborder)
	r.HandleFunc("/exchange/v1/diag/", diag)
	http.Handle("/", r)
}

/*
 * Main function that initialized the application
 */
func main() {
	handler()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
