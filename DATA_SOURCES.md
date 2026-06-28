# Data Sources

## Recommended MVP dataset stack

For the first buildathon version, the project should stay focused on a small, credible data stack that supports climate-risk intelligence without becoming too heavy to build in 21 days.

## Core sources

### 1. RAPID Array
**Purpose:** AMOC and Atlantic overturning circulation context

Why use it:
- recognized benchmark dataset for AMOC behavior
- strong scientific credibility
- useful for long-range ocean circulation context and anomaly interpretation

Use in MVP:
- ingest summarized AMOC transport timeseries
- use as a background climate/ocean signal rather than a direct short-term weather feed

### 2. Copernicus Marine
**Purpose:** operational ocean and sea-surface intelligence

Why use it:
- excellent APIs and tooling
- strong support for SST, reanalysis, altimetry, and regional ocean-state products
- practical for automation and scheduled ingestion

Use in MVP:
- sea surface temperature anomalies
- regional ocean-state snapshots
- reanalysis context for Caribbean risk scoring

### 3. CariCOOS and NOAA NDBC
**Purpose:** Caribbean and local near-real-time observational data

Why use it:
- regionally relevant
- practical for real-world marine and weather observation inputs
- useful for current conditions and localized signal confirmation

Use in MVP:
- buoy observations
- wind / wave / SST context
- local observational grounding for resilience alerts

### 4. Argo Floats
**Purpose:** temperature and salinity profile context

Why use it:
- valuable for density, heat, and salinity trend analysis
- good for deeper climate context beyond surface-only signals

Use in MVP:
- optional regional subset
- anomaly context and research-grade signal support

## Supporting sources for later phases

### OSNAP
Use later for broader North Atlantic circulation context if the MVP expands into system-wide pattern modeling.

### World Ocean Database / World Ocean Atlas
Use later for historical baseline and anomaly comparisons.

## Data-layer strategy

The product should treat data in three layers:

### Layer 1: short-term operational awareness
- CariCOOS
- NDBC
- Copernicus near-real-time products

### Layer 2: climate and ocean pattern intelligence
- RAPID
- Argo
- Copernicus reanalysis
- later OSNAP / WOD / WOA

### Layer 3: local exposure and actionability
- property/business profile
- storm/flood exposure
- backup power / drainage / roof condition inputs
- business continuity readiness information

## MVP guidance

Keep the first prototype centered on:
- RAPID
- Copernicus Marine
- CariCOOS / NDBC
- optional lightweight Argo subset

This is enough to support a believable agentic workflow without overbuilding the ingestion pipeline.
