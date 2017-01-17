// Package router has the necessary functions for creating the router
// and provides the HTTP method handlers
package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"github.com/cbarlas/PeakAPI/model"
	"github.com/cbarlas/PeakAPI/apis"
	"github.com/cbarlas/PeakAPI/db"
	"github.com/wcharczuk/go-chart"
	"fmt"
)

// Creates a new router, initializes its Handler functions and returns it
func CreateRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	HandleRouterFuncs(router)
	return router
}

// Sets the handler functions of the given router
func HandleRouterFuncs(router *mux.Router) {
	router.HandleFunc("/Events/{ApiKey}", GetEvents).Methods("GET")
	router.HandleFunc("/Events", CreateEvent).Methods("POST")
	router.HandleFunc("/Histogram", DrawChart)
}

// Handler function of the GET method upon the /Events/ path
// Returns the Events with the API_KEY given in the request
func GetEvents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	key := params["ApiKey"]	

	if IsApiKeyValid(key) {
		evs := db.GetEventsWithApiKey(key)
		// if the API_KEY is valid but there is no record in database
		if len(evs) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(CreateHTTPError(404))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(evs)
	} else {
		// if the API_KEY is invalid
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(CreateHTTPError(403))
	}

}

// Creates and returns a new HTTPError pointer according to the given status code
func CreateHTTPError(code int) *HTTPError {
	switch code {
	case 400:
		return &HTTPError{code, "Bad Request"}
	case 401:
		return &HTTPError{code, "Unauthorized"}
	case 403:
		return &HTTPError{code, "Forbidden"}
	case 404:
		return &HTTPError{code, "Not Found"}
	default:
		return nil
	
	}
}

// Struct for representing HTTP Error
type HTTPError struct {
	StatusCode int	`json:"statusCode"`
	Message string	`json:"message"`
}

// Checks the API_KEY and returns if it is valid or not
func IsApiKeyValid(s string) bool {
	return apis.API_KEYS[s]
}

// Checks the structure of the Event data
func CheckEvent(e *model.Event) bool {
	return e.ApiKey != "" && e.UserID != 0 && e.Timestamp != 0
}

// Draws a chart of response times if a request is made on /Histogram path
func DrawChart(w http.ResponseWriter, r *http.Request) {
	bars := CreateChartValues()
	sbc := chart.BarChart{
		Height:   512,
		BarWidth: 60,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: bars,
	}

	w.Header().Set("Content-Type", "image/png")
	err := sbc.Render(chart.PNG, w)
	if err != nil {
		fmt.Printf("Error rendering chart: %v\n", err)
	}
}

// Designs the chart according to the histogram buckets
func CreateChartValues() []chart.Value {
	times := db.GetResponseTimes()
	bars := []chart.Value{
		{Value: 0, Label: "<1ms"},
		{Value: 0, Label: "<5ms"},
		{Value: 0, Label: "<10ms"},
		{Value: 0, Label: "<20ms"},
		{Value: 0, Label: "<50ms"},
		{Value: 0, Label: "<100ms"},
	}

	for _, value := range times {
		if value <= 1 {
			bars[0].Value++
		} else if value <= 5 {
			bars[1].Value++
		} else if value <= 10 {
			bars[2].Value++
		} else if value <= 20 {
			bars[3].Value++
		} else if value <= 50 {
			bars[4].Value++
		} else if value <= 100 {
			bars[5].Value++
		}
	}
	return bars
}

// Handler function of the POST method upon the /Events path
// Stores the Event data given in the request body if the data includes valid values
// and returns it as JSON data
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var ev model.Event

	w.Header().Set("Content-Type", "application/json")
	// if there is an error parsing the JSON Event data or the data has an empty value
	if err := json.NewDecoder(r.Body).Decode(&ev); err != nil || !CheckEvent(&ev) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CreateHTTPError(400))
	} else if !IsApiKeyValid(ev.ApiKey) { // if the api key of the data is invalid
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(CreateHTTPError(403))
	} else {
		db.AddEvent(&ev)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ev)
	}
	
}

