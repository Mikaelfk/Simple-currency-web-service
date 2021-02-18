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
	"time"

	"github.com/gorilla/mux"
)

// Currency is a struct used to store information about a currency
// retrieved from the third-party api's
type Currency struct {
	Code   string
	Name   string
	Symbol string
}

// Information is a struct used to store information about a country
// retrieved from the third-party api's
type Information struct {
	Currencies []Currency
	Borders    []string
}

/*
 * A function which retrieves the exchange history of a currency of a requested country
 * between a requested time period
 */
func exchangehistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Exchange History Endpoint")
	vars := mux.Vars(r)
	countryName := vars["country_name"]
	// Requests the country's currencies
	body, err := getResponse("https://restcountries.eu/rest/v2/name/"+countryName+"?fields=currencies", w)
	// If there is an error, log it and return
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	// Saves the currencies in an array of Information structs
	// Because of the way the third-party API works, information needs to be an array
	var information []Information
	if err := json.Unmarshal([]byte(string(body)), &information); err != nil {
		// Handles json parsing error
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	beginDateEndDate := vars["begin_date-end_date"]
	// Splits the dates into an array
	splitdate := strings.Split(beginDateEndDate, "-")

	if len(splitdate) < 6 {
		// Handles string error
		err = errors.New("Error in date query")
		log.Printf("Error, %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Splits the dates into two differen strings
	beginDate := splitdate[0] + "-" + splitdate[1] + "-" + splitdate[2]
	endDate := splitdate[3] + "-" + splitdate[4] + "-" + splitdate[5]
	body, err = getResponse("https://api.exchangeratesapi.io/history?start_at="+beginDate+"&end_at="+endDate+"&symbols="+information[0].Currencies[0].Code, w)
	// If any errors occur, log them and return
	if err != nil {
		log.Printf("Error, %v", err)
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
	// Gets the requested country's currency and its bordering countries
	body, err := getResponse("https://restcountries.eu/rest/v2/name/"+countryName+"?fields=borders;currencies", w)

	// If any errors occur, log them and return
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Stores the JSON information in the information variable
	var information []Information
	if err = json.Unmarshal([]byte(string(body)), &information); err != nil {
		// Handles json parsing error
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	limit := len(information[0].Borders)
	// Checks if the user has requested a limit on how many countries should be checked
	val, ok := vars["limit"]
	newLimit, err := strconv.Atoi(val)
	if ok && newLimit < limit {
		limit = newLimit
	}
	var currencies []string
	// Saves the requested country's currency in index 0 in the currencies array
	currencies = append(currencies, information[0].Currencies[0].Code)
	// Requests all the bordering countries' currencies
	// This variable does not need to be a slice, again because of how to API works
	var information2 Information
	for i := 0; i < limit; i++ {
		body, err = getResponse("https://restcountries.eu/rest/v2/alpha/"+information[0].Borders[i]+"?fields=currencies", w)
		// If any errors occur, log them and return
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		if err := json.Unmarshal([]byte(string(body)), &information2); err != nil {
			// Handles json parsing error
			log.Printf("Error: %v", err)
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if len(information2.Currencies) != 0 {
			currencies = append(currencies, information2.Currencies[0].Code)
		}
	}
	// Removes all duplicate currencies in the slice
	currencies = unique(currencies)
	var currenciesRequest string
	validCodes := []string{"CAD", "HKD", "ISK", "PHP", "DKK", "HUF", "CZK", "GBP", "RON", "SEK", "IDR", "INR", "BRL",
		"RUB", "HRK", "JPY", "THB", "CHF", "EUR", "MYR", "BGN", "TRY", "CNY", "NOK", "NZD", "ZAR", "USD", "MXN", "SGD", "AUD", "ILS", "KRW", "PLN"}
	// Saves the currencies in a single string with a comma between each currency
	for i := 1; i < len(currencies); i++ {
		if stringInSlice(currencies[i], validCodes) {
			currenciesRequest += currencies[i] + ","
		}
	}
	currenciesRequest = strings.TrimRight(currenciesRequest, ",")
	// Requests the bordering countries exchange rates
	if len(currenciesRequest) == 0 {
		log.Printf("Error: No bordering countries or no bordering countries with an available currency")
		http.Error(w, "Error: No bordering countries or no bordering countries with an available currency", http.StatusBadRequest)
		return
	}
	body, err = getResponse("https://api.exchangeratesapi.io/latest?symbols="+currenciesRequest+";base="+currencies[0], w)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Fprintf(w, string(body))
}

/*
 * a function which returns some information about the program, including status codes for the API's used,
 * version and uptime
 */
func diag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Diag Endpoint")
	var exchangeStatusCode int
	var countriesStatusCode int
	// Does a basic request to the exchange rates api.
	respExchange, err := http.Get("https://api.exchangeratesapi.io/")
	// If any errors occur, log it and set the status code to 500,
	// otherwise set the status code to the recieved status code
	if err != nil {
		log.Printf("Something went wrong with the exchange rates api, %v", err)
		exchangeStatusCode = 500
	} else {
		exchangeStatusCode = respExchange.StatusCode
		defer respExchange.Body.Close()
	}
	// Does a basic request to the countries information api.
	respCountries, err := http.Get("https://restcountries.eu")
	// If any errors occur, log it and set the status code to 500,
	// otherwise set the status code to the recieved status code
	if err != nil {
		log.Printf("Something went wrong with the countries api, %v", err)
		countriesStatusCode = 500
	} else {
		countriesStatusCode = respCountries.StatusCode
		defer respCountries.Body.Close()
	}
	// Print the information in JSON format
	fmt.Fprintf(w, `{"exchangeratesapi": "%v", "restcountries": "%v", "version": "v1", "uptime": "%v Seconds"}`,
		exchangeStatusCode, countriesStatusCode, int(time.Since(startTime)/time.Second))
}

/*
 * a function for getting the string response from a http get request
 */
func getResponse(request string, w http.ResponseWriter) ([]byte, error) {
	resp, err := http.Get(request)

	if err != nil {
		// Handles retrieval errors
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// Handles user input error
		http.Error(w, "Error: Error in query", resp.StatusCode)
		return nil, errors.New("Error in query")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// Handles body read error
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return body, nil
}

/*
 * A function for removing all duplicate elements from a string slice
 */
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

/*
 * A function for checking if a string is in a slice
 * Returns a boolean value
 */
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
