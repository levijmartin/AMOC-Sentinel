# SAEONYX Integration Notes

## Source and evidence status

Source reviewed: `SAEONYX_Full_Technical_Operational_Monograph_DEFINITIVE_CODE_GROUNDED.docx`, prepared 2026-09-01.

The monograph describes `SX_DEFINITIVE_FRESH.zip` with SHA-256 `986b02e01800862c772404319e4736a7acd76a6bc35ed302175ee0da1815f80c`, a 146-member manifest, 113 Python files, eight test modules, and 29 explicit tests. That source archive was not supplied with the document, so AMOC must treat its implementation claims as **documented assertions, not independently verified code evidence**.

Use four capability labels in AMOC documentation and APIs:

- **active** — reachable from the current AMOC production/demo entrypoint and covered by a current test or live check;
- **prototype** — executable in the repository but not production-qualified;
- **retained/reference** — code or design material exists but is not on the active path;
- **theory** — conceptual language with no validated operational implementation.

Do not promote a capability merely because the monograph or a repository file names it.

## Reusable architectural patterns

### 1. Separate runtime authority from capability delegation

The monograph distinguishes one runtime authority from a surrounding specialist-agent fabric. AMOC should reuse the boundary, not duplicate the described 62-agent system.

Recommended AMOC interpretation:

- deterministic ingestion, validation, decoding, thresholds, and persistence remain authoritative code paths;
- specialist or language-model components may research, summarize, explain, and propose actions;
- specialist output cannot silently replace sensor evidence, decoder output, or provenance records;
- every delegated result identifies the responsible component and source evidence.

### 2. Make trajectory evidence first-class

The described GTC pattern records whether a process remained inside an admissible corridor, rather than judging only its final state. AMOC can apply this to environmental observations and Sentinel control work:

- record source, observed time, ingest time, units, transforms, thresholds, and decision path;
- preserve each validation failure and fallback decision;
- distinguish stale, partial, offline, simulated, and live inputs;
- generate a digest over the normalized input and decision envelope;
- never allow a final green status to erase an earlier invalid transition.

GTC/UEF terminology may label an experimental governance view, but physical alerts must remain driven by validated measurements and explicit thresholds.

### 3. Reuse exact Tryte9 invariants

The monograph states the same exact Tryte9 invariant already implemented in AMOC:

- nine balanced trits use values `{-1, 0, +1}`;
- the state space contains exactly `3^9 = 19,683` values;
- the signed range is `-9,841` through `+9,841`;
- encode/decode must round-trip losslessly.

AMOC's current Go implementation and exhaustive round-trip test remain the authority for the demo. Neural-state, qutrit, entropy, or consciousness claims from the monograph do not extend that implementation automatically.

### 4. Preserve provenance and durable consequences

The monograph's strongest reusable pattern is transactional provenance. AMOC should define one evidence envelope shared by regional environmental profiles and future Sentinel components:

```json
{
  "event_id": "...",
  "component": "great-salt-lake|caribbean|sentinel-decode",
  "capability_status": "active|prototype|retained|theory",
  "source_uri": "...",
  "source_observed_at": "...",
  "ingested_at": "...",
  "freshness": "live|stale|partial|offline|simulated",
  "input_digest": "sha256:...",
  "transform_version": "...",
  "decision": "...",
  "decision_basis": ["..."],
  "output_digest": "sha256:..."
}
```

For the current SQLite prototype, append-only event rows plus a hash chain are sufficient. PostgreSQL, byte-level model provenance, and a separate ledger service are future-state options, not MVP requirements.

### 5. Never fabricate absent sensor input

The monograph explicitly describes sensory snapshots as expiring and not inventing values when hardware is absent. Apply the same rule to AMOC data adapters:

- missing NOAA, USGS, RAPID, Copernicus, or local sensor values remain missing;
- test fixtures and demonstrations are visibly marked simulated;
- stale data cannot be described as live;
- confidence falls as coverage or freshness degrades;
- an ethical/governance score may require review but must not downgrade a measurement-supported hazard alert.

### 6. Keep evidence acquisition bounded

If AMOC later adopts public-web observation or document ingestion, require:

- HTTP(S)-only allowlists for operational data sources;
- private-network and loopback denial;
- redirect, timeout, response-size, archive-expansion, and MIME limits;
- explicit source citations and retrieval timestamps;
- untrusted-content handling;
- no automatic conversion of retrieved text into truth or training data.

## Security and deployment guidance

The monograph distinguishes active controls from retained security modules. AMOC should do the same and avoid security-by-inventory.

Minimum AMOC controls:

- load tokens and keys only from environment or protected secret references;
- default local demos to `127.0.0.1` rather than public bind;
- authenticate administrative, ingestion, task-execution, and webhook-management routes;
- do not auto-install packages during application startup;
- use pinned dependencies and repeatable deployment manifests;
- validate webhook destinations and block private-network SSRF;
- keep blockchain/payment code outside environmental and decoder decision paths;
- verify checkpoints, manifests, and evidence logs before claiming continuity;
- surface terminal failures rather than silently replacing an authoritative component.

The monograph's Internet-on default and unauthenticated general API posture should **not** be copied into AMOC without a threat model and explicit deployment approval.

## Capability mapping

### Suitable for the current AMOC MVP

- capability-status taxonomy;
- source freshness and no-fabrication rules;
- normalized provenance envelope;
- append-only decision history and digest chaining;
- deterministic Tryte9 invariants;
- trajectory/admissibility logging;
- region-specific NOAA/USGS adapters;
- human-readable dashboard evidence.

### Future-state or separately qualified

- PostgreSQL-backed durable event service;
- specialist-agent routing and conservative confidence correction;
- document ingestion with bounded extraction;
- hardware sensory projections;
- cryptographic checkpoint and manifest verification;
- air-gapped deployment profile;
- FPGA/ASIC verification harnesses.

### Not established for AMOC by this document

- a production-ready native 8-layer GRU language system;
- 62 operational specialists;
- validated UEF/K7 moral geometry;
- exact-qutrit physical hardware;
- consciousness or ontological-randomness claims;
- production security certification;
- cryogenic, radiation, or energy-recovery performance;
- compatibility with the current AMOC runtime.

Those claims require the referenced archive, dependency lock, fresh tests, runtime traces, security review, and performance evidence before adoption.

## Minimum integration sequence

1. Add the shared evidence envelope and capability-status labels.
2. Refactor NOAA and USGS inputs behind source-adapter interfaces.
3. Add stale/partial/offline handling and fixture-labelled simulation.
4. Store normalized observations and decisions with SHA-256 digests.
5. Expose provenance in regional status endpoints and the dashboard.
6. Keep Tryte9 as an optional Sentinel Decode feature with exhaustive equivalence tests.
7. Evaluate specialist-agent, PostgreSQL, sensory, and training surfaces only as separate future tranches.

## Acceptance gate for SAEONYX-derived code

Before importing implementation from SAEONYX, require:

- the exact `SX_DEFINITIVE_FRESH.zip` archive;
- SHA-256 verification against the monograph value;
- manifest verification;
- license and ownership review;
- dependency and container review;
- fresh execution of all stated tests;
- API authentication and network-policy review;
- proof that claimed active paths are reachable from deployment entrypoints;
- an explicit mapping showing which code belongs to AMOC Sentinel versus BioForge.
