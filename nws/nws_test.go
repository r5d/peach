// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package nws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestPoints(t *testing.T) {
	// Test valid lat,lng.
	np, err := Points(41.115, -83.177)
	if err != nil {
		t.Errorf("points: %v", err)
		return
	}
	if np.Properties.Forecast != "https://api.weather.gov/gridpoints/CLE/33,42/forecast" {
		t.Errorf("points: forcecast link: '%v'", np.Properties.Forecast)
	}
	if np.Properties.ForecastHourly != "https://api.weather.gov/gridpoints/CLE/33,42/forecast/hourly" {
		t.Errorf("points: forcecast link: '%v'", np.Properties.ForecastHourly)
	}
	if np.Properties.GridId != "CLE" {
		t.Errorf("points: gridid: %v", np.Properties.GridId)
	}
	if np.Properties.GridX != 33 {
		t.Errorf("points: gridx: %v", np.Properties.GridX)
	}
	if np.Properties.GridY != 42 {
		t.Errorf("points: gridy: %v", np.Properties.GridY)
	}

	if np.Properties.RelativeLocation.Properties.City != "Tiffin" {
		t.Errorf("points: location: city: %v", np.Properties)
	}
	if np.Properties.RelativeLocation.Properties.State != "OH" {
		t.Errorf("points: location: state: %v", np.Properties)
	}

	// Test invalid lat,lng
	np, err = Points(115.0, -83.177)
	if err == nil {
		t.Errorf("points: %v", np)
	}
}

func TestGetForecast(t *testing.T) {
	// Get point.
	np, nwsErr := Points(41.115, -83.177)
	if nwsErr != nil {
		t.Errorf("error: %v", nwsErr)
		return
	}

	// Get forecast.
	fc, nwsErr := GetForecast(np)
	if nwsErr != nil {
		t.Errorf("error: %v", nwsErr)
		return
	}

	// Verify periods.
	for i, period := range fc.Properties.Periods {
		if period.Number < 1 {
			t.Errorf("period: %d: number invalid: %v", i, period.Number)
		}
		if len(period.Name) < 1 {
			t.Errorf("period: %d: name invalid: %v", i, period.Name)
		}
		if len(period.StartTime) < 1 {
			t.Errorf("period: %d: start time invalid: %v", i,
				period.StartTime)
		}
		if len(period.EndTime) < 1 {
			t.Errorf("period: %d: end time invalid: %v", i,
				period.EndTime)
		}
		if len(period.TemperatureUnit) < 1 {
			t.Errorf("period: %d: temperature unit invalid: %v",
				i, period.TemperatureUnit)
		}
		if len(period.WindSpeed) < 1 {
			t.Errorf("period: %d: wind speed invalid: %v",
				i, period.WindSpeed)
		}
		if len(period.WindDirection) < 1 {
			t.Errorf("period: %d: wind direction invalid: %v",
				i, period.WindDirection)
		}
		if len(period.ShortForecast) < 1 {
			t.Errorf("period: %d: short forecast invalid: %v",
				i, period.ShortForecast)
		}
		if len(period.DetailedForecast) < 1 {
			t.Errorf("period: %d: detailed forecast invalid: %v",
				i, period.DetailedForecast)
		}
	}
}

func TestGetForecastHourly(t *testing.T) {
	// Get point.
	np, nwsErr := Points(41.115, -83.177)
	if nwsErr != nil {
		t.Errorf("error: %v", nwsErr)
		return
	}

	// Get forecast hourly.
	fc, nwsErr := GetForecastHourly(np)
	if nwsErr != nil {
		t.Errorf("error: %v", nwsErr)
		return
	}

	// Verify periods.
	for i, period := range fc.Properties.Periods {
		if period.Number < 1 {
			t.Errorf("period: %d: number invalid: %v", i, period.Number)
		}
		if len(period.StartTime) < 1 {
			t.Errorf("period: %d: start time invalid: %v", i,
				period.StartTime)
		}
		if len(period.EndTime) < 1 {
			t.Errorf("period: %d: end time invalid: %v", i,
				period.EndTime)
		}
		if len(period.TemperatureUnit) < 1 {
			t.Errorf("period: %d: temperature unit invalid: %v",
				i, period.TemperatureUnit)
		}
		if len(period.WindSpeed) < 1 {
			t.Errorf("period: %d: wind speed invalid: %v",
				i, period.WindSpeed)
		}
		if len(period.WindDirection) < 1 {
			t.Errorf("period: %d: wind direction invalid: %v",
				i, period.WindDirection)
		}
		if len(period.ShortForecast) < 1 {
			t.Errorf("period: %d: short forecast invalid: %v",
				i, period.ShortForecast)
		}
	}
}

