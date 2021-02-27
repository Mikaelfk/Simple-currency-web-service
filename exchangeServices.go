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
	Name       string
	Currencies []Currency
	Borders    []string
}

// Rates is a struct used to store all the exchange rates available in a map
type Rates struct {
	Rates map[string]float32
	Base  string
}

// Exchangeborder is a struct used to store information which will be converted to json
// for displaying exchange rates for bordering countries
type Exchangeborder struct {
	Rates map[string]Countryinfo `json:"rates"`
	Base  string                 `json:"base"`
}

// Countryinfo is a strcut used to store the name of the currency and the exchange rate of that currency
type Countryinfo struct {
	Currency string  `json:"currency"`
	Rate     float32 `json:"rate"`
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
	// If any errors occur, log them and return
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
	// This variable has to be a slice, I dont know why
	var information []Information
	if err = json.Unmarshal([]byte(string(body)), &information); err != nil {
		// Handles json parsing error
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Initialized variable that will be used to save all information that will be shown to the user
	var exchangeborderResponse Exchangeborder

	// Sets the base currency
	exchangeborderResponse.Base = information[0].Currencies[0].Code

	limit := len(information[0].Borders)
	// Checks if the user has requested a limit on how many countries should be checked
	val, ok := vars["limit"]
	newLimit, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Error: %v", err)
	} else if ok && newLimit < limit {
		limit = newLimit
	}
	// Makes a map that will be used to store what currency each country has
	// The key in the map is the country, while the entry is the currency
	var allCurrencies map[string]string
	allCurrencies = make(map[string]string)

	// Creates a string that will store all the bordering countries country code
	var allBorderCodes string

	// Adds all the bordering countries country codes to the string
	for i := 0; i < limit; i++ {
		allBorderCodes += information[0].Borders[i] + ";"
	}
	allBorderCodes = strings.TrimRight(allBorderCodes, ";")
	// Checks if the string is empty
	if allBorderCodes == "" {
		err = errors.New("No bordering countries")
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Requests the currencies and the name of the bordering countries
	body, err = getResponse("https://restcountries.eu/rest/v2/alpha?codes="+allBorderCodes+"&fields=currencies;name", w)
	// If any errors occur, log them and return
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	// Unmarshal the information into information2
	var information2 []Information
	if err := json.Unmarshal([]byte(string(body)), &information2); err != nil {
		// Handles json parsing error
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// If there are fewer countries than the set limit, reduce the limit to the amount of countries
	if len(information2) < limit {
		limit = len(information2)
	}
	// Save the currencies in the allCurrencies map
	for i := 0; i < limit; i++ {
		allCurrencies[information2[i].Name] = information2[i].Currencies[0].Code
	}

	// A request to get all the exchange rates available
	body, err = getResponse("https://api.exchangeratesapi.io/latest", w)
	// If any errors occur, log them and return
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	var rates Rates
	if err := json.Unmarshal([]byte(string(body)), &rates); err != nil {
		// Handles json parsing error
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Adds all the valid currency codes to the validCodes slice
	var validCodes []string
	for key := range rates.Rates {
		validCodes = append(validCodes, key)
	}
	// Checks if all the currencies are valid.
	var allCurrenciesSlice []string
	for _, entry := range allCurrencies {
		if stringInSlice(entry, validCodes) {
			allCurrenciesSlice = append(allCurrenciesSlice, entry)
		}
	}
	// Checks if the base currency is valid
	if stringInSlice(exchangeborderResponse.Base, validCodes) != true {
		err = errors.New("Chosen country's currency is not available")
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Makes sure that there are bordering countries with available currencies
	if len(allCurrenciesSlice) == 0 {
		err = errors.New("No bordering countries with an available currency")
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Create a string uniqueCurrencies, which will be used to request the exchange rates for all the currencies.
	// Add all the currencies to this string, as long as they are unique
	var uniqueCurrencies string
	allCurrenciesSlice = unique(allCurrenciesSlice)
	for i := 0; i < len(allCurrenciesSlice); i++ {
		if stringInSlice(allCurrenciesSlice[i], validCodes) && allCurrenciesSlice[i] != exchangeborderResponse.Base {
			uniqueCurrencies += allCurrenciesSlice[i] + ","
		}
	}
	uniqueCurrencies = strings.TrimRight(uniqueCurrencies, ",")

	// Requests the exchange rates of the currencies we have
	body, err = getResponse("https://api.exchangeratesapi.io/latest?symbols="+uniqueCurrencies+"&base="+exchangeborderResponse.Base, w)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	// Saves the information in the rates2 variable
	var rates2 Rates
	if err := json.Unmarshal([]byte(string(body)), &rates2); err != nil {
		// Handles json parsing error
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize a map which will store the information that will be shown to the user.
	var responseMap map[string]Countryinfo
	responseMap = make(map[string]Countryinfo)

	// Add all exchange rates and currencies to the responseMap variable, delete countries that do not have available currencies
	for key, entry := range allCurrencies {
		if stringInSlice(entry, validCodes) {
			exchangeRate := rates2.Rates[entry]
			if entry == exchangeborderResponse.Base {
				exchangeRate = 1
			}
			responseMap[key] = Countryinfo{entry, exchangeRate}
		} else {
			delete(allCurrencies, key)
		}
	}

	exchangeborderResponse.Rates = responseMap
	// Marshal the information into a variable, that will be shown to the user
	b, err := json.Marshal(exchangeborderResponse)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, string(b))
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
	fmt.Fprintf(w, `{"exchangeratesapi": "%v", "restcountries": "%v", "version": "v1", "uptime": %v}`,
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
