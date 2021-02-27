package main

import (
	"errors"
	"io/ioutil"
	"net/http"
)

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
