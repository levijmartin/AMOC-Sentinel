# SAEONYX Caribbean Ocean Agent

Runnable FastAPI prototype for a Caribbean ocean coordination agent using live NOAA NDBC buoy data.

## What it does
- fetches live NOAA NDBC station data
- projects real observations into an 8D substrate state
- derives NORMAL / WATCH / ALERT / ESCALATE decisions
- serves a small Leaflet dashboard and JSON endpoints
- persists decisions, actions, and subscribers to SQLite

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
