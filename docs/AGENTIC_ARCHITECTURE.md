# Agentic Architecture

## Purpose
This document sharpens AMOC Sentinel's agentic system design so the project reads as a true multi-agent intelligence workflow rather than a dashboard with AI attached.

## Core principle
Each agent should have a narrow responsibility, a clear input/output contract, and a reason to exist independently. The value of the system comes from how agents collaborate across ingestion, validation, interpretation, hazard translation, simulation, and delivery.

## Primary agent roles

### 1. Source Ingestion Agent
Responsibilities:
- pull data from RAPID, Copernicus Marine, Argo, CariCOOS, NDBC, and related sources
- schedule or trigger ingestion jobs
- normalize metadata about source freshness and availability

Outputs:
- raw datasets
- source status metadata
- ingestion logs

### 2. Provenance and Trust Agent
Responsibilities:
- verify source identity and freshness
- assign confidence and trust labels
- track lineage from source to downstream output
- flag incomplete or conflicting data

Outputs:
- trusted signal set
- provenance records
- confidence labels

### 3. Ocean State Agent
Responsibilities:
- evaluate AMOC, Cold Blob, SST, salinity, and related ocean indicators
- compute anomalies and current-state summaries
- maintain the current ocean condition picture for the region

Outputs:
- ocean state summaries
- anomaly indicators
- state-change flags

### 4. Hazard Translation Agent
Responsibilities:
- convert ocean and climate signals into Caribbean-relevant hazards
- map environmental state into categories such as storm, marine heat, fisheries stress, sargassum risk, sea-level stress, or port disruption

Outputs:
- hazard classifications
- regional or sector-specific risk indicators

### 5. Operator Context Agent
Responsibilities:
- incorporate local exposure details
- understand whether the user is a port, fishery, hospitality operator, coastal facility, or resilience stakeholder
- attach context-specific operational consequences to hazards

Outputs:
- localized risk context
- stakeholder-specific impact framing

### 6. Action Brief Agent
Responsibilities:
- generate operator-facing advisories and action briefs
- prioritize what matters now vs later
- produce concise, usable output rather than technical raw data

Outputs:
- briefs
- advisories
- action summaries

### 7. Scenario Simulation Agent
Responsibilities:
- explore what-if response paths
- simulate possible operational outcomes under selected hazard conditions
- provide alternate decision paths when confidence is high enough

Outputs:
- scenario summaries
- response comparisons
- simulation-backed recommendations

### 8. Human Review Agent / Approval Gate
Responsibilities:
- intercept high-risk or low-confidence outputs
- require expert review before escalation when appropriate
- maintain trust in the system for sensitive decisions

Outputs:
- approved outputs
- held outputs
- review notes

### 9. API and Commerce Agent
Responsibilities:
- package premium outputs for API or stakeholder delivery
- support monetizable surfaces such as premium briefs, alerts, reports, and machine-payable endpoints
- connect with x402-style access controls when the product is ready

Outputs:
- API responses
- billable intelligence products
- premium delivery events

## Agent handoff chain

```text
Source Ingestion Agent
  -> Provenance and Trust Agent
  -> Ocean State Agent
  -> Hazard Translation Agent
  -> Operator Context Agent
  -> Action Brief Agent
  -> Scenario Simulation Agent (optional branch)
  -> Human Review Agent (conditional)
  -> API and Commerce Agent
```

## Why this is stronger than a dashboard
A dashboard displays data. An agentic system:
- validates trust
- reasons across multiple sources
- adapts outputs to the stakeholder
- branches into simulation or review when needed
- delivers decisions, not just information

## Buildathon framing
This architecture is useful in a buildathon because it shows:
- clear multi-agent decomposition
- defensible workflow design
- real-world trust and review handling
- a believable path from open data to operator action
- a commercial layer without making the MVP payment-heavy from day one

## MVP subset
The full architecture is modular, but the first MVP should focus on:
- Source Ingestion Agent
- Provenance and Trust Agent
- Ocean State Agent
- Hazard Translation Agent
- Action Brief Agent

Optional if time allows:
- Scenario Simulation Agent
- Human Review Agent
- API and Commerce Agent

## Bottom line
AMOC Sentinel should present itself as a coordinated intelligence workflow where specialized agents convert fragmented ocean signals into trusted, localized, operator-facing decisions.