func TestNWSGetWrapper(t *testing.T) {
	// Initialize test NWS server.
	fails := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fails > 0 {
			fails -= 1
			http.Error(w, `{"type":"urn:noaa:nws:api:UnexpectedProblem","title":"Unexpected Problem","status":500,"detail":"An unexpected problem has occurred.","instance":"urn:noaa:nws:api:request:493c3a1d-f87e-407f-ae2c-24483f5aab63","correlationId":"493c3a1d-f87e-407f-ae2c-24483f5aab63","additionalProp1":{}}`, 500)
			return
		}

		// Add expires header.
		w.Header().Set("expires",
			time.Now().Add(time.Second*60).Format(time.RFC1123))

		// Success.
		fmt.Fprintln(w, `{"@context":[],"properties":{"gridId":"CLE","gridX":82,"gridY":64,"forecast":"https://api.weather.gov/gridpoints/CLE/82,64/forecast","forecastHourly":"https://api.weather.gov/gridpoints/CLE/82,64/forecast/hourly","relativeLocation":{"properties":{"city":"Cleveland","state":"OH"}}}}`)
	}))
	defer ts.Close()

	// Test 1 - Server fails 5 times.
	fails = 5
	_, _, err := get(ts.URL)
	if err != nil {
		t.Errorf("get failed: %v", err)
		return
	}

	// Test 2 - Server fails 6 times.
	fails = 6
	respBody, _, err := get(ts.URL)
	if err == nil {
		t.Errorf("get did not fail: %s", respBody)
		return
	}
	if err != nil && respBody != nil {
		t.Errorf("body is not nil: %s", respBody)
	}
	if err.Title != "Unexpected Problem" {
		t.Errorf("err title: %s", err.Title)
		return
	}
	if err.Type != "urn:noaa:nws:api:UnexpectedProblem" {
		t.Errorf("err type: %s", err.Type)
		return
	}
	if err.Status != 500 {
		t.Errorf("err status: %d", err.Status)
		return
	}
	if err.Detail != "An unexpected problem has occurred." {
		t.Errorf("err detail: %s", err.Detail)
		return
	}

	// Test 3 - Server fails 1 time.
	fails = 1
	respBody, expires, err := get(ts.URL)
	if err != nil {
		t.Errorf("get failed: %v", err)
		return
	}
	if respBody == nil {
		t.Errorf("body: %s", respBody)
		return
	}
	if time.Until(expires).Seconds() < 1 {
		t.Errorf("points: expires in not in the future")
		return
	}
	point := new(Point)
	jerr := json.Unmarshal(respBody, point)
	if jerr != nil {
		t.Errorf("points: decode: %v", jerr)
		return
	}
	if point.Properties.Forecast == "" {
		t.Errorf("points: forecast empty")
		return
	}
	if point.Properties.ForecastHourly == "" {
		t.Errorf("points: forecasthourly empty")
		return
	}
	if point.Properties.RelativeLocation.Properties.City == "" {
		t.Errorf("points: city empty")
		return
	}
	if point.Properties.RelativeLocation.Properties.State == "" {
		t.Errorf("points: state empty")
		return
	}
}

