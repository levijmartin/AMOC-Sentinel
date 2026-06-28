# Implementation Notes

## Purpose
These notes capture practical engineering guidance for building StormShield Caribbean as an agentic climate-risk system without overusing language models for tasks better handled by standard code and data tooling.

## Core principle
Use conventional software and scientific computing tools for data ingestion, parsing, normalization, quality checks, and numerical calculations. Use agentic and language-model layers for interpretation, prioritization, explanation, and action planning.

## Recommended split of responsibilities

### Standard code should handle
- downloading data from Copernicus, RAPID, CariCOOS, NDBC, and related sources
- parsing NetCDF, CSV, JSON, or other structured data formats
- subsetting Caribbean geographic regions
- anomaly calculations and summary metrics
- threshold checks and rule-based triggers
- retries, caching, scheduling, and durable state management

### Agent / LLM layers should handle
- deciding which signals matter most to the user
- turning data summaries into understandable risk narratives
- asking follow-up questions about local exposure
- prioritizing resilience actions
- generating operator-facing reports, alerts, and action plans

## Do not send giant raw datasets into an LLM
Large raw tables or gridded datasets should not be passed directly into the reasoning layer. Instead:
- compute compact summaries first
- extract anomaly indicators and threshold events
- pass only the most relevant structured context into the agent workflow

## Suggested workflow state
Define a clean structured state object for the pipeline.

Possible state fields:
- request parameters
- region or island selection
- source datasets requested
- download status
- local file paths or cached objects
- normalized metrics
- anomaly summaries
- risk scores
- alert/report generation status
- user profile or property exposure inputs

## Reliability and resiliency
Expect external data systems to be imperfect.

Plan for:
- slow API responses
- download failures
- partial data availability
- retries with backoff
- caching of latest successful pulls
- durable checkpointing between workflow stages

## Data pipeline guidance
For the MVP:
- start with bounded Caribbean subsets
- prefer one or two strong signals over too many weak ones
- store reduced outputs for agent use rather than only raw data
- keep the ingestion path transparent and debuggable

## Agent design guidance
The agentic system should not pretend to be doing scientific computation directly. It should operate on prepared summaries and structured outputs from the data-processing layer.

That makes the system:
- more reliable
- easier to debug
- easier to explain
- cheaper to run
- more credible in a buildathon demo

## Recommended stack direction
- Python for ingestion, processing, and scientific data handling
- pandas / xarray / netCDF4 for climate and ocean data
- lightweight storage and caching for MVP
- agent layer for explanation, planning, and user-facing outputs

## Bottom line
StormShield Caribbean should behave like a disciplined data pipeline with an agentic decision layer on top, not like a language model trying to do raw climate/ocean computation by itself.
