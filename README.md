# forecast-service
A pass-through microservice for the [forecast.io](http://forecast.io/) API, using [mlbright/forecast](https://github.com/mlbright/forecast)

[![Build Status](https://drone.io/github.com/danesparza/forecast-service/status.png)](https://drone.io/github.com/danesparza/forecast-service/latest)

To build, make sure you have the latest version of [Go](http://golang.org/) installed.  If you've never used Go before, it's a quick install and [there are installers for multiple platforms](http://golang.org/doc/install), including Windows, Linux and OSX.

### Quick Start

Run the following commands get latest and build.

```bash
go get github.com/danesparza/forecast-service
go build
```

### Starting and testing the service
To start the service, just run `forecast-service`.  

If you need help, just run `forecast-service --help`.

There are a few command line parameters available:

Parameter       | Description
----------      | -----------
apikey          | The Forecast.io api key use for making calls.  You'll need to supply your own key, but they are free.  You can get one at the [Forecast.io developer site](https://developer.forecast.io/)
port            | The port the service listens on.  
allowedOrigins  | comma seperated list of [CORS](http://en.wikipedia.org/wiki/Cross-origin_resource_sharing) origins to allow.  In order to access the service directly from a javascript application, you'll need to specify the origin you'll be running the javascript site on.  For example: http://www.myjavascriptapplication.com

Once the service is up and running, you can connect to it using
`http://yourhostname:3000/forecast/lat,long` where `lat` and `long` are the latitude and longitude you'd like to get weather for.  

Example: `http://yourdomain.com:3000/forecast/34.0487043,-84.22674289999999`

To test your service quickly, you can use the [Postman Google Chrome Extension](https://chrome.google.com/webstore/detail/postman-rest-client/fdmmgilgnpjigdojojpjoooidkmcomcm?hl=en) to call the service and see the JSON return format.

Weather information will be returned as a [JSON payload outlined on the Forecast.io website](https://developer.forecast.io/docs/v2#forecast_call).
