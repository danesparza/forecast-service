# forecast-service
A pass-through microservice for the [forecast.io](http://forecast.io/) API, using [mlbright/forecast](https://github.com/mlbright/forecast)

[![Build Status](https://drone.io/github.com/danesparza/forecast-service/status.png)](https://drone.io/github.com/danesparza/forecast-service/latest)

Command line:
`-apikey="YOURAPIKEYHERE"`
`-port=3000`
`-allowedOrigins="*"`

`apikey` is the [Forecast.io](https://developer.forecast.io/docs/v2) API key you'll need to get forecast information.  

`port` is the port number you'd like the service to listen on. 

`allowedOrigins` is the comma seperated list of [CORS](http://en.wikipedia.org/wiki/Cross-origin_resource_sharing) origins to allow

Once the service is up and running, you can connect to it using
`http://yourhostname:3000/forecast/lat,long` where `lat` and `long` are the latitude and longitude you'd like to get weather for.  Weather information will be returned as a [JSON payload outlined on the Forecast.io website](https://developer.forecast.io/docs/v2#forecast_call).
