# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# To configure the app, set these environment variables or use the command line flags
ENV FORECAST_ALLOWED_ORIGINS *
ENV FORECAST_KEY YOUR_API_KEY
ENV FORECAST_EXPVAR_USER YOUR_USERNAME
ENV FORECAST_EXPVAR_PASS YOUR_PASSWORD

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/danesparza/forecast-service

# Build and install the app inside the container.
RUN go get github.com/danesparza/forecast-service/...

# Run the app by default when the container starts.
ENTRYPOINT /go/bin/forecast-service

# Document that the app listens on port 3000.
EXPOSE 3000