package main

import "fmt"
import "net/http"
import "encoding/json"
import "strings"

func main() {
	fmt.Println("Hello World!!!")
	http.HandleFunc("/hello", defaultHandler)
	http.HandleFunc("/weather/", weatherHandler)
	http.ListenAndServe(":8000", nil)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!!!\n")
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {

	city := strings.SplitN(r.URL.Path, "/", 3)[2]

	data, err := query(city)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)

}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func query(city string) (weatherData, error) {

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	return d, nil
}