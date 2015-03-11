package main

import (
	"encoding/json"
	"expvar"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	forecast "github.com/mlbright/forecast/v2"
	"github.com/pmylund/go-cache"
	"github.com/rs/cors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//	Expvars for cache hits and misses
var pubCacheHits = expvar.NewInt("Cache hits")
var pubCacheMisses = expvar.NewInt("Cache misses")
var cacheHits int64 = 0
var cacheMisses int64 = 0

func main() {

	//	Set up our flags
	port := flag.Int("port", 3000, "The port to listen on")
	key := flag.String("apikey", "ReplaceWithYourKey", "Your Forecast.io API key")
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
		fcast, found := c.Get("forecast")
		if !found {
			//	We didn't find it in cache.
			cacheMisses++
			pubCacheMisses.Set(cacheMisses)

			//	Call the API with the key and the lat/long
			f, err := forecast.Get(*key, lat, long, "now", forecast.AUTO)

			//	If we have errors, return them using standard HTTP service method
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//	Set the item in cache:
			fcast = f
			c.Set("forecast", fcast, cache.DefaultExpiration)
		} else {
			cacheHits++
			pubCacheHits.Set(cacheHits)
		}

		//	Set the content type header and return the JSON
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(fcast)
	})

	//	We want to expose our debug variables:
	//	r.Handle("/debug/vars", http.DefaultServeMux)

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
