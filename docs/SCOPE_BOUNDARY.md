# MVP vs Future-State Scope Boundary

## Purpose
This document makes the scope boundary explicit so AMOC Sentinel reads as both ambitious and believable.

The buildathon project has two architecture layers:
- an MVP architecture that is intentionally scoped for the sprint
- a future-state architecture that captures the larger platform vision

## MVP architecture
The MVP is the smallest credible system that can ingest ocean and observational data, reason over trusted signals, and produce operator-facing outputs.

### In scope for the buildathon MVP
- bounded ingestion of public data sources such as RAPID, Copernicus Marine, CariCOOS, and NOAA NDBC
- provenance, trust, and quality validation
- AMOC and thermal risk interpretation
- hazard translation for Caribbean-relevant outputs
- action briefs, advisories, or risk summaries
- a clean multi-agent workflow with explicit role separation

### Optional if time allows
- lightweight scenario simulation
- API-ready output packaging
- limited premium or machine-payable endpoint planning
- limited review workflow or expert-gated output path

## Future-state architecture
The future-state architecture is the broader platform vision beyond the 21-day sprint.

### Future-state elements
- deeper edge and IoT sensing networks
- packetization or compression layers
- broader simulation systems
- richer human approval and review flows
- stakeholder hubs and larger API ecosystems
- x402-style machine-payable distribution
- premium data, simulation, and operator products

## Why this distinction matters
This boundary helps the project communicate two things at once:
- the MVP is realistic and buildable
- the long-term platform vision is much bigger than the first sprint

That is the right posture for judges, mentors, and potential partners.

## Recommended communication line
Use this sentence whenever needed:

"This is our full future-state architecture. For the buildathon MVP, we are implementing the ingestion, trust/provenance, AMOC/current-state reasoning, hazard translation, and operator briefing layers first."

## Bottom line
AMOC Sentinel should present a disciplined MVP and an ambitious long-term architecture, without confusing the two.
