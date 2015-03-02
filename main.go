package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	forecast "github.com/mlbright/forecast/v2"
	"net/http"
	"strconv"
)

func main() {

	//	Set up our flags
	port := flag.Int("port", 3000, "The port to listen on")
	key := flag.String("apikey", "ReplaceWithYourKey", "Your Forecast.io API key")

	//	Parse the command line for flags:
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/forecast/{lat},{long}", func(w http.ResponseWriter, r *http.Request) {

		//	Parse the lat & long from the url
		lat := mux.Vars(r)["lat"]
		long := mux.Vars(r)["long"]

		//	Call the API with the key and the lat/long
		f, err := forecast.Get(*key, lat, long, "now", forecast.AUTO)

		//	If we have errors, return them using standard HTTP service method
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//	Set the content type header and return the JSON
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(f)
	})

	//	Indicate what port we're starting the service on
	portString := strconv.Itoa(*port)
	fmt.Println("Starting server on :", portString)
	http.ListenAndServe(":"+portString, r)
}
