package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", HomePage)
}

type LocationResponse struct {
	Key  string `json:"Key"`
	Name string `json:"LocalizedName"`
}

type WeatherResponse struct {
	Temperature struct {
		Metric struct {
			Value float64 `json:"Value"`
		} `json:"Metric"`
	} `json:"Temperature"`
	WeatherText string `json:"WeatherText"`
}

func getLocationKey(city, apiKey string) (string, error) {
	url := fmt.Sprintf("http://dataservice.accuweather.com/locations/v1/cities/search?apikey=%s&q=%s", apiKey, city)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var locations []LocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil || len(locations) == 0 {
		return "", err
	}
	return locations[0].Key, nil
}

func getWeather(locationKey, apiKey string) (*WeatherResponse, error) {
	url := fmt.Sprintf("http://dataservice.accuweather.com/currentconditions/v1/%s?apikey=%s", locationKey, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weather []WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil || len(weather) == 0 {
		return nil, err
	}
	return &weather[0], nil
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		city := r.FormValue("city")
		apiKey := os.Getenv("ACCUWEATHER_API_KEY")

		if locationKey, err := getLocationKey(city, apiKey); err == nil {
			if weatherData, err := getWeather(locationKey, apiKey); err == nil {
				tmpl, _ := template.ParseFiles("index.html")
				tmpl.Execute(w, weatherData)
				return
			}
		}
		http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
	} else {
		tmpl, _ := template.ParseFiles("index.html")
		tmpl.Execute(w, nil)
	}
}
