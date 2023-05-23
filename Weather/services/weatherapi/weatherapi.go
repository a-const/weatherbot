package weatherapi

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// reference
// https://www.weatherapi.com/api-explorer.aspx
// cabinet:
// https://www.weatherapi.com/my/
// max requests for free: 1 million per month

type WeatherApiWeather struct {
	Forecast struct {
		Forecastday []struct {
			Day struct {
				Maxtemp_c      float64 `json:"maxtemp_c"`
				Mintemp_c      float64 `json:"mintemp_c"`
				Totalprecip_mm float64 `json:"totalprecip_mm"`
				Maxwind_kph    float64 `json:"maxwind_kph"`
			} `json:"day"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func Load() *WeatherApiWeather {
	//struct for results
	api_weather := &WeatherApiWeather{}
	// build request
	req_text := "http://api.weatherapi.com/v1/forecast.json?key="
	req_text += viper.GetString("weatherapi_key")
	req_text += "&q=yerevan&days=1&aqi=no&alerts=no"
	req, _ := http.NewRequest("GET", req_text, nil)
	// send req, receive response, and unpack it to struct
	client := &http.Client{}
	if resp, err := client.Do(req); err != nil {
		logrus.Printf("weatherApi GET err! %s", err.Error())
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&api_weather); err != nil {
			logrus.Printf("error while decoding WeatherApi response! %s", err.Error())
		}
	}
	return api_weather
}
