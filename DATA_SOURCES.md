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

Practical note:
- the RAPID download bundle appears to include NetCDF, Matlab, ASCII, vertical transport, grid timeseries, uncertainty-related files, meridional transports, and 2D gridded data
- for the MVP, start with summarized or selected timeseries before taking on heavier processing

### 2. Copernicus Marine
**Purpose:** operational ocean and sea-surface intelligence

General Copernicus access reference:
- https://www.copernicus.eu/en/access-data

Why use it:
- excellent APIs and tooling
- strong support for SST, reanalysis, altimetry, and regional ocean-state products
- practical for automation and scheduled ingestion

Use in MVP:
- sea surface temperature anomalies
- regional ocean-state snapshots
- reanalysis context for Caribbean risk scoring

#### Copernicus starter products for the MVP
Start narrow. For the first buildathon version, prioritize a small group of product types rather than trying to ingest everything at once.

Recommended first pulls:
- **Sea Surface Temperature (SST)**
  - use for warm-water conditions, anomaly detection, and marine heat context
- **Regional or global ocean reanalysis / ocean state**
  - use for broader Caribbean environmental context and background state reasoning
- **Sea level anomaly / altimetry**
  - use for circulation and coastal context where relevant
- **Currents**
  - use selectively if the demo needs movement or flow context

Suggested MVP role for each:
- SST -> near-surface heat and storm-supportive condition context
- Reanalysis -> broader background ocean conditions
- Sea level anomaly -> circulation and coastal signal support
- Currents -> optional second-phase enhancement

Practical ingestion guidance:
- begin with one SST product and one reanalysis product
- subset to a Caribbean bounding box to keep the prototype lightweight
- store derived anomaly summaries instead of only raw full-resolution grids
- use scheduled pulls rather than fully continuous ingestion for the MVP

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

### AMOC Atlas (`AMOCcommunity/amocatlas`)
Use as a research/reference source for AMOC-related analysis patterns, scientific framing, and possible visualization ideas. Best treated as an upstream reference, not a core MVP dependency.

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
