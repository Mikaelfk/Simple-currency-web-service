package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler() {
	http.HandleFunc("/exchange/v1/exchangehistory/", exchangehistory)
	http.HandleFunc("/exchange/v1/exchangeborder/", exchangeborder)
	http.HandleFunc("/exchange/v1/diag/", diag)
}

func main() {
	handler()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func exchangehistory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exchange History Endpoint")
}

func exchangeborder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exchange border Endpoint")
}

func diag(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Diag Endpoint")
}
