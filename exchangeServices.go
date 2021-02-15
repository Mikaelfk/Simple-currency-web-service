package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
		// Handles retrieval errors
		log.Printf("Bad request, %v", err)
		return
	}

	if resp.StatusCode != 200 {
		// Handles user input error
		log.Println("Could not retrieve country")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// Handles body read error
		log.Printf("Body read error, %v", err)
		return
	}

	var information []Information

	if err := json.Unmarshal([]byte(string(body)), &information); err != nil {
		// Handles json parsing error
		log.Printf("Body parse error, %v", err)
		return
	}
	beginDateEndDate := vars["begin_date-end_date"]
	splitdate := strings.Split(beginDateEndDate, "-")

	if len(splitdate) < 6 {
		// Handles string error
		log.Printf("Error in date query")
		return
	}
	beginDate := splitdate[0] + "-" + splitdate[1] + "-" + splitdate[2]
	endDate := splitdate[3] + "-" + splitdate[4] + "-" + splitdate[5]
	respEx, err := http.Get("https://api.exchangeratesapi.io/history?start_at=" + beginDate + "&end_at=" + endDate + "&symbols=" + information[0].Currencies[0].Code)
	if err != nil {
		// Handles retrieval error
		log.Printf("Bad Request, %v", err)
		return
	}
	if respEx.StatusCode != 200 {
		// Handles user input error
		log.Println("Could not retrieve exchange rates")
		return
	}
	bodyEx, err := ioutil.ReadAll(respEx.Body)

	if err != nil {
		// Handles body read error
		log.Printf("Body read error, %v", err)
		return
	}

	fmt.Fprintf(w, "%s", string(bodyEx))
}

func exchangeborder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Exchange border Endpoint")

	vars := mux.Vars(r)
	countryName := vars["country_name"]
	resp, err := http.Get("https://restcountries.eu/rest/v2/name/" + countryName + "?fields=borders")

	if err != nil {
		// Handles retrieval errors
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Bad request, %v", err)
		return
	}

	if resp.StatusCode != 200 {
		// Handles user input error
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Could not retrieve country")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// Handles body read error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Body read error, %v", err)
		return
	}

	var information []Information

	if err := json.Unmarshal([]byte(string(body)), &information); err != nil {
		// Handles json parsing error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Body parse error, %v", err)
		return
	}

	limit := len(information[0].Borders)
	if val, ok := vars["limit"]; ok {
		limit, _ = strconv.Atoi(val)
	}
	fmt.Println(limit)
	currencies := make([]string, 0, limit)
	var information2 Information
	for i := 0; i < limit; i++ {
		resp, err = http.Get("https://restcountries.eu/rest/v2/alpha/" + information[0].Borders[i] + "?fields=currencies")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Bad request, %v", err)
			return
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			// Handles body read error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Body read error, %v", err)
			return
		}
		if err := json.Unmarshal([]byte(string(body)), &information2); err != nil {
			// Handles json parsing error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Body parse error, %v", err)
			return
		}
		currencies[i] = information2.Currencies[0].Code
		fmt.Println(information2.Currencies[0].Code)
	}
	fmt.Println(currencies)

}

func diag(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Diag Endpoint")
}
