# SAEONYX Caribbean Ocean Agent

Runnable FastAPI prototype for a Caribbean ocean coordination agent using live NOAA NDBC buoy data.

## What it does
- fetches live NOAA NDBC station data
- projects real observations into an 8D substrate state
- derives NORMAL / WATCH / ALERT / ESCALATE decisions
- serves a small Leaflet dashboard and JSON endpoints
- persists decisions, actions, and subscribers to SQLite
- exposes an independent prototype Great Salt Lake profile backed by public USGS NWIS observations

## Run

```bash
cd apps/caribbean-agent
python -m pip install -r requirements.txt
python main.py
```

Then open:

```text
http://localhost:7735
```

## Regional profile endpoints

- `GET /regions/great-salt-lake` — profile metadata, public sites, and supported parameters
- `GET /regions/great-salt-lake/status` — current normalized USGS snapshot and provenance digest

The Great Salt Lake profile keeps USGS gage height separate from absolute water-surface elevation, never fills missing measurements with synthetic values, and labels its output `prototype`. Set `GSL_ANALYSIS_REFERENCE_ELEVATION_FT` only when an analytical comparison reference is desired; it is not treated as a sensor reading or hazard threshold.

## Config
Set environment variables as needed:

- `HOST`
- `PORT`
- `DATABASE_URL`
- `POLL_INTERVAL`
- `ADMIN_TOKEN`

Example PowerShell:

```powershell
$env:ADMIN_TOKEN = "change-me"
python main.py
```

## Notes
- Uses real NOAA NDBC text feeds only.
- If a station is unavailable, it is marked offline instead of being simulated.
- The `/admin/status` endpoint requires the `admin-token` header and a configured `ADMIN_TOKEN` environment variable.
