package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type IWeatherForecast interface {
	Formated() string
}

type WeatherForecast struct {
	DateTime     time.Time
	TemperatureC int32
	TemperatureF int32
	Summary      string
}

func (forecast WeatherForecast) Formated() string {
	formated := fmt.Sprintf(
		"%s | %s | C° %d | F° %d",
		forecast.DateTime,
		forecast.Summary,
		forecast.TemperatureC,
		forecast.TemperatureF,
	)

	return formated
}

func GetTemperatureF(temperatureC int32) int32 {
	farenheit := float32(temperatureC) / 0.5556
	roundedFahrenheit := 32 + int(farenheit)
	return int32(roundedFahrenheit)
}

func authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing authorization middleware")

		bearer := r.Header.Get("Authorization")

		if bearer == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		parts := strings.Split(bearer, " ")

		if len(parts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func greet(w http.ResponseWriter, r *http.Request) {
	log.Print("Executing greet")
	sumaries := []string{
		"Freezing",
		"Bracing",
		"Chilly",
		"Cool",
		"Mild",
		"Warm",
		"Balmy",
		"Hot",
		"Sweltering",
		"Scorching",
	}

	index := rand.Intn(len(sumaries))

	maxTemperature := int32(50)
	minTemperature := int32(-22)
	temperatureC := rand.Int31n(maxTemperature-minTemperature) + int32(minTemperature)

	forecast := WeatherForecast{
		DateTime:     time.Now(),
		TemperatureC: temperatureC,
		TemperatureF: GetTemperatureF(temperatureC),
		Summary:      sumaries[index],
	}

	log.Print(forecast.Formated())

	json.NewEncoder(w).Encode(forecast)
}

func main() {
	greetHandle := http.HandlerFunc(greet)
	http.Handle("/greet", authorize(greetHandle))
	http.ListenAndServe(":3000", nil)

	log.Print("Listening on :3000...")
}
