# Idea Canvas

## Project name
AMOC Sentinel

## Track
Climate Risk

## One-line pitch
AMOC Sentinel is a cold-blob-to-Caribbean risk operating system that tracks ocean and weather shifts, verifies data quality and provenance, tokenizes trusted data packets, and delivers actionable risk intelligence for Caribbean operators.

## Who is the user?
Primary users:
- small property owners in hurricane-prone Caribbean markets
- small business operators exposed to storm, flood, and outage risk
- hospitality operators such as guesthouses, villas, and small lodging businesses

Secondary users:
- local operators responsible for facilities or continuity planning
- resilience-focused partners such as contractors, insurers, lenders, or advisors

## What painful problem are we solving?
Caribbean operators are exposed to storm, flood, outage, and coastal risk, but the intelligence chain is fragmented. Ocean and weather shifts, AMOC-related context, local observations, and operator decisions are spread across disconnected systems.

Today the problem is fragmented:
- ocean and weather information is scattered across technical sources
- operators receive data, but not trusted, verified decisions
- there is limited linkage between Atlantic-scale context and Caribbean-local action
- most users do not have a practical, phased resilience roadmap tied to their actual operating risk

## Why is this a Caribbean opportunity?
The Caribbean faces recurring hurricane, flooding, coastal, and infrastructure risk, while many communities and operators remain underserved by practical resilience intelligence tools.

This is a strong regional opportunity because:
- climate risk is immediate and economically meaningful
- island businesses and coastal operators are highly exposed
- public datasets already exist, but are not translated into operator-ready guidance
- a solution can create both public-value and commercial-value outcomes

## Why does this need agentic AI?
This is not just a dashboard problem. The system needs to:
- gather environmental and ocean data from multiple public sources
- track ocean and weather shifts across Atlantic and Caribbean layers
- verify data quality, provenance, and confidence before acting on it
- convert trusted signals into compact, structured intelligence objects
- prioritize vulnerabilities, thresholds, and tradeoffs
- generate phased action plans based on the user’s actual situation

An agentic workflow is useful because the system must move from raw signals to verified, user-specific decisions rather than just present charts.

## Core workflow
Describe the agent loop step by step.

1. Ingestion agent pulls data from sources such as RAPID, Copernicus Marine, CariCOOS, NOAA NDBC, and optional Argo subsets.
2. Validation and provenance agent cleans, aligns, signs, and structures environmental signals for Caribbean use.
3. Signal analysis agent computes anomalies, AMOC/current-state context, and relevant climate/ocean conditions.
4. Risk-intelligence agent combines those signals with user-provided property, business, or local operator details.
5. Action-planning agent generates immediate, near-term, and longer-term resilience actions.
6. Reporting and delivery layer produces practical readiness summaries, action briefs, alerts, and API-ready intelligence outputs.

## MVP for the 21-day sprint
- bounded Caribbean data ingestion using a small set of public sources
- trusted signal intake with validation, provenance, and confidence flags
- simplified environmental signal summaries and anomaly context
- AMOC/current-state-informed risk reasoning
- phased resilience action plans and action briefs
- exportable readiness report, alert summary, or API-ready output

## Success metric
A user can complete intake in under 10 minutes and receive a useful, understandable, prioritized resilience action plan tied to both environmental context and local vulnerability.

## Demo scenario
A Caribbean operator such as a guesthouse owner, port-adjacent business, or local facility manager enters location and exposure details.

AMOC Sentinel then:
- pulls relevant Atlantic, ocean, and Caribbean environmental context
- verifies the trusted signal set and confidence level
- identifies elevated vulnerabilities and threshold risks
- produces a phased plan showing what to do now, what to do before the next storm season, and what longer-term resilience investments to prioritize

## Data sources / APIs
Core sources:
- RAPID Array
- Copernicus Marine
- CariCOOS
- NOAA NDBC

Optional or later-stage supporting sources:
- Argo
- OSNAP
- World Ocean Database / World Ocean Atlas
- AMOC Atlas as an upstream research and scientific framing reference

## Risks / blockers
- balancing scientific credibility with MVP simplicity
- keeping recommendations understandable and actionable
- avoiding over-engineering the ingestion layer during the sprint
- limited time for deep calibration across multiple islands and exposure types
- ensuring the system gives decision support, not false precision

## Why this can become a company
AMOC Sentinel can evolve into a climate-risk and resilience intelligence platform for Caribbean operators.

Possible expansion paths:
- premium action briefs, dashboards, and alert subscriptions
- insurer, lender, advisor, and public-sector partnerships
- API and data sales for researchers, AI models, and operators
- sensor bounty or node-network models for trusted local observations
- x402-enabled per-report, per-alert, or per-query access

The long-term value is not just climate monitoring. The company opportunity is turning environmental complexity into trusted, actionable operational intelligence.
