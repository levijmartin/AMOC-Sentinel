# MVP Architecture

## Project name
StormShield Caribbean

## Technical concept
An agentic climate-risk intelligence system that combines ocean and regional observational data with local property or business inputs to generate practical resilience guidance.

## System goal
Turn complex Atlantic and Caribbean environmental signals into understandable, prioritized action plans for property owners and small businesses.

## High-level architecture

```text
Data Sources
  -> Ingestion Agent
  -> Validation / Normalization Agent
  -> Signal Analysis Agent
  -> Climate Risk Reasoning Agent
  -> Action Planning Agent
  -> Report / Alert Output
```

## Agents

### 1. Ingestion Agent
Responsibilities:
- fetch RAPID timeseries
- fetch Copernicus SST and regional ocean-state data
- fetch CariCOOS / NDBC observational data
- optionally fetch a lightweight Argo subset

Outputs:
- normalized time series and snapshots
- cached raw source data

### 2. Validation / Normalization Agent
Responsibilities:
- clean missing or malformed data
- standardize timestamps and units
- align regional subsets to project geography
- create simplified internal data objects

Outputs:
- normalized climate/ocean observations
- source confidence flags

### 3. Signal Analysis Agent
Responsibilities:
- compute SST anomalies
- identify risk-relevant ocean signal changes
- summarize local observational conditions
- detect abnormal conditions against baseline assumptions

Outputs:
- machine-readable signal summaries
- anomaly and trend indicators

### 4. Climate Risk Reasoning Agent
Responsibilities:
- combine environmental signals with user-provided exposure details
- assess likely business or property vulnerabilities
- translate scientific inputs into understandable risk categories

Inputs:
- location / island / region
- property or business type
- roof, drainage, power, opening protection, exposure details

Outputs:
- wind / flood / outage / continuity risk scores
- top vulnerabilities
- reasoning summary

### 5. Action Planning Agent
Responsibilities:
- prioritize next best actions
- split recommendations into immediate, near-term, and long-term
- adapt recommendations to budget or readiness constraints

Outputs:
- prioritized resilience roadmap
- checklist of recommended actions
- phased upgrade plan

### 6. Report / Alert Output
Responsibilities:
- generate a user-readable resilience report
- provide short summary alerts
- export a simple shareable document for owners, operators, contractors, or investors

Outputs:
- dashboard summary
- downloadable report
- alert summaries

## Suggested MVP stack

### Front end
- simple web interface
- form-based intake for property or business details
- report/results page

### Back end
- Python or TypeScript service
- scheduled ingestion jobs
- lightweight API layer

### Data tools
- pandas
- xarray
- netCDF4
- Copernicus Marine package
- argopy if Argo is included in MVP

### Storage
- lightweight database or cached file store
- JSON / CSV / SQLite acceptable for MVP

## MVP workflow

1. User enters property or business details
2. System fetches or uses latest climate/ocean signals
3. Agents compute relevant anomalies and contextual risk
4. Risk reasoning agent scores vulnerabilities
5. Action planning agent generates a phased plan
6. User receives a resilience report

## Recommended demo scenario

A Caribbean guesthouse owner enters:
- island/location
- flood exposure
- roof type
- drainage quality
- storm shutters status
- backup power availability

The system returns:
- environmental context summary
- vulnerability highlights
- immediate preparedness actions
- 30 to 90 day resilience plan
- longer-term hardening roadmap

## MVP success criteria

- user can complete intake in under 10 minutes
- system returns a prioritized resilience plan
- recommendations are understandable and actionable
- demo clearly shows agentic workflow, not just static reporting
