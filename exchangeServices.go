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
	resp, err := http.Get("https://restcoutries.eu/rest/v2/name/" + countryName + "?fields=borders;currencies")

	if err != nil {
		// Handles retrieval errors
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Bad request, %v", err)
		w.WriteHeader(500)
		return
	}

	if resp.StatusCode != 200 {
		// Handles user input error
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Could not retrieve country")
		w.WriteHeader(400)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// Handles body read error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Body read error, %v", err)
		w.WriteHeader(500)
		return
	}

	var information []Information

	if err := json.Unmarshal([]byte(string(body)), &information); err != nil {
		// Handles json parsing error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Body parse error, %v", err)
		w.WriteHeader(500)
		return
	}
	beginDateEndDate := vars["begin_date-end_date"]
	splitdate := strings.Split(beginDateEndDate, "-")

	if len(splitdate) < 6 {
		// Handles tring error
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error in date query")
		w.WriteHeader(400)
		return
	}
	beginDate := splitdate[0] + "-" + splitdate[1] + "-" + splitdate[2]
	endDate := splitdate[3] + "-" + splitdate[4] + "-" + splitdate[5]
	respEx, err := http.Get("https://api.exchangeratesapi.io/history?start_at=" + beginDate + "&end_at=" + endDate + "&symbols=" + information[0].Currencies[0].Code)
	if err != nil {
		// Handles retrieval error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Bad Request, %v", err)
		w.WriteHeader(500)
		return
	}
	if respEx.StatusCode != 200 {
		// Handles user input error
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Could not retrieve exchange rates")
		w.WriteHeader(400)
		return
	}
	bodyEx, err := ioutil.ReadAll(respEx.Body)

	if err != nil {
		// Handles body read error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Body read error, %v", err)
		w.WriteHeader(500)
		return
	}
	var prettyJSON bytes.Buffer
	jsonErr := json.Indent(&prettyJSON, bodyEx, "", "\t")
	if jsonErr != nil {
		// Handles json indenting error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("JSON indenting error, %v", err)
		w.WriteHeader(500)
		return
	}
	fmt.Fprintf(w, "%s", string(prettyJSON.Bytes()))
}

func exchangeborder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exchange border Endpoint")
}

func diag(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Diag Endpoint")
}
