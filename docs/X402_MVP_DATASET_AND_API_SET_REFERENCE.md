# x402 MVP Dataset and API Set Reference

This document adds a structured reference for the uploaded paper **"MVP Dataset and API Set on the Measurement Coordinate System"**.

## Added artifact

- `docs/x402-mvp-dataset-and-api-set.json`

## What it contains

The JSON extraction currently captures the portions visible in the provided text, including:

- the introductory MVP framing
- the one-mechanic walkthrough
- packetization, progressive truncation, alarm semantics, and cross-check rules
- base seating and inversion wave notes
- the electrical quantities walkthrough:
  - voltage
  - current
  - resistance
  - power
- detailed voltage packet examples
- Source 1 (Copernicus Marine Toolbox) mapping and examples
- Source 9 (OpenTelemetry Collector) and Source 10 (x402) rail roles

## Important note

This is a **structured extraction**, not yet a full canonical transcription of the entire paper.

The uploaded text in this session appears to include:
- the full introductory and mechanics sections
- the electrical quantities walkthrough
- Source 1

It does **not** appear to include the full remaining source sections in the visible excerpt. Because of that, the JSON file intentionally marks itself as a partial extraction where appropriate.

## Suggested next steps

1. Continue extracting the remaining source sections if additional pages/text are provided.
2. Add a master JSON Schema once the full MVP paper is captured.
3. Cross-link this file from `docs/X402_REFERENCE.md` or `README.md` if it becomes the canonical structured reference.
