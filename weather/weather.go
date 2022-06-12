// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package weather

import (
	"fmt"
	"strings"
	"time"

	"ricketyspace.net/peach/nws"
	"ricketyspace.net/peach/version"
)

type Weather struct {
	Title           string
	Version         string
	Location        string
	Now             WeatherNow
	Q2HTimeline     WeatherTimeline // Q2H forecast of the next 12 hours.
	BiDailyTimeline WeatherTimeline // BiDaily forecast for the next 3 days.
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

func NewWeather(point *nws.Point, f, fh *nws.Forecast) (*Weather, error) {
	w := new(Weather)
	w.Location = fmt.Sprintf("%s, %s",
		strings.ToLower(point.Properties.RelativeLocation.Properties.City),
		strings.ToLower(point.Properties.RelativeLocation.Properties.State),
	)
	w.Title = w.Location
	w.Version = version.Version
	w.Now = WeatherNow{
		Temperature:     fh.Properties.Periods[0].Temperature,
		TemperatureUnit: fh.Properties.Periods[0].TemperatureUnit,
		Forecast:        fh.Properties.Periods[0].ShortForecast,
		WindSpeed:       fh.Properties.Periods[0].WindSpeed,
		WindDirection:   fh.Properties.Periods[0].WindDirection,
	}

	// Build Q2H timeline for the 12 hours.
	q2hPeriods := []WeatherPeriod{}
	max := 6
	for i, period := range fh.Properties.Periods {
		if i%2 != 0 {
			continue // Take every other period
		}
		t, err := time.Parse(time.RFC3339, period.StartTime)
		if err != nil {
			return nil, err
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
	for _, period := range f.Properties.Periods {
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

	return w, nil
}
