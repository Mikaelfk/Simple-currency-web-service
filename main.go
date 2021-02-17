package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var startTime time.Time

/*
 * A handler function for handling what functions are called when the different urls are visited
 */
func handler() {
	startTime = time.Now()
	r := mux.NewRouter()
	r.HandleFunc("/exchange/v1/exchangehistory/{country_name}/{begin_date-end_date}", exchangehistory)
	r.HandleFunc("/exchange/v1/exchangeborder/{country_name}", exchangeborder).Queries("limit", "{limit}")
	r.HandleFunc("/exchange/v1/exchangeborder/{country_name}", exchangeborder)
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
