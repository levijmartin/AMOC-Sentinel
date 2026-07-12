# Mem0 Integration Plan for AMOC Sentinel

## Goal

Add persistent operator/context memory without mixing it into the scientific data pipeline.

## Scope

### In scope
- operator profile memory
- user preference memory
- prior recommendation summaries
- open follow-up items
- optional shared agent context retrieval

### Out of scope
- raw climate/ocean data storage
- provenance logs
- x402 packet authority
- large telemetry history

## Recommended architecture placement

```text
Data Sources
  -> Ingestion / Validation / Analysis
  -> Structured risk outputs + provenance records
  -> Agent workflow
       -> mem0 retrieval (operator context)
       -> reasoning / action brief generation
       -> mem0 write-back (new stable facts, action summaries)
```

## First implementation slice

### Read before advisory generation
Retrieve:
- operator profile
- location
- business/facility type
- known vulnerabilities
- prior recommendations
- unresolved next steps

### Write after advisory generation
Store only stable, useful memory such as:
- updated operator profile facts
- constraints the user stated
- concise summary of recommendations delivered
- actions the user says are complete
- actions still pending

## Example memory payloads

### Operator profile
```json
{
  "type": "operator_profile",
  "operator_id": "guesthouse-001",
  "region": "Barbados",
  "facility_type": "guesthouse",
  "vulnerabilities": ["poor drainage", "no backup generator"],
  "preferences": ["short alerts", "practical action lists"]
}
```

### Recommendation history summary
```json
{
  "type": "report_history_summary",
  "operator_id": "guesthouse-001",
  "summary": "Previous brief emphasized drainage clearing, shutter inspection, and generator planning.",
  "open_actions": ["obtain generator quotes", "clear roof drains monthly"]
}
```

## Integration notes

- Keep mem0 writes concise and selective.
- Do not write every model output verbatim.
- Summarize before storing.
- Keep canonical reports in normal application storage.
- Use mem0 retrieval only as supporting context for agent output generation.

## Suggested next repo step

When implementation starts, add:
- a memory service wrapper
- clear schemas for stored memory types
- a toggle so mem0 is optional in MVP deployments

## Bottom line

Use mem0 to remember the operator.
Do not use mem0 to remember the ocean.
