package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Currency struct {
	Code   string
	Name   string
	Symbol string
}

type Information struct {
	Currencies []Currency
	Borders    []string
}

/*
 * A function which retrieves the exchange history of a currency of a requested country
 * between a requested time period.
 */
func exchangehistory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exchange History Endpoint")
	vars := mux.Vars(r)
	countryName := vars["country_name"]
	resp, err := http.Get("https://restcountries.eu/rest/v2/name/" + countryName + "?fields=borders;currencies")
	body, err := ioutil.ReadAll(resp.Body)

	var information []Information

	if resp.StatusCode != 200 {
		fmt.Fprint(w, "Error")
	} else if err != nil {
		// handle error
		fmt.Fprint(w, "Error")
	} else {
		fmt.Fprint(w, string(body)+"\n")
		json.Unmarshal([]byte(string(body)), &information)
		fmt.Fprint(w, information)
	}
}

func exchangeborder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exchange border Endpoint")
}

func diag(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Diag Endpoint")
}