func TestAlerts(t *testing.T) {
	// Initialize test NWS server.
	fail := false
	alert := `{"@context":["https://geojson.org/geojson-ld/geojson-context.jsonld",{"@version":"1.1","wx":"https://api.weather.gov/ontology#","@vocab":"https://api.weather.gov/ontology#"}],"type":"FeatureCollection","features":[{"id":"https://api.weather.gov/alerts/urn:oid:2.49.0.1.840.0.b17674895f0389f4cb276c6db9c15d0aed41460f.001.1","type":"Feature","geometry":null,"properties":{"@id":"https://api.weather.gov/alerts/urn:oid:2.49.0.1.840.0.b17674895f0389f4cb276c6db9c15d0aed41460f.001.1","@type":"wx:Alert","id":"urn:oid:2.49.0.1.840.0.b17674895f0389f4cb276c6db9c15d0aed41460f.001.1","areaDesc":"Coffee; Dale; Henry; Geneva; Houston; North Walton; Central Walton; Holmes; Washington; Jackson; Inland Bay; Calhoun; Inland Gulf; Inland Franklin; Gadsden; Leon; Inland Jefferson; Madison; Liberty; Inland Wakulla; Inland Taylor; Lafayette; Inland Dixie; South Walton; Coastal Bay; Coastal Gulf; Coastal Franklin; Coastal Jefferson; Coastal Wakulla; Coastal Taylor; Coastal Dixie; Quitman; Clay; Randolph; Calhoun; Terrell; Dougherty; Lee; Worth; Turner; Tift; Ben Hill; Irwin; Early; Miller; Baker; Mitchell; Colquitt; Cook; Berrien; Seminole; Decatur; Grady; Thomas; Brooks; Lowndes; Lanier","geocode":{"SAME":["001031","001045","001067","001061","001069","012131","012059","012133","012063","012005","012013","012045","012037","012039","012073","012065","012079","012077","012129","012123","012067","012029","013239","013061","013243","013037","013273","013095","013177","013321","013287","013277","013017","013155","013099","013201","013007","013205","013071","013075","013019","013253","013087","013131","013275","013027","013185","013173"],"UGC":["ALZ065","ALZ066","ALZ067","ALZ068","ALZ069","FLZ007","FLZ008","FLZ009","FLZ010","FLZ011","FLZ012","FLZ013","FLZ014","FLZ015","FLZ016","FLZ017","FLZ018","FLZ019","FLZ026","FLZ027","FLZ028","FLZ029","FLZ034","FLZ108","FLZ112","FLZ114","FLZ115","FLZ118","FLZ127","FLZ128","FLZ134","GAZ120","GAZ121","GAZ122","GAZ123","GAZ124","GAZ125","GAZ126","GAZ127","GAZ128","GAZ129","GAZ130","GAZ131","GAZ142","GAZ143","GAZ144","GAZ145","GAZ146","GAZ147","GAZ148","GAZ155","GAZ156","GAZ157","GAZ158","GAZ159","GAZ160","GAZ161"]},"affectedZones":["https://api.weather.gov/zones/forecast/ALZ065","https://api.weather.gov/zones/forecast/ALZ066","https://api.weather.gov/zones/forecast/ALZ067","https://api.weather.gov/zones/forecast/ALZ068","https://api.weather.gov/zones/forecast/ALZ069","https://api.weather.gov/zones/forecast/FLZ007","https://api.weather.gov/zones/forecast/FLZ008","https://api.weather.gov/zones/forecast/FLZ009","https://api.weather.gov/zones/forecast/FLZ010","https://api.weather.gov/zones/forecast/FLZ011","https://api.weather.gov/zones/forecast/FLZ012","https://api.weather.gov/zones/forecast/FLZ013","https://api.weather.gov/zones/forecast/FLZ014","https://api.weather.gov/zones/forecast/FLZ015","https://api.weather.gov/zones/forecast/FLZ016","https://api.weather.gov/zones/forecast/FLZ017","https://api.weather.gov/zones/forecast/FLZ018","https://api.weather.gov/zones/forecast/FLZ019","https://api.weather.gov/zones/forecast/FLZ026","https://api.weather.gov/zones/forecast/FLZ027","https://api.weather.gov/zones/forecast/FLZ028","https://api.weather.gov/zones/forecast/FLZ029","https://api.weather.gov/zones/forecast/FLZ034","https://api.weather.gov/zones/forecast/FLZ108","https://api.weather.gov/zones/forecast/FLZ112","https://api.weather.gov/zones/forecast/FLZ114","https://api.weather.gov/zones/forecast/FLZ115","https://api.weather.gov/zones/forecast/FLZ118","https://api.weather.gov/zones/forecast/FLZ127","https://api.weather.gov/zones/forecast/FLZ128","https://api.weather.gov/zones/forecast/FLZ134","https://api.weather.gov/zones/forecast/GAZ120","https://api.weather.gov/zones/forecast/GAZ121","https://api.weather.gov/zones/forecast/GAZ122","https://api.weather.gov/zones/forecast/GAZ123","https://api.weather.gov/zones/forecast/GAZ124","https://api.weather.gov/zones/forecast/GAZ125","https://api.weather.gov/zones/forecast/GAZ126","https://api.weather.gov/zones/forecast/GAZ127","https://api.weather.gov/zones/forecast/GAZ128","https://api.weather.gov/zones/forecast/GAZ129","https://api.weather.gov/zones/forecast/GAZ130","https://api.weather.gov/zones/forecast/GAZ131","https://api.weather.gov/zones/forecast/GAZ142","https://api.weather.gov/zones/forecast/GAZ143","https://api.weather.gov/zones/forecast/GAZ144","https://api.weather.gov/zones/forecast/GAZ145","https://api.weather.gov/zones/forecast/GAZ146","https://api.weather.gov/zones/forecast/GAZ147","https://api.weather.gov/zones/forecast/GAZ148","https://api.weather.gov/zones/forecast/GAZ155","https://api.weather.gov/zones/forecast/GAZ156","https://api.weather.gov/zones/forecast/GAZ157","https://api.weather.gov/zones/forecast/GAZ158","https://api.weather.gov/zones/forecast/GAZ159","https://api.weather.gov/zones/forecast/GAZ160","https://api.weather.gov/zones/forecast/GAZ161"],"references":[],"sent":"2022-06-17T16:00:00-04:00","effective":"2022-06-17T16:00:00-04:00","onset":"2022-06-18T12:00:00-04:00","expires":"2022-06-17T22:00:00-04:00","ends":"2022-06-18T20:00:00-04:00","status":"Actual","messageType":"Alert","category":"Met","severity":"Moderate","certainty":"Likely","urgency":"Expected","event":"Heat Advisory","sender":"w-nws.webmaster@noaa.gov","senderName":"NWS Tallahassee FL","headline":"Heat Advisory issued June 17 at 4:00PM EDT until June 18 at 8:00PM EDT by NWS Tallahassee FL","description":"* WHAT...Heat index values of 108 to 112 expected.\n\n* WHERE...Portions of southeast Alabama, south central and\nsouthwest Georgia, and the Big Bend and Panhandle of Florida.\n\n* WHEN...For the first Heat Advisory, until 8 PM EDT /7 PM CDT/\nthis evening. For the second Heat Advisory, from noon EDT /11\nAM CDT/ to 8 PM EDT /7 PM CDT/ Saturday.\n\n* IMPACTS...Hot temperatures and high humidity may cause heat\nillnesses to occur.","instruction":"Drink plenty of fluids, stay in an air-conditioned room, stay out\nof the sun, and check up on relatives and neighbors. Young\nchildren and pets should never be left unattended in vehicles\nunder any circumstances.\n\nTake extra precautions if you work or spend time outside. When\npossible reschedule strenuous activities to early morning or\nevening. Know the signs and symptoms of heat exhaustion and heat\nstroke. Wear lightweight and loose fitting clothing when\npossible. To reduce risk during outdoor work, the Occupational\nSafety and Health Administration recommends scheduling frequent\nrest breaks in shaded or air conditioned environments. Anyone\novercome by heat should be moved to a cool and shaded location.\nHeat stroke is an emergency! Call 9 1 1.","response":"Execute","parameters":{"AWIPSidentifier":["NPWTAE"],"WMOidentifier":["WWUS72 KTAE 172000"],"NWSheadline":["HEAT ADVISORY REMAINS IN EFFECT UNTIL 8 PM EDT /7 PM CDT/ THIS EVENING... ...HEAT ADVISORY IN EFFECT FROM NOON EDT /11 AM CDT/ TO 8 PM EDT /7 PM CDT/ SATURDAY"],"BLOCKCHANNEL":["EAS","NWEM","CMAS"],"VTEC":["/O.NEW.KTAE.HT.Y.0005.220618T1600Z-220619T0000Z/"],"eventEndingTime":["2022-06-19T00:00:00+00:00"]}}}],"title":"current watches, warnings, and advisories for 30.7848 N, 83.5601 W","updated":"2022-06-18T00:00:00+00:00"}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fail {
			http.Error(w, `{"type":"urn:noaa:nws:api:UnexpectedProblem","title":"Unexpected Problem","status":500,"detail":"An unexpected problem has occurred.","instance":"urn:noaa:nws:api:request:493c3a1d-f87e-407f-ae2c-24483f5aab63","correlationId":"493c3a1d-f87e-407f-ae2c-24483f5aab63","additionalProp1":{}}`, 500)
			return
		}

		// Add expires header.
		w.Header().Set("expires",
			time.Now().Add(time.Second*60).Format(time.RFC1123))

		// Success.
		fmt.Fprint(w, alert)
	}))
	defer ts.Close()
	baseUrl, _ = url.Parse(ts.URL)

	// Hit it.
	fail = true
	_, nwsErr := GetAlerts(33.2938, -83.9674)
	if nwsErr == nil {
		t.Errorf("alerts: expected it to fail")
		return
	}
	if nwsErr.Title != "Unexpected Problem" {
		t.Errorf("alerts: Error.Title: %v", nwsErr.Title)
		return
	}
	if nwsErr.Type != "urn:noaa:nws:api:UnexpectedProblem" {
		t.Errorf("alerts: Error.Type: %v", nwsErr.Type)
		return
	}
	if nwsErr.Status != 500 {
		t.Errorf("alerts: Error.Status: %v", nwsErr.Status)
		return
	}
	if nwsErr.Detail != "An unexpected problem has occurred." {
		t.Errorf("alerts: Error.Detail: %v", nwsErr.Detail)
		return
	}

	// Hit it again.
	fail = false
	fc, nwsErr := GetAlerts(33.2938, -83.9674)
	if nwsErr != nil {
		t.Errorf("alerts: %v", nwsErr)
		return
	}

	// Validate the feature collection.
	if len(fc.Features) < 1 {
		t.Errorf("alerts: fc.Features < 1")
		return
	}

	// Validate feature.
	f := fc.Features[0]
	if f.Properties.Event != "Heat Advisory" {
		t.Errorf("alerts: feture.Properties.Event: %v", f.Properties.Event)
		return
	}
	if f.Properties.Severity != "Moderate" {
		t.Errorf("alerts: feture.Properties.Severity: %v", f.Properties.Severity)
		return
	}
	if f.Properties.Description != "* WHAT...Heat index values of 108 to 112 expected.\n\n* WHERE...Portions of southeast Alabama, south central and\nsouthwest Georgia, and the Big Bend and Panhandle of Florida.\n\n* WHEN...For the first Heat Advisory, until 8 PM EDT /7 PM CDT/\nthis evening. For the second Heat Advisory, from noon EDT /11\nAM CDT/ to 8 PM EDT /7 PM CDT/ Saturday.\n\n* IMPACTS...Hot temperatures and high humidity may cause heat\nillnesses to occur." {
		t.Errorf("alerts: feture.Properties.Description: %v", f.Properties.Description)
		return
	}
	if f.Properties.Instruction != "Drink plenty of fluids, stay in an air-conditioned room, stay out\nof the sun, and check up on relatives and neighbors. Young\nchildren and pets should never be left unattended in vehicles\nunder any circumstances.\n\nTake extra precautions if you work or spend time outside. When\npossible reschedule strenuous activities to early morning or\nevening. Know the signs and symptoms of heat exhaustion and heat\nstroke. Wear lightweight and loose fitting clothing when\npossible. To reduce risk during outdoor work, the Occupational\nSafety and Health Administration recommends scheduling frequent\nrest breaks in shaded or air conditioned environments. Anyone\novercome by heat should be moved to a cool and shaded location.\nHeat stroke is an emergency! Call 9 1 1." {
		t.Errorf("alerts: feture.Properties.Instruction: %v", f.Properties.Instruction)
		return
	}

	// Check if the alert is cached.
	ll := fmt.Sprintf("%.4f,%.4f", 33.2938, -83.9674)
	cachedRawAlert := aCache.Get(ll)
	if len(cachedRawAlert) < 1 {
		t.Errorf("alerts: cache: empty: %v", cachedRawAlert)
		return
	}
	if len(cachedRawAlert) != len(alert) {
		t.Errorf("alerts: cached entry size: %d != %d",
			len(cachedRawAlert), len(alert))
		return
	}
	if string(cachedRawAlert[:]) != alert {
		t.Errorf("alerts: cached entry does not match alert: %s != %s",
			cachedRawAlert, alert)
		return
	}
}
