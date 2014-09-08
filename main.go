package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Booting up the server....")
	r := gin.Default()
	r.GET("/weather/:city", weatherHandler)
	r.Run(":8000")
}

func weatherHandler(c *gin.Context) {
	r := c.Request
	begin := time.Now()
	city := c.Params.ByName("city")

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
		c.Fail(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"name": city,
		"temp": data,
		"took": time.Since(begin).String(),
	})
}
