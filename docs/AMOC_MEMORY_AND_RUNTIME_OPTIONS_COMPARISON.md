# AMOC Memory and Runtime Options Comparison

## Purpose

This short note compares four options already discussed for AMOC Sentinel:
- mem0
- FlashDB
- Needle
- Metauto Neural Computer

These do **not** solve the same problem, so the right question is not "which one wins overall?" but **which one fits which layer**.

## Quick summary

### mem0
Best for:
- operator/user memory
- preference recall
- advisory continuity
- shared agent context retrieval

Not best for:
- raw climate/ocean data
- provenance logs
- packet authority
- primary analytical storage

### FlashDB
Best for:
- lightweight embedded persistence
- forecast-service local state
- cached API responses
- simple operator profile records
- app-level durability without a separate DB service

Not best for:
- large scientific datasets
- rich analytics
- relational reporting complexity
- platform-scale backbone storage

### Needle
Best for:
- document retrieval
- repo knowledge grounding
- context-aware agent responses
- searching reference material and implementation notes
- internal knowledge-assistant workflows

Not best for:
- primary app persistence
- operator memory by itself
- scientific data backbone storage
- deterministic packet/x402 execution logic

### Metauto Neural Computer
Best for:
- future-state inspiration
- autonomous workflow ideas
- shared-memory/orchestration concepts
- long-term agent runtime thinking

Not best for:
- immediate MVP implementation
- concrete persistence needs
- explicit provenance/state storage
- practical near-term app plumbing

## Best fit by AMOC layer

### Operator memory layer
**Best fit: mem0**

Why:
- purpose-built for persistent user/operator memory
- helps agents remember context across sessions
- aligns with Action Brief and Operator Context use cases

### Forecast-service persistence layer
**Best fit: FlashDB**

Why:
- lightweight embedded durability
- simpler than standing up a separate database service for the MVP
- clean fit behind the current placeholder `MemoryStore` interface

### Knowledge retrieval / grounding layer
**Best fit: Needle**

Why:
- useful for retrieving the right AMOC docs, schemas, and notes at response time
- helps ground agent answers in repo material
- strong fit for internal knowledge-assistant and documentation-navigation use cases

### Scientific data / ingestion backbone
**Best fit: none of mem0, FlashDB, Needle, or Metauto as the primary answer**

Better direction:
- normal files/object storage
- SQLite/Postgres where needed
- explicit data pipeline code

### Long-term orchestration vision
**Best fit: Metauto as inspiration**

Why:
- strongest value is conceptual
- useful for thinking about shared state and more autonomous agent workflows
- not yet the most concrete delivery tool for this repo

## Recommended near-term combination

For the AMOC MVP, the cleanest working split is:

- **mem0** for operator/user memory
- **FlashDB** for lightweight local app persistence
- **Needle** for knowledge retrieval and document grounding
- **standard code + conventional storage** for climate/ocean data and workflow state
- **Metauto ideas** only as future architecture inspiration

## Decision shorthand

If the question is:

- **"How do we remember the operator?"** -> mem0
- **"How do we persist local app state simply?"** -> FlashDB
- **"How do we retrieve the right docs and reference knowledge?"** -> Needle
- **"How do we store climate/ocean data?"** -> conventional data storage, not these
- **"How do we think about future autonomous agent architecture?"** -> Metauto

## Bottom line

These options are complementary more than competing:

- mem0 remembers the user
- FlashDB persists the app
- Needle retrieves the knowledge
- Metauto informs the future vision

For AMOC Sentinel right now, the practical path is:
**build with mem0 + FlashDB + Needle-style grounding where useful, and treat Metauto as inspiration rather than dependency.**
