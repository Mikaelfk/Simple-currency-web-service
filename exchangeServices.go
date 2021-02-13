package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

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
