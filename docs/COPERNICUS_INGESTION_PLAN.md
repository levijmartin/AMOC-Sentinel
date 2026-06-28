# Copernicus Ingestion Plan

## Goal
Use Copernicus Marine as a primary ocean-intelligence layer for StormShield Caribbean without overbuilding the MVP.

## Recommended first product categories

### 1. Sea Surface Temperature (SST)
Why:
- easiest high-value signal for climate and storm-context interpretation
- useful for anomaly detection and marine heat context
- understandable to judges and users

Use in MVP:
- track Caribbean SST anomaly conditions
- flag elevated heat or unusual warm-water patterns
- support climate-risk summaries

### 2. Ocean Reanalysis / Ocean State
Why:
- gives broader background ocean conditions
- useful for contextualizing localized risk
- stronger system-wide signal layer than surface-only data

Use in MVP:
- provide regional baseline and anomaly context
- support background climate/ocean reasoning

### 3. Sea Level Anomaly / Altimetry
Why:
- useful for circulation context and some coastal reasoning
- supports richer ocean-state interpretation

Use in MVP:
- optional but useful enhancement
- include if team capacity allows

### 4. Currents
Why:
- can strengthen advanced ocean-state interpretation
- useful if later versions need movement/transport logic

Use in MVP:
- optional
- lower priority than SST and reanalysis

## Recommended implementation order

### Phase 1
- one SST product
- one ocean reanalysis product

### Phase 2
- sea level anomaly product
- optional current field product

## Geographic scope
Keep the initial scope bounded to a Caribbean region to stay lightweight.

Suggested logic:
- use a Caribbean bounding box
- later add island-specific subregions for focused analysis

## Output strategy
Do not treat Copernicus as just a raw-data dump.

Instead, transform product outputs into:
- anomaly summaries
- regional risk indicators
- simplified climate/ocean context objects
- user-readable explanations

## Agent workflow fit

### Ingestion agent
- fetch Copernicus subsets on schedule
- cache latest snapshots

### Signal analysis agent
- compute SST and ocean-state anomalies
- identify elevated environmental conditions

### Climate-risk reasoning agent
- combine Copernicus signals with local exposure inputs
- decide what matters for the user

### Action-planning agent
- turn elevated signals into preparation guidance

## MVP success criteria
A successful Copernicus-backed MVP should be able to:
- pull a bounded Caribbean dataset reliably
- compute at least one useful anomaly or summary
- feed that result into the climate-risk reasoning workflow
- explain the result in plain English to a user
