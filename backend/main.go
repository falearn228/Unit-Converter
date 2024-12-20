package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

var lengthConversionFactors = map[string]float64 {
	"millimeter": 1,            // базовая единица
	"centimeter": 10,           // 1 cm = 10 mm
	"meter":      1000,         // 1 m = 1000 mm
	"kilometer":  1000000,      // 1 km = 1000000 mm
	"inch":       25.4,         // 1 inch = 25.4 mm
	"foot":       304.8,        // 1 foot = 304.8 mm
	"yard":       914.4,        // 1 yard = 914.4 mm
	"mile":       1609344,  
}

var weightConversionFactors = map[string]float64{
	"milligram": 1,
	"gram":      1000,
	"kilogram":  1000000,
	"ounce":     28349.5,
	"pound":     453592,
}

type ConversionRequest struct {
	Value float64 `json:"value"` // Значение для перевода
	From  string  `json:"from"`  // Исходная единица
	To    string  `json:"to"`    // Конечная единица
	Type  string  `json:"type"`  // Тип (length, weight, temperature)
}

type ConversionResponse struct {
	Result float64 `json:"result"` // Результат перевода
}

func convertLengthOrWeight(value float64, from string, to string, factors map[string]float64) (float64, error) {
	fromFactor, fromExists := factors[from]
	toFactor, toExists := factors[to]
	if !fromExists || !toExists {
		return 0, fmt.Errorf("invalid units")
	}
	return (value * fromFactor) / toFactor, nil
}

func convertTemperature(value float64, from string, to string) (float64, error) {
	if from == to {
		return value, nil
	}

	switch from {
	case "Celsius":
		if to == "Fahrenheit" {
			return (value * 9 / 5) + 32, nil
		} else if to == "Kelvin" {
			return value + 273.15, nil
		}
	case "Fahrenheit":
		if to == "Celsius" {
			return (value - 32) * 5 / 9, nil
		} else if to == "Kelvin" {
			return (value-32)*5/9 + 273.15, nil
		}
	case "Kelvin":
		if to == "Celsius" {
			return value - 273.15, nil
		} else if to == "Fahrenheit" {
			return (value-273.15)*9/5 + 32, nil
		}
	}
	return 0, fmt.Errorf("invalid temperature units")
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request ConversionRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}


	var result float64

	// Определяем тип перевода
	switch request.Type {
	case "length":
		result, err = convertLengthOrWeight(request.Value, request.From, request.To, lengthConversionFactors)
	case "weight":
		result, err = convertLengthOrWeight(request.Value, request.From, request.To, weightConversionFactors)
	case "temperature":
		result, err = convertTemperature(request.Value, request.From, request.To)
	default:
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := ConversionResponse{
		Result: math.Round(result*100) / 100, // Округляем до 2 знаков
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/convert", handler)

	fmt.Println("Starting HTTPS server on https://localhost:8443")
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	if err != nil {
		fmt.Println("Error while starting server...", err)
	}
}