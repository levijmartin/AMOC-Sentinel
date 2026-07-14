# x402 Two-Phase Replay Risk Note

## Why this note exists

An NDSS 2026 poster, **"Exploiting the Two-Phase Gap in the x402 Protocol for Autonomous AI Payments,"** highlights a practical security risk in x402-style monetization flows.

The issue is not that x402 is unusable. The issue is that **naive implementations of the verify/settle split can allow replay-driven over-delivery of paid resources**.

## Core finding

The poster describes a **Two-Phase Gap** between:
- payment proof verification
- on-chain settlement confirmation

If a service delivers protected resources after verification but before settlement finality, and it does not enforce strict single-use proof consumption under concurrency, the system can become vulnerable to a **TOCTOU replay attack**.

In the reported black-box evaluation:
- seven public x402-enabled services were tested
- two showed replay vulnerability under synchronized parallel bursts
- multiple protected responses were returned for a single proof
- while settlement still finalized only once on-chain

## Why this matters for AMOC Sentinel

AMOC Sentinel has been treating x402 as a plausible future layer for:
- premium forecast endpoints
- pay-per-report access
- machine-payable intelligence products
- API/commercialization surfaces

That means AMOC should assume the following risk exists unless explicitly designed away:

- one proof may trigger multiple premium deliveries
- one paid request may fan out into repeated expensive work
- repeated premium outputs may be served before proof consumption is safely recorded

## Highest-risk AMOC cases

This matters most for **non-idempotent or costly premium actions**, such as:
- generated reports
- premium forecast runs
- paid agent queries
- credits/tokens/access grants
- expensive downstream processing triggered by a paid request

## Recommended design rules for AMOC

### 1. Atomic single-use consumption
Do not treat proof verification as sufficient on its own.

AMOC should atomically:
- verify the proof
- mark it consumed
- reject concurrent re-use

If those are not tied together safely, replay windows remain.

### 2. Strict idempotency
The same proof should map to the same logical request outcome.

If a duplicate request arrives with the same proof, the service should:
- reject it, or
- return the already-recorded outcome,

but **not** trigger a second billable delivery path.

### 3. Concurrency-safe proof tracking
Consumption state must be checked and updated in a concurrency-safe way.

This means avoiding designs where two requests can both observe:
- "proof not yet consumed"

and then both proceed.

### 4. Confirmation-aware provisioning
If possible, prefer delivery after stronger settlement assurance.

If AMOC chooses lower-latency pre-confirmation delivery, it should add compensating controls that prevent cumulative replay effects.

### 5. Extra caution for expensive endpoints
Premium endpoints that trigger expensive inference, report generation, or resource-heavy processing should be treated as especially sensitive.

Even a small replay window can multiply cost.

## Practical AMOC implementation takeaway

For AMOC, the premium-access rule should be:

**Do not just verify payment. Atomically consume payment proof before allowing repeated delivery paths.**

That is the key control that turns x402 from a monetization idea into a defensible production design.

## Scope note

This note is about:
- x402 monetization security
- concurrency and replay safety
- premium endpoint design

It is **not** a comment on:
- climate/ocean ingestion
- operator memory
- retrieval systems
- local persistence
- x402 packet/schema research outside the payment flow

## Bottom line

**AMOC Sentinel can still use x402-style monetization, but only if replay protection, atomic proof consumption, and idempotency are treated as first-class design requirements.**
