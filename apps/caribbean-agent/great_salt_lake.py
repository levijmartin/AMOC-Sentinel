"""Great Salt Lake regional profile for AMOC Sentinel.

This module is an independent AMOC implementation built against the public
USGS NWIS JSON contract. It does not contain SAEONYX GSL package source.
Missing measurements stay missing; gage height is never treated as absolute
lake elevation.
"""

from __future__ import annotations

import hashlib
import json
import os
import time
import urllib.parse
import urllib.request
from dataclasses import asdict, dataclass
from datetime import datetime, timezone
from typing import Any, Dict, Iterable, List, Optional

USGS_NWIS_IV_URL = "https://waterservices.usgs.gov/nwis/iv/"

# Public USGS monitoring locations relevant to the Great Salt Lake profile.
GSL_SITES: Dict[str, Dict[str, str]] = {
    "10010000": {"name": "Great Salt Lake at Saltair Boat Harbor", "role": "lake"},
    "10126000": {"name": "Bear River near Corinne", "role": "inflow"},
    "10141000": {"name": "Weber River near Plain City", "role": "inflow"},
    "10170500": {"name": "Surplus Canal at Salt Lake City", "role": "inflow"},
}

PARAMETERS = {
    "62614": ("water_surface_elevation_ft", "ft"),
    "00065": ("gage_height_ft", "ft"),
    "00060": ("discharge_cfs", "ft3/s"),
    "00010": ("water_temperature_c", "deg C"),
}


@dataclass(frozen=True)
class GSLObservation:
    site_id: str
    site_name: str
    role: str
    parameter_code: str
    field: str
    value: Optional[float]
    unit: str
    observed_at: Optional[str]
    latitude: Optional[float]
    longitude: Optional[float]
    previous_value: Optional[float] = None
    delta: Optional[float] = None
    source: str = "USGS NWIS instantaneous values"


def _float_or_none(value: Any) -> Optional[float]:
    try:
        if value is None or str(value).strip() in {"", "-999999", "Ice", "Eqp", "Dis"}:
            return None
        return float(value)
    except (TypeError, ValueError):
        return None


def _parse_timestamp(value: Optional[str]) -> Optional[datetime]:
    if not value:
        return None
    try:
        return datetime.fromisoformat(value.replace("Z", "+00:00")).astimezone(timezone.utc)
    except ValueError:
        return None


def parse_usgs_payload(payload: Dict[str, Any]) -> List[GSLObservation]:
    """Parse NWIS time-series data without manufacturing missing values."""
    observations: List[GSLObservation] = []
    series = payload.get("value", {}).get("timeSeries", [])
    for item in series:
        source_info = item.get("sourceInfo", {})
        site_codes = source_info.get("siteCode", [])
        site_id = str(site_codes[0].get("value", "")) if site_codes else ""
        if site_id not in GSL_SITES:
            continue

        variable_codes = item.get("variable", {}).get("variableCode", [])
        parameter_code = str(variable_codes[0].get("value", "")) if variable_codes else ""
        if parameter_code not in PARAMETERS:
            continue
        field, default_unit = PARAMETERS[parameter_code]

        geog = source_info.get("geoLocation", {}).get("geogLocation", {})
        latitude = _float_or_none(geog.get("latitude"))
        longitude = _float_or_none(geog.get("longitude"))
        unit = item.get("variable", {}).get("unit", {}).get("unitCode") or default_unit

        values: List[Dict[str, Any]] = []
        for block in item.get("values", []):
            values.extend(block.get("value", []))
        if not values:
            value = previous = None
            observed_at = None
        else:
            latest = values[-1]
            value = _float_or_none(latest.get("value"))
            observed_at = latest.get("dateTime")
            previous = _float_or_none(values[-2].get("value")) if len(values) > 1 else None
        delta = value - previous if value is not None and previous is not None else None

        meta = GSL_SITES[site_id]
        observations.append(
            GSLObservation(
                site_id=site_id,
                site_name=source_info.get("siteName") or meta["name"],
                role=meta["role"],
                parameter_code=parameter_code,
                field=field,
                value=value,
                unit=str(unit),
                observed_at=observed_at,
                latitude=latitude,
                longitude=longitude,
                previous_value=previous,
                delta=delta,
            )
        )
    return observations


