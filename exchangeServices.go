package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Exchange History Endpoint")
	vars := mux.Vars(r)
	countryName := vars["country_name"]
	resp, err := http.Get("https://restcountries.eu/rest/v2/name/" + countryName + "?fields=borders;currencies")

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	var information []Information

	if resp.StatusCode != 200 {
		fmt.Fprint(w, "Error while retrieving country currency")
	} else {
		json.Unmarshal([]byte(string(body)), &information)
		beginDateEndDate := vars["begin_date-end_date"]
		splitdate := strings.Split(beginDateEndDate, "-")
		if len(splitdate) < 6 {
			// String error
			fmt.Fprint(w, "Error in string")
		} else {
			beginDate := splitdate[0] + "-" + splitdate[1] + "-" + splitdate[2]
			endDate := splitdate[3] + "-" + splitdate[4] + "-" + splitdate[5]
			respEx, err := http.Get("https://api.exchangeratesapi.io/history?start_at=" + beginDate + "&end_at=" + endDate + "&symbols=" + information[0].Currencies[0].Code)
			if err != nil {
				// handle error
				log.Fatal(err)
			}

			bodyEx, _ := ioutil.ReadAll(respEx.Body)

			if respEx.StatusCode != 200 {
				//Error while retrieving exchange rates
				fmt.Fprint(w, "Error while retrieving exchange rates")
			} else {
				var prettyJSON bytes.Buffer
				json.Indent(&prettyJSON, bodyEx, "", "\t")

				fmt.Fprintf(w, "%s", string(prettyJSON.Bytes()))
			}
		}
	}
}

func exchangeborder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exchange border Endpoint")
}

func diag(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Diag Endpoint")
}
