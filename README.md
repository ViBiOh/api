# go-api

[![Build Status](https://travis-ci.org/ViBiOh/go-api.svg?branch=master)](https://travis-ci.org/ViBiOh/go-api)
[![codecov](https://codecov.io/gh/ViBiOh/go-api/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/go-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/go-api)](https://goreportcard.com/report/github.com/ViBiOh/go-api)

## Usage

```bash
Usage of api:
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
  -crudDefaultPage uint
      [crud] Default page (default 1)
  -crudDefaultPageSize uint
      [crud] Default page size (default 20)
  -crudMaxPageSize uint
      [crud] Max page size (default 500)
  -crudPath string
      [crud] HTTP Path prefix (default "/crud")
  -csp string
      [owasp] Content-Security-Policy (default "default-src 'self'; base-uri 'self'")
  -frameOptions string
      [owasp] X-Frame-Options (default "deny")
  -hsts
      [owasp] Indicate Strict Transport Security (default true)
  -location string
      [hello] TimeZone for displaying current time (default "Europe/Paris")
  -port int
      Listen port (default 1080)
  -prometheusPath string
      [prometheus] Path for exposing metrics (default "/metrics")
  -tls
      Serve TLS content (default true)
  -tlsCert string
      [tls] PEM Certificate file
  -tlsHosts string
      [tls] Self-signed certificate hosts, comma separated (default "localhost")
  -tlsKey string
      [tls] PEM Key file
  -tlsOrganization string
      [tls] Self-signed certificate organization (default "ViBiOh")
  -tracingAgent string
      [opentracing] Jaeger Agent (e.g. host:port) (default "jaeger:6831")
  -tracingName string
      [opentracing] Service name
  -url string
      [health] URL to check
  -userAgent string
      [health] User-Agent for check (default "Golang alcotest")
```
