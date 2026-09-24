# AMOC / BioForge Token Budget Ledger

Baseline established: 2026-09-23 00:16 PDT
Model: `openai/gpt-5.6-sol`
Initial OpenClaw snapshot: 29k input / 3.4k output tokens; 5-hour allowance 87% remaining; weekly allowance 78% remaining.

> The displayed input/output counters are per-run/context workload indicators, not reliable cumulative billing counters: a later snapshot can increase, decrease, or be recomputed after context injection or compaction. Therefore this ledger tracks (a) implementation-budget estimates and (b) changes in the authoritative timed allowance percentages. It does not claim billing-exact per-project token accounting.

## Contract boundary

- BioForge: ALIGN, COMPLETE, ISSUE
- AMOC Sentinel: Sentinel Decode, Sentinel Pulse, Sentinel CryoPath
- Tryte9: optional bypassable Decode feature, not a SKU
- Sentinel Stack: bundle label, not an additional SKU unless separately priced

## Minimum viable integration budgets

| Item | Minimum scope | Estimated tokens |
|---|---|---:|
| Shared schemas/registry | Product IDs, contract boundary, common events, tests | 8k–12k |
| 1. ALIGN | Intake schema, evidence-pack API, validation, tests | 6k–10k |
| 2. COMPLETE | Append-only attestation log, O(1) lookup index, TDA adapter boundary, tests | 7k–12k |
| 3. ISSUE | Vault/escrow state machine, BT9-F job codec, release rules, tests | 9k–15k |
| 4. Sentinel Decode | Software reference decoder, latency harness, hardware interface contract, tests | 14k–24k |
| 5. Sentinel Pulse | Vendor-neutral pulse adapter, mock OPX/Qblox connectors, frame-update tests | 10k–18k |
| 6. Sentinel CryoPath | Energy/thermal accounting model and qualification harness only | 7k–12k |
| Tryte9 option | Bypassable combinational prefilter plus equivalence tests | 4k–7k |

Minimum staged total: **65k–110k tokens**.

The minimum does not include production FPGA/ASIC RTL, live OPX/Qblox hardware certification, cryogenic hardware qualification, security audit, or custom edge/GPU kernels.

## Running usage

| Date | Work item | Input delta | Output delta | Approx. total | Status |
|---|---|---:|---:|---:|---|
| 2026-09-23 00:16 | Baseline and planning | n/a | n/a | n/a | 5h 87% left; week 78% left |
| 2026-09-23 00:21 | Tryte9 visual-spec intake | n/a | n/a | n/a | 5h 84% left; week 78% left |
| 2026-09-23 00:32 | Tryte9 minimum demo implementation | n/a | n/a | n/a | 5h 62% left; week 74% left; tests/build/live API verified; local commit `73aa051` |
| 2026-09-24 12:33 | Great Salt Lake Sentinel integration assessment | n/a | n/a | n/a | 5h 95% left; week 66% left; architecture/security review only; no runtime code added |

## Token-control rules

1. Receive and inspect one framework at a time.
2. Define an acceptance test before editing code.
3. Reuse shared schemas and fixtures across products.
4. Avoid subagents unless parallel research materially reduces total work.
5. Implement the smallest vertical slice, test it, then stop for review.
6. Record a fresh OpenClaw usage snapshot before and after each numbered item.
7. Do not count Tryte9 or Sentinel Stack as separate SKUs.
8. Keep BioForge and AMOC Sentinel APIs and contracts separate.