def summarize_observations(
    observations: Iterable[GSLObservation],
    *,
    fetched_at: Optional[float] = None,
    baseline_elevation_ft: Optional[float] = None,
) -> Dict[str, Any]:
    fetched_at = fetched_at or time.time()
    rows = list(observations)
    now = datetime.fromtimestamp(fetched_at, timezone.utc)

    ages: List[float] = []
    for row in rows:
        observed = _parse_timestamp(row.observed_at)
        if observed:
            ages.append(max(0.0, (now - observed).total_seconds()))
    newest_age = min(ages) if ages else None
    freshness = "offline"
    if rows:
        freshness = "live" if newest_age is not None and newest_age <= 2 * 60 * 60 else "stale"
        if any(row.value is None for row in rows):
            freshness = "partial" if freshness == "live" else freshness

    lake_elevation = next(
        (
            row
            for row in rows
            if row.site_id == "10010000"
            and row.parameter_code == "62614"
            and row.value is not None
        ),
        None,
    )
    # Deliberately exclude parameter 00065: gage height is not absolute elevation.
    inflows = [
        row
        for row in rows
        if row.role == "inflow" and row.parameter_code == "00060" and row.value is not None
    ]
    total_inflow = sum(row.value for row in inflows) if inflows else None
    inflow_delta = (
        sum(row.delta for row in inflows if row.delta is not None)
        if inflows and all(row.delta is not None for row in inflows)
        else None
    )

    context: Dict[str, Any] = {
        "lake_surface_elevation_ft": lake_elevation.value if lake_elevation else None,
        "lake_surface_elevation_delta_ft": lake_elevation.delta if lake_elevation else None,
        "combined_inflow_cfs": total_inflow,
        "combined_inflow_delta_cfs": inflow_delta,
        "inflow_sites_reporting": len(inflows),
        "freshness": freshness,
        "newest_observation_age_seconds": newest_age,
    }
    if baseline_elevation_ft is not None and lake_elevation:
        context["analysis_reference_elevation_ft"] = baseline_elevation_ft
        context["elevation_vs_reference_ft"] = lake_elevation.value - baseline_elevation_ft

    normalized = [asdict(row) for row in sorted(rows, key=lambda x: (x.site_id, x.parameter_code))]
    digest_payload = {"observations": normalized, "context": context}
    digest = hashlib.sha256(
        json.dumps(digest_payload, sort_keys=True, separators=(",", ":")).encode("utf-8")
    ).hexdigest()

    notices: List[str] = []
    if freshness == "offline":
        notices.append("No supported USGS observations were available; no physical condition was inferred.")
    elif freshness == "stale":
        notices.append("The newest supported USGS observation is stale; review upstream availability before acting.")
    elif freshness == "partial":
        notices.append("The snapshot is partial; missing measurements remain missing.")
    if baseline_elevation_ft is not None:
        notices.append("The configured elevation reference is analytical context, not a sensor substitute or hazard threshold.")

    return {
        "profile": "great-salt-lake",
        "capability_status": "prototype",
        "source": USGS_NWIS_IV_URL,
        "fetched_at": fetched_at,
        "sites_requested": list(GSL_SITES),
        "observation_count": len(rows),
        "context": context,
        "observations": normalized,
        "notices": notices,
        "provenance": {
            "transform_version": "amoc-gsl-profile-v1",
            "output_digest": f"sha256:{digest}",
            "measurement_policy": "missing values remain missing; 00065 gage height is not 62614 elevation",
        },
    }


class GreatSaltLakeProfile:
    def __init__(self, timeout: float = 15.0):
        self.timeout = timeout
        raw_baseline = os.getenv("GSL_ANALYSIS_REFERENCE_ELEVATION_FT", "").strip()
        self.baseline_elevation_ft = _float_or_none(raw_baseline)

    def fetch_payload(self) -> Dict[str, Any]:
        query = urllib.parse.urlencode(
            {
                "format": "json",
                "sites": ",".join(GSL_SITES),
                "parameterCd": ",".join(PARAMETERS),
                "siteStatus": "all",
                "period": "P1D",
            }
        )
        request = urllib.request.Request(
            f"{USGS_NWIS_IV_URL}?{query}",
            headers={"User-Agent": "AMOC-Sentinel/Great-Salt-Lake-Profile-v1"},
        )
        with urllib.request.urlopen(request, timeout=self.timeout) as response:
            if response.status != 200:
                raise RuntimeError(f"USGS NWIS returned HTTP {response.status}")
            return json.loads(response.read().decode("utf-8"))

    def snapshot(self) -> Dict[str, Any]:
        fetched_at = time.time()
        payload = self.fetch_payload()
        observations = parse_usgs_payload(payload)
        return summarize_observations(
            observations,
            fetched_at=fetched_at,
            baseline_elevation_ft=self.baseline_elevation_ft,
        )

    def metadata(self) -> Dict[str, Any]:
        return {
            "profile": "great-salt-lake",
            "capability_status": "prototype",
            "source": USGS_NWIS_IV_URL,
            "sites": GSL_SITES,
            "parameters": {
                code: {"field": field, "unit": unit} for code, (field, unit) in PARAMETERS.items()
            },
            "analysis_reference_elevation_ft": self.baseline_elevation_ft,
        }
