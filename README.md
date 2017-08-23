# go-api

[![Build Status](https://travis-ci.org/ViBiOh/go-api.svg?branch=master)](https://travis-ci.org/ViBiOh/go-api)
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
  -prometheusMetricsIP string
      Prometheus - Allowed regex IP to call metrics endpoint (default "*")
  -prometheusMetricsPath string
      Prometheus - Allowed regex IP to call metrics endpoint (default "/metrics")
```