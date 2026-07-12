# Forecast Service

Simple Go-based weather/risk forecast form app for AMOC Sentinel.

## What it does
- serves a basic HTML form for forecast lookup
- exposes a free JSON endpoint at `/api/forecast`
- exposes an x402-style premium test endpoint at `/api/premium/forecast`
- exposes a placeholder operator-memory endpoint at `/api/memory/operator`
- simulates premium access using the `X-X402-Payment` header
- includes a lightweight in-memory `MemoryStore` interface that can later be swapped for mem0

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

The response now includes:
- `rememberedContext` when memory exists
- `memoryImplementation` showing the active memory backend

## Test premium x402-style API

```bash
curl -H "X-X402-Payment: test-paid" "http://localhost:8080/api/premium/forecast?location=Bridgetown&region=Barbados&date=2026-07-01"
```

## Test placeholder memory API

Save operator context:

```bash
curl -X POST "http://localhost:8080/api/memory/operator" \
  -H "Content-Type: application/json" \
  -d '{
    "location": "Bridgetown",
    "region": "Barbados",
    "facilityType": "guesthouse",
    "preferredAlertStyle": "short alerts",
    "knownVulnerabilities": ["poor drainage", "no backup generator"]
  }'
```

Read operator context:

```bash
curl "http://localhost:8080/api/memory/operator?location=Bridgetown&region=Barbados"
```

## Notes
This is still a scaffold for testing the AMOC Sentinel service shape and monetization boundary. It uses mock forecast logic right now, not real Copernicus or NOAA ingestion yet.

The memory layer is intentionally a placeholder. It demonstrates where persistent operator context fits, while keeping scientific data and telemetry out of the memory subsystem.
