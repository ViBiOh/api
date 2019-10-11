# api

[![Build Status](https://travis-ci.org/ViBiOh/api.svg?branch=master)](https://travis-ci.org/ViBiOh/api)
[![codecov](https://codecov.io/gh/ViBiOh/api/branch/master/graph/badge.svg)](https://codecov.io/gh/ViBiOh/api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/api)](https://goreportcard.com/report/github.com/ViBiOh/api)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=ViBiOh/api)](https://dependabot.com)

## Usage

```bash
Usage of api:
  -address string
        [http] Listen address {API_ADDRESS}
  -cert string
        [http] Certificate file {API_CERT}
  -corsCredentials
        [cors] Access-Control-Allow-Credentials {API_CORS_CREDENTIALS}
  -corsExpose string
        [cors] Access-Control-Expose-Headers {API_CORS_EXPOSE}
  -corsHeaders string
        [cors] Access-Control-Allow-Headers {API_CORS_HEADERS} (default "Content-Type")
  -corsMethods string
        [cors] Access-Control-Allow-Methods {API_CORS_METHODS} (default "GET")
  -corsOrigin string
        [cors] Access-Control-Allow-Origin {API_CORS_ORIGIN} (default "*")
  -crudDefaultPage uint
        [crud] Default page {API_CRUD_DEFAULT_PAGE} (default 1)
  -crudDefaultPageSize uint
        [crud] Default page size {API_CRUD_DEFAULT_PAGE_SIZE} (default 20)
  -crudMaxPageSize uint
        [crud] Max page size {API_CRUD_MAX_PAGE_SIZE} (default 500)
  -csp string
        [owasp] Content-Security-Policy {API_CSP} (default "default-src 'self'; base-uri 'self'")
  -frameOptions string
        [owasp] X-Frame-Options {API_FRAME_OPTIONS} (default "deny")
  -hsts
        [owasp] Indicate Strict Transport Security {API_HSTS} (default true)
  -key string
        [http] Key file {API_KEY}
  -location string
        [hello] TimeZone for displaying current time {API_LOCATION} (default "Europe/Paris")
  -port int
        [http] Listen port {API_PORT} (default 1080)
  -prometheusPath string
        [prometheus] Path for exposing metrics {API_PROMETHEUS_PATH} (default "/metrics")
  -tracingAgent string
        [tracing] Jaeger Agent (e.g. host:port) {API_TRACING_AGENT} (default "jaeger:6831")
  -tracingName string
        [tracing] Service name {API_TRACING_NAME}
  -url string
        [alcotest] URL to check {API_URL}
  -userAgent string
        [alcotest] User-Agent for check {API_USER_AGENT} (default "Golang alcotest")
```
