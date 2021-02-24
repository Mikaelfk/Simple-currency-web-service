package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var startTime time.Time

/*
 * A handler function for handling what functions are called when the different urls are visited
 */
func handler() {
	r := mux.NewRouter()
	r.HandleFunc("/exchange/v1/exchangehistory/{country_name}/{begin_date-end_date}", exchangehistory)
	r.HandleFunc("/exchange/v1/exchangeborder/{country_name}", exchangeborder).Queries("limit", "{limit}")
	r.HandleFunc("/exchange/v1/exchangeborder/{country_name}", exchangeborder)
	r.HandleFunc("/exchange/v1/diag/", diag)
	r.HandleFunc("/exchange/v1/diag", diag)
	http.Handle("/", r)
}

/*
 * Main function that initialized the application
 */
func main() {
	startTime = time.Now()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	handler()
	log.Printf("Listening on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
