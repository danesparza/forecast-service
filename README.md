# forecast-service
A pass-through microservice for Forecast.io, using [mlbright/forecast](https://github.com/mlbright/forecast)

Command line:
`--apikey="YOURAPIKEYHERE"`
`--port=3000`

API key is the [Forecast.io](https://developer.forecast.io/docs/v2) API key you'll need to get forecast information.  

Port is the port number you'd like the service to listen on. 

Once the service is up and running, you can connect to it using
`http://yourhostname:3000/forecast/lat,long` where `lat` and `long` are the latitude and longitude you'd like to get weather for.  Weather information will be returned as a [JSON payload outlined on the Forecast.io website](https://developer.forecast.io/docs/v2#forecast_call).
