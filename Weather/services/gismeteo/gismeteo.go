package gismeteo

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// reference
// https://www.gismeteo.ru/api/
// key through mail

type GismeteoWeather struct {
	Temperature struct {
		Air struct {
			Celcius float64 `json:"C"`
		} `json:"air"`
	} `json:"temperature"`
	Pressure struct {
		Mm_hg_atm float64 `json:"mm_hg_atm"`
	} `json:"pressure"`
	Wind struct {
		Speed struct {
			M_s int `json:"m_s"`
		} `json:"speed"`
	} `json:"wind"`
}

func Load() *GismeteoWeather {
	//struct for results
	gis_weather := &GismeteoWeather{}
	// build request
	req, _ := http.NewRequest("GET", "https://api.gismeteo.net/v2/weather/current/4368/?lang=ru&units=metric", nil)
	req.Header.Set("X-Gismeteo-Token", "56b30cb255.3443075")
	// send req, receive response, and unpack it to struct
	client := &http.Client{}
	if resp, err := client.Do(req); err != nil {
		logrus.Printf("gismeteo GET err! %s", err.Error())
	} else {
		defer resp.Body.Close()
		//bodyBytes, _ := io.ReadAll(resp.Body)
		//bodyString := string(bodyBytes)
		//fmt.Print(bodyString)
		if err := json.NewDecoder(resp.Body).Decode(&gis_weather); err != nil {
			logrus.Printf("error while decoding gismeteo response! %s", err.Error())
		}
	}
	return gis_weather
}
