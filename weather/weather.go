// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package weather

import (
	"fmt"
	"strings"
	"time"

	"ricketyspace.net/peach/nws"
	"ricketyspace.net/peach/photon"
	"ricketyspace.net/peach/version"
)

type Weather struct {
	Title           string
	Version         string
	Location        string
	Now             WeatherNow
	Q2HTimeline     WeatherTimeline // Q2H forecast of the next 12 hours.
	BiDailyTimeline WeatherTimeline // BiDaily forecast for the next 3 days.
	SearchEnabled   bool
	Alerts          []Alert
}

type WeatherNow struct {
	Temperature     int
	TemperatureUnit string
	Forecast        string
	WindSpeed       string
	WindDirection   string
}

type WeatherPeriod struct {
	Name            string
	Forecast        string
	Hour            int
	Temperature     int
	TemperatureUnit string
}

type WeatherTimeline struct {
	Periods []WeatherPeriod
}

type Alert struct {
	Event       string
	Severity    string
	Description []string
	Instruction []string
}

func NewWeather(lat, lng float32) (*Weather, error, int) {
	fBundle, nwsErr := nws.GetForecastBundle(lat, lng)
	if nwsErr != nil {
		return nil, nwsErr, nwsErr.Status
	}

	w := new(Weather)
	w.Location = fmt.Sprintf("%s, %s",
		strings.ToLower(fBundle.Point.Properties.RelativeLocation.Properties.City),
		strings.ToLower(fBundle.Point.Properties.RelativeLocation.Properties.State),
	)
	w.Title = w.Location
	w.Version = version.Version
	w.Now = WeatherNow{
		Temperature:     fBundle.ForecastHourly.Properties.Periods[0].Temperature,
		TemperatureUnit: fBundle.ForecastHourly.Properties.Periods[0].TemperatureUnit,
		Forecast:        fBundle.ForecastHourly.Properties.Periods[0].ShortForecast,
		WindSpeed:       fBundle.ForecastHourly.Properties.Periods[0].WindSpeed,
		WindDirection:   fBundle.ForecastHourly.Properties.Periods[0].WindDirection,
	}

	// Build Q2H timeline for the 12 hours.
	q2hPeriods := []WeatherPeriod{}
	max := 6
	for i, period := range fBundle.ForecastHourly.Properties.Periods {
		if i%2 != 0 {
			continue // Take every other period
		}
		t, err := time.Parse(time.RFC3339, period.StartTime)
		if err != nil {
			return nil, err, 500
		}
		p := WeatherPeriod{
			Forecast:        period.DetailedForecast,
			Hour:            t.Hour(),
			Temperature:     period.Temperature,
			TemperatureUnit: period.TemperatureUnit,
		}
		q2hPeriods = append(q2hPeriods, p)
		if len(q2hPeriods) == max {
			break
		}
	}
	w.Q2HTimeline = WeatherTimeline{
		Periods: q2hPeriods,
	}

	// Build BiDaily  timeline for the next 3 days.
	bdPeriods := []WeatherPeriod{}
	max = 8
	for _, period := range fBundle.Forecast.Properties.Periods {
		p := WeatherPeriod{
			Name:            period.Name,
			Forecast:        period.DetailedForecast,
			Temperature:     period.Temperature,
			TemperatureUnit: period.TemperatureUnit,
		}
		bdPeriods = append(bdPeriods, p)
		if len(bdPeriods) == max {
			break
		}
	}
	w.BiDailyTimeline = WeatherTimeline{
		Periods: bdPeriods,
	}
	w.SearchEnabled = photon.Enabled()

	// Add alerts if they exist.
	if len(fBundle.Alerts.Features) > 0 {
		w.Alerts = make([]Alert, 0)
		am := make(map[string]bool, 0) // Alerts map.
		for _, f := range fBundle.Alerts.Features {
			if _, ok := am[f.Id]; ok {
				continue // Duplicate; skip.
			}
			a := Alert{
				Event:       f.Properties.Event,
				Severity:    f.Properties.Severity,
				Description: strings.Split(f.Properties.Description, "\n\n"),
				Instruction: strings.Split(f.Properties.Instruction, "\n\n"),
			}
			w.Alerts = append(w.Alerts, a)
			am[f.Id] = true
		}
	}
	return w, nil, 200
}
