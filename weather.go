package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type weatherProvider interface {
	temperature(city string) (kelvin float64, err error)
}

type openWeatherMap struct{}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func (w openWeatherMap) temperature(city string) (float64, error) {

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
	if err != nil {
		log.Printf("error :: %v", err)
		return 0.0, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		log.Printf("error :: %v", err)
		return 0.0, err
	}

	log.Printf("Got the temperature from OpenWeatherMap t :: %v", d)
	return d.Main.Kelvin, nil
}

type weatherUnderGround struct {
	apiKey string
}

func (w weatherUnderGround) temperature(city string) (float64, error) {
	resp, err := http.Get("http://api.wunderground.com/api/" + w.apiKey + "/conditions/q/" + city + ".json")
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Observation struct {
			Celsius float64 `json:"temp_c"`
		} `json:"current_observation"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	kelvin := d.Observation.Celsius + 273.15
	log.Printf("Got the temperature from weatherUnderground: %s: %.2f", city, kelvin)
	return kelvin, nil
}

type serialMetaWeatherProvider []weatherProvider

func (w serialMetaWeatherProvider) temperature(city string) (float64, error) {

	sum := 0.0

	for _, provider := range w {
		data, err := provider.temperature(city)

		if err != nil {
			log.Println("error in fetching temperature data")
			return 0.0, err
		}

		sum += data
	}

	return sum / float64(len(w)), nil
}

type parallelMetaWeatherProvider []weatherProvider

func (w parallelMetaWeatherProvider) temperature(city string) (float64, error) {

	temps := make(chan float64, len(w))
	errs := make(chan error, len(w))

	for _, provider := range w {
		go func(p weatherProvider) {
			data, err := p.temperature(city)
			if err != nil {
				log.Println("error in fetching temperature data")
				errs <- err
			}
			temps <- data
		}(provider)
	}

	sum := 0.0
	for i := 0; i < len(w); i++ {
		select {
		case temp := <-temps:
			sum += temp
		case err := <-errs:
			return 0, err
		}
	}
	return sum / float64(len(w)), nil
}
