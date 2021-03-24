package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"unicode"
)

var ErrInvalidZip = errors.New("invalid zip code")
var ErrLookupFailure = errors.New("unable to retrieve data")

var key = os.Getenv("WEATHER_TOKEN")
var weatherURL = "https://api.openweathermap.org/data/2.5/weather?units=imperial&zip=%s,US&appid=%s"

// Generated with the help of https://mholt.github.io/json-to-go/
type weatherReport struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
}

// Weather contains weather data for a zip code.
type Weather struct {
	Temperature   float64   `json:"temperature,omitempty"`
	Humidity      int       `json:"humidity,omitempty"`
	WindSpeed     float64   `json:"wind_speed,omitempty"`
	WindDirection int       `json:"wind_direction,omitempty"`
	Timestamp     time.Time `json:"timestamp,omitempty"`
}

// Do a live lookup.
func lookup(zip string) (Weather, error) {
	var w Weather
	resp, err := http.Get(fmt.Sprintf(weatherURL, zip, key))
	if err != nil {
		return w, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return w, err
	}

	var report weatherReport
	err = json.Unmarshal(b, &report)
	w.Temperature = report.Main.Temp
	w.Humidity = report.Main.Humidity
	w.WindSpeed = report.Wind.Speed
	w.WindDirection = report.Wind.Deg
	w.Timestamp = time.Now()
	return w, err
}

// Get accepts a zip code and provides weather data and an error.
func Get(zip string) (Weather, error) {
	zip = strip(zip)
	if !isValidZip(zip) {
		return Weather{}, ErrInvalidZip
	}
	w, found := getCache(zip)
	if found {
		return w, nil
	}

	w, err := lookup(zip)
	if err != nil {
		return w, ErrLookupFailure
	}
	setCache(zip, w)
	return w, nil
}

// Strip non-ints and return first five characters.
func strip(zip string) string {
	out := make([]rune, 0, len(zip))
	for _, r := range []rune(zip) {
		if unicode.IsDigit(r) {
			out = append(out, r)
		}
	}
	if len(out) > 5 {
		out = out[:5]
	}
	return string(out)
}

func isValidZip(zip string) bool {
	if len(zip) != 5 {
		return false
	}
	_, found := us[zip[:3]]
	return found
}

func pruneCache() {
	for {
		time.Sleep(cacheDuration)
		mu.Lock()
		cutoff := time.Now().Add(-cacheDuration)
		for zip, w := range cache {
			if w.Timestamp.Before(cutoff) {
				go deleteCache(zip)
			}
		}
		mu.Unlock()
	}
}
