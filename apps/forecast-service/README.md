# Forecast Service

Simple Go-based weather/risk forecast form app for AMOC Sentinel.

## What it does
- serves a basic HTML form for forecast lookup
- exposes a free JSON endpoint at `/api/forecast`
- exposes an x402-style premium test endpoint at `/api/premium/forecast`
- simulates premium access using the `X-X402-Payment` header

## Run locally

```bash
cd apps/forecast-service
go run .
```

Then open:

```text
http://localhost:8080
```

## Test free API

```bash
curl "http://localhost:8080/api/forecast?location=Bridgetown&region=Barbados&date=2026-07-01"
```

## Test premium x402-style API

```bash
curl -H "X-X402-Payment: test-paid" "http://localhost:8080/api/premium/forecast?location=Bridgetown&region=Barbados&date=2026-07-01"
```

## Notes
This is a scaffold for testing the AMOC Sentinel service shape and monetization boundary. It uses mock forecast logic right now, not real Copernicus or NOAA ingestion yet.
