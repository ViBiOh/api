# go-api

[![Build Status](https://travis-ci.org/ViBiOh/go-api.svg?branch=master)](https://travis-ci.org/ViBiOh/go-api)
[![codecov](https://codecov.io/gh/ViBiOh/go-api/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/go-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/go-api)](https://goreportcard.com/report/github.com/ViBiOh/go-api)

## Usage

```
Usage of api:
  -c string
    	[health] URL to check
  -corsCredentials
    	[cors] Access-Control-Allow-Credentials
  -corsExpose string
    	[cors] Access-Control-Expose-Headers
  -corsHeaders string
    	[cors] Access-Control-Allow-Headers (default "Content-Type")
  -corsMethods string
    	[cors] Access-Control-Allow-Methods (default "GET")
  -corsOrigin string
    	[cors] Access-Control-Allow-Origin (default "*")
  -csp string
    	[owasp] Content-Security-Policy (default "default-src 'self'")
  -hsts
    	[owasp] Indicate Strict Transport Security (default true)
  -location string
    	TimeZone for displaying current time (default "Europe/Paris")
  -port string
    	Listen port (default "1080")
  -tls
    	Serve TLS content
  -tlsCert string
    	[tls] PEM Certificate file
  -tlsHosts string
    	[tls] Self-signed certificate hosts, comma separated (default "localhost")
  -tlsKey string
    	[tls] PEM Key file
```