package main

import (
	"fmt"
	"io/ioutil"
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

/*
 * A function which retrieves the exchange history of a currency of a requested country
 * between a requested time period.
 */
func exchangehistory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exchange History Endpoint")
	vars := mux.Vars(r)
	countryName := vars["country_name"]
	resp, err := http.Get("https://restcountries.eu/rest/v1/name/" + countryName + "?fields=borders;currencies")
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Fprint(w, "Error")
	} else if err != nil {
		// handle error
		fmt.Fprint(w, "Error")
	} else {
		fmt.Println(vars)
		fmt.Fprint(w, string(body))
	}
}

func exchangeborder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exchange border Endpoint")
}

func diag(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Diag Endpoint")
}
