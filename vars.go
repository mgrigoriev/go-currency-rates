package main

const cbrApiUrl = "https://www.cbr-xml-daily.ru/daily.xml"

const httpAddr = "127.0.0.1:9999"

var ratesCache map[string]float64
