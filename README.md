# Assignment 1


This application has three different endpoints:

```
/exchange/v1/exchangehistory/
/exchange/v1/exchangeborder/
/exchange/v1/diag/
```
# How to use

## Exchange History
To get the exchange rate history of a certain country's currency, use the endpoint
```
exchange/v1/exchangehistory/
```
### Request
```
METHOD: GET
Path: exchangehistory/{:country_name}/{:begin_date-end_date}
```
The request requires two mandatory inputs: country_name and begin_date-end_date.

```{:country_name}``` is the english name of the country.

```{:begin_date-end_date}``` refers to the begin date and end date of the period over which exchange rates are shown.

If a country has multiple currencies, only the first one is reported.

The date format is yyyy-mm-dd.

Example request: ```exchangehistory/sweden/2020-01-31-2020-02-29```


## Exchange Rates of Bordering Countries Currencies
To get the exchange rates of bordering countries currencies, use the endpoint 
```
/exchange/v1/exchangeborder/
```
### Request
```
METHOD: GET
Path: exchangeborder/{:country_name}{?limit={:number}}
```
The request requires one mandatory input:

```{:country_name}``` is the english name of the country.

and one optional input:

```{?limit={:number}}``` limits the number of currencies of surrounding countries to be reported. 

If a country has multiple currencies, only the first one is reported.

If no currency is reported, the country is ignored.

Example request: ```/exchangeborder/norway?limit=4```


## Diagnostics interface
To access the diagnostics interface, use the endpoint
```
/exchange/v1/diag/
```
### Request
```
METHOD: GET
Path: diag/
```

## Deployment
The service is deployed with [Heroku](https://heroku.com):

https://assignment-1-mikaelfk.herokuapp.com

Examples:

https://assignment-1-mikaelfk.herokuapp.com/exchange/v1/exchangehistory/norway/2020-12-01-2021-01-31

https://assignment-1-mikaelfk.herokuapp.com/exchange/v1/exchangeborder/sweden

https://assignment-1-mikaelfk.herokuapp.com/exchange/v1/diag/
