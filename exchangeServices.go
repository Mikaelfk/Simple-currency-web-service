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
	//Requests the country's currencies
	body, err := getResponse("https://restcountries.eu/rest/v2/name/" + countryName + "?fields=currencies")

	if err != nil {
		fmt.Fprintf(w, `{"error":"%v"}`, err)
		return
	}
	//Saves the currencies in an array of Information structs
	//Because of the way the third-party API works, information needs to be an array.
	var information []Information
	if err := json.Unmarshal([]byte(string(body)), &information); err != nil {
		// Handles json parsing error
		log.Printf("Body parse error, %v", err)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
		return
	}

	beginDateEndDate := vars["begin_date-end_date"]
	//Splits the dates into an array
	splitdate := strings.Split(beginDateEndDate, "-")

	if len(splitdate) < 6 {
		// Handles string error
		err = errors.New("Error in date query")
		log.Printf("Error, %v", err)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
		return
	}
	//Splits the dates into two differen strings
	beginDate := splitdate[0] + "-" + splitdate[1] + "-" + splitdate[2]
	endDate := splitdate[3] + "-" + splitdate[4] + "-" + splitdate[5]
	body, err = getResponse("https://api.exchangeratesapi.io/history?start_at=" + beginDate + "&end_at=" + endDate + "&symbols=" + information[0].Currencies[0].Code)
	if err != nil {
		fmt.Fprintf(w, `{"error":"%v"}`, err)
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
	//Gets the requested country's currency and its bordering countries
	body, err := getResponse("https://restcountries.eu/rest/v2/name/" + countryName + "?fields=borders;currencies")

	if err != nil {
		// Handles body read error
		log.Printf("Body read error, %v", err)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
		return
	}

	//Stores the JSON information in the information variable
	var information []Information
	if err = json.Unmarshal([]byte(string(body)), &information); err != nil {
		// Handles json parsing error
		log.Printf("Body parse error, %v", err)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
		return
	}

	limit := len(information[0].Borders)
	//Checks if the user has requested a limit on how many countries should be checked
	val, ok := vars["limit"]
	newLimit, err := strconv.Atoi(val)
	if ok && newLimit < limit {
		limit = newLimit
	}
	var currencies []string
	//Saves the requested country's currency in index 0 in the currencies array.
	currencies = append(currencies, information[0].Currencies[0].Code)
	//Requests all the bordering countries' currencies
	//This variable does not need to be a slice, again because of how to API works
	var information2 Information
	for i := 0; i < limit; i++ {
		body, err = getResponse("https://restcountries.eu/rest/v2/alpha/" + information[0].Borders[i] + "?fields=currencies")
		if err != nil {
			fmt.Fprintf(w, `{"error":"%v"}`, err)
			return
		}
		if err := json.Unmarshal([]byte(string(body)), &information2); err != nil {
			// Handles json parsing error
			log.Printf("Body parse error, %v", err)
			fmt.Fprintf(w, `{"error":"%v"}`, err)
			return
		}
		if len(information2.Currencies) != 0 {
			currencies = append(currencies, information2.Currencies[0].Code)
		}
	}
	//Removes all duplicate currencies in the slice
	currencies = unique(currencies)
	var currenciesRequest string
	validCodes := []string{"CAD", "HKD", "ISK", "PHP", "DKK", "HUF", "CZK", "GBP", "RON", "SEK", "IDR", "INR", "BRL",
		"RUB", "HRK", "JPY", "THB", "CHF", "EUR", "MYR", "BGN", "TRY", "CNY", "NOK", "NZD", "ZAR", "USD", "MXN", "SGD", "AUD", "ILS", "KRW", "PLN"}
	//Saves the currencies in a single string with a comma between each currency.
	for i := 1; i < len(currencies); i++ {
		if stringInSlice(currencies[i], validCodes) {
			currenciesRequest += currencies[i] + ","
		}
	}
	currenciesRequest = strings.TrimRight(currenciesRequest, ",")
	//Requests the bordering countries exchange rates
	fmt.Println(currenciesRequest)
	if len(currenciesRequest) == 0 {
		fmt.Fprint(w, `{"error":"Country has no bordering countries, or no bordering countries with available currencies"}`)
		return
	}
	body, err = getResponse("https://api.exchangeratesapi.io/latest?symbols=" + currenciesRequest + ";base=" + currencies[0])
	if err != nil {
		fmt.Fprintf(w, `{"error":"%v"}`, err)
		return
	}
	fmt.Fprintf(w, string(body))
}

func diag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Diag Endpoint")
	var exchangeStatusCode int
	var countriesStatusCode int
	respExchange, err := http.Get("https://api.exchangeratesapi.io/")
	if err != nil {
		log.Printf("Bad request, %v", err)
		exchangeStatusCode = 404
	} else {
		exchangeStatusCode = respExchange.StatusCode
		defer respExchange.Body.Close()
	}
	respCountries, err := http.Get("https://restcountries.eu")
	if err != nil {
		log.Printf("Bad request, %v", err)
		countriesStatusCode = 404
	} else {
		countriesStatusCode = respCountries.StatusCode
		defer respCountries.Body.Close()
	}
	fmt.Fprintf(w, `{"exchangeratesapi": "%v", "restcountries": "%v", "version": "v1", "uptime": "%v Seconds"}`,
		exchangeStatusCode, countriesStatusCode, int(time.Since(startTime)/time.Second))
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

	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		// Handles user input error
		log.Printf("Status code is not 2xx")
		return nil, errors.New("Status code is not 200")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// Handles body read error
		log.Printf("Body read error, %v", err)
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
