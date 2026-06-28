# AMOC Atlas Reference

## Source
GitHub: `AMOCcommunity/amocatlas`

Repository URL:
https://github.com/AMOCcommunity/amocatlas

## Why it matters
AMOC Atlas appears relevant as a scientific and technical reference point for Atlantic Meridional Overturning Circulation (AMOC) analysis, organization, and possibly visualization workflows. For StormShield Caribbean, this makes it a useful upstream reference for the climate/ocean intelligence layer.

## Best role in this project
This repository should be treated as a **reference source**, not a core MVP dependency.

Use it for:
- understanding AMOC-related data and analysis patterns
- informing climate/ocean signal interpretation
- improving scientific framing around Atlantic circulation context
- borrowing ideas for analysis structure or visual presentation

Do not use it for:
- blindly copying architecture into the buildathon MVP
- replacing the operational Caribbean signal stack
- making the project overly academic at the expense of usability

## How it fits StormShield Caribbean
A clean layering for the project is:

- **AMOC Atlas / AMOC research sources** -> upstream scientific context
- **RAPID / Copernicus / Argo** -> ocean-state and climate signal layer
- **CariCOOS / NOAA NDBC** -> regional and local operational observations
- **StormShield Caribbean agents** -> localized risk reasoning and action planning

## Practical value
The project can use AMOC Atlas as a reference for:
- AMOC framing language
- possible analysis workflows
- interpretation of circulation-related context
- chart or atlas-style presentation ideas

## Recommendation
Preserve this as a documented reference and revisit it when:
- refining the ocean-signal layer
- improving scientific credibility in the application
- designing visualizations or data products

## Bottom line
AMOC Atlas is a useful research and reference input for the project’s climate/ocean intelligence layer, but it should support the MVP rather than define it.
