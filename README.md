# go-api

[![Build Status](https://travis-ci.org/ViBiOh/go-api.svg?branch=master)](https://travis-ci.org/ViBiOh/go-api)
[![codecov](https://codecov.io/gh/ViBiOh/go-api/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/go-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/go-api)](https://goreportcard.com/report/github.com/ViBiOh/go-api)

## Usage

```
Usage of api:
  -c string
      URL to check
  -corsHeaders string
      Access-Control-Allow-Headers (default "Content-Type")
  -corsMethods string
      Access-Control-Allow-Methods (default "GET")
  -corsOrigin string
      Access-Control-Allow-Origin (default "*")
  -csp string
      Content-Security-Policy (default "default-src 'self'")
  -hsts
      Indicate Strict Transport Security (default true)
  -prometheusMetricsPath string
      Prometheus - Metrics endpoint path (default "/metrics")
  -prometheusMetricsRemoteHost string
      Prometheus - Regex of allowed hosts to call metrics endpoint (default ".*")
```