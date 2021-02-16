package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//Currency is a struct used to store information about a currency
//retrieved from the third-party api's
type Currency struct {
	Code   string
	Name   string
	Symbol string
}

//Information is a struct used to store information about a country
//retrieved from the third-party api's
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

	body, err := getResponse("https://restcountries.eu/rest/v2/name/" + countryName + "?fields=borders;currencies")

	if err != nil {
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
	body, err = getResponse("https://api.exchangeratesapi.io/history?start_at=" + beginDate + "&end_at=" + endDate + "&symbols=" + information[0].Currencies[0].Code)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "%s", string(body))
}

/*
 * A function used to gather current exchange rates from the bordering countries of
 * a requested country
 */
func exchangeborder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Exchange border Endpoint")

	vars := mux.Vars(r)
	countryName := vars["country_name"]
	body, err := getResponse("https://restcountries.eu/rest/v2/name/" + countryName + "?fields=borders")

	if err != nil {
		// Handles body read error
		log.Printf("Body read error, %v", err)
		return
	}

	var information []Information

	if err = json.Unmarshal([]byte(string(body)), &information); err != nil {
		// Handles json parsing error
		log.Printf("Body parse error, %v", err)
		return
	}

	limit := len(information[0].Borders)
	val, ok := vars["limit"]
	newLimit, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Conversion error, %v", err)
		return
	}
	if ok && newLimit < limit {
		limit = newLimit
	}
	currencies := make([]string, limit, limit)
	var information2 Information
	for i := 0; i < limit; i++ {
		body, err = getResponse("https://restcountries.eu/rest/v2/alpha/" + information[0].Borders[i] + "?fields=currencies")
		if err != nil {
			return
		}
		if err := json.Unmarshal([]byte(string(body)), &information2); err != nil {
			// Handles json parsing error
			log.Printf("Body parse error, %v", err)
			return
		}
		currencies[i] = information2.Currencies[0].Code
	}
	fmt.Print(currencies)
}

func diag(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Diag Endpoint")
}

/*
 * a function for getting the string response from a http get request
 */
func getResponse(request string) ([]byte, error) {
	resp, err := http.Get(request)

	if err != nil {
		// Handles retrieval errors
		log.Printf("Bad request, %v", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		// Handles user input error
		log.Println("Could not retrieve country")
		return nil, errors.New("Could not retrieve country")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// Handles body read error
		log.Printf("Body read error, %v", err)
		return nil, err
	}
	return body, nil
}
