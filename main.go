package main

import (
	"encoding/json"
	"expvar"
	"flag"
	"fmt"
	forecast "github.com/danesparza/forecast/v2"
	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
	"github.com/pmylund/go-cache"
	"github.com/rs/cors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/* The data from the pollen service */
type PollenResponse struct {
	StatusMessage string `json:"WebMessage"`
	Status        int    `json:"WebStatus"`
	Successful    bool   `json:"IsSuccess"`

	Entries []PollenData `json:"Entries"`
}

type PollenData struct {
	CityState string `json:"CityState"`
	City      string `json:"City"`
	State     string `json: "State"`

	PredominantPollen string `json: "PredominantPollen"`

	Today     float64 `json:"Today"`
	Tomorrow  float64 `json: "Tomorrow"`
	TwoDays   float64 `json: "TwoDays"`
	ThreeDays float64 `json: "ThreeDays"`
}

/* Pollen return type */
type PollenInfo struct {
	City  string `json:"City"`
	State string `json: "State"`

	PredominantPollen string    `json: "PredominantPollen"`
	PollenCount       []float64 `json: "PollenCount"`
}

const (
	// URL example:  "https://nasacort.com/Ajax/PollenResults.aspx?ZipCode=30022"
	POLLEN_BASEURL = "https://nasacort.com/Ajax/PollenResults.aspx?ZipCode="
)

//	Expvars for cache hits and misses
var forecastCacheHits = expvar.NewInt("Forecast cache hits")
var forecastCacheMisses = expvar.NewInt("Forecast cache misses")
var pollenCacheHits = expvar.NewInt("Pollen cache hits")
var pollenCacheMisses = expvar.NewInt("Pollen cache misses")

func main() {

	//	Set up our flags
	port := flag.Int("port", 3000, "The port to listen on")
	key := flag.String("apikey", "ReplaceWithYourKey", "Your Forecast.io API key")
	expvarUser := flag.String("expvarUser", "changeme", "The username to access expvar stats")
	expvarPass := flag.String("expvarPass", "changeme", "The password to access expvar stats")
	allowedOrigins := flag.String("allowedOrigins", "*", "A comma-separated list of valid CORS origins")

	//	Parse the command line for flags:
	flag.Parse()

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 30 seconds
	c := cache.New(5*time.Minute, 30*time.Second)

	r := mux.NewRouter()
	r.HandleFunc("/forecast/{lat},{long}", func(w http.ResponseWriter, r *http.Request) {

		//	Parse the lat & long from the url
		lat := mux.Vars(r)["lat"]
		long := mux.Vars(r)["long"]

		// 	See if we have the forecast in the cache
		fcast, found := c.Get("forecast-" + lat + "-" + long)
		if !found {
			//	We didn't find it in cache.
			forecastCacheMisses.Add(1)

			//	Call the API with the key and the lat/long
			f, err := forecast.Get(*key, lat, long, "now", forecast.AUTO)

			//	If we have errors, return them using standard HTTP service method
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//	Set the item in cache:
			fcast = f
			c.Set("forecast-"+lat+"-"+long, fcast, cache.DefaultExpiration)
		} else {
			forecastCacheHits.Add(1)
		}

		//	Set the content type header and return the JSON
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(fcast)
	})

	r.HandleFunc("/pollen/{zip}", func(w http.ResponseWriter, r *http.Request) {

		//	Parse the zipcode from the url
		zip := mux.Vars(r)["zip"]

		// 	See if we have the pollen in the cache
		fcast, found := c.Get("pollen-" + zip)
		if !found {
			//	We didn't find it in cache.
			pollenCacheMisses.Add(1)

			//	Call the API with the key and the lat/long
			f, err := GetPollenInfo(zip)

			//	If we have errors, return them using standard HTTP service method
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//	Set the item in cache:
			fcast = f
			c.Set("pollen-"+zip, fcast, cache.DefaultExpiration)
		} else {
			pollenCacheHits.Add(1)
		}

		//	Set the content type header and return the JSON
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, fcast)
	})

	//	Expose debug variables with a simple user/password:
	r.Handle("/debug/vars", httpauth.SimpleBasicAuth(*expvarUser, *expvarPass)(http.DefaultServeMux))

	//	CORS handler
	ch := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(*allowedOrigins, ","),
		AllowCredentials: true,
	})
	handler := ch.Handler(r)

	//	Indicate what port we're starting the service on
	portString := strconv.Itoa(*port)
	fmt.Println("Starting server on :", portString)
	http.ListenAndServe(":"+portString, handler)
}

func GetPollenInfo(zipcode string) (string, error) {

	//	Construct the complete url
	url := POLLEN_BASEURL + zipcode

	//	Go fetch the response from the server:
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	//	Read the body of the response if we have one:
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	//	Unmarshall from JSON into our struct:
	pres := &PollenResponse{}
	if err := json.Unmarshal(body, &pres); err != nil {
		return "", err
	}

	//	Put data into our response struct
	pollenCounts := make([]float64, 4)
	pollenCounts[0] = pres.Entries[0].Today
	pollenCounts[1] = pres.Entries[0].Tomorrow
	pollenCounts[2] = pres.Entries[0].TwoDays
	pollenCounts[3] = pres.Entries[0].ThreeDays
	PollenReturnData := PollenInfo{
		City:              pres.Entries[0].City,
		State:             pres.Entries[0].State,
		PollenCount:       pollenCounts,
		PredominantPollen: pres.Entries[0].PredominantPollen}

	//	Marshall into a JSON string
	PollenReturn, err := json.Marshal(PollenReturnData)
	if err != nil {
		return "", err
	}

	//	Return the JSON string
	return string(PollenReturn), nil
}
