package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	fmt.Println("Booting up the server....")
	http.HandleFunc("/weather/", weatherHandler)
	http.ListenAndServe(":8000", nil)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()
	city := strings.SplitN(r.URL.Path, "/", 3)[2]

	pt := r.URL.Query().Get("pt")
	log.Printf("passed query parameter pt :: %v", pt)

	var m weatherProvider

	ovm := openWeatherMap{}
	wug := weatherUnderGround{apiKey: "c5e7c706ec76acd6"}

	if pt == "s" {
		m = serialMetaWeatherProvider{ovm, wug}
	} else {
		// default means parallel fetching of weather info
		m = parallelMetaWeatherProvider{ovm, wug}
	}
	data, err := m.temperature(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name": city,
		"temp": data,
		"took": time.Since(begin).String(),
	})
}
