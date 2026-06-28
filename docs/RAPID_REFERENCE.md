# RAPID Reference

## Source
RAPID Array data portal:
https://rapid.ac.uk/

Referenced download bundle:
https://rapid.ac.uk/download-confirm/moc_netcdf+moc_matlab+moc_ascii+moc_vertical_netcdf+gridts_netcdf+uncertain_matlab+meridional_transports+2d_gridded

## Why it matters
RAPID is one of the strongest scientific sources for AMOC and Atlantic overturning circulation context. For StormShield Caribbean, it is a high-value upstream source for broad ocean-state interpretation and long-range circulation context.

## What it appears to include
The referenced bundle appears to include:
- AMOC NetCDF outputs
- Matlab files
- ASCII data
- vertical transport data
- grid time series
- uncertainty-related files
- meridional transport data
- 2D gridded outputs

## Best role in this project
RAPID should be treated as a **core climate and ocean context source**.

Use it for:
- AMOC trend context
- circulation change interpretation
- large-scale Atlantic ocean-state framing
- scientific credibility in the climate/ocean intelligence layer

Do not rely on it for:
- hyperlocal daily weather
- short-horizon consumer forecasting
- immediate storm alerting by itself

## How it fits the project stack
Recommended layering:
- **RAPID** -> Atlantic circulation backbone
- **Copernicus Marine** -> operational ocean-state layer
- **CariCOOS / NOAA NDBC** -> regional and local observational layer
- **StormShield Caribbean agents** -> local risk reasoning and resilience planning

## Practical MVP stance
For the MVP:
- preserve RAPID as a core reference and signal source
- use summarized or selected timeseries before attempting full heavy processing
- combine RAPID context with more operational datasets rather than using it alone

## Recommendation
RAPID should remain one of the named core sources in the repo and architecture, especially when explaining the scientific basis of the Atlantic-context layer.
