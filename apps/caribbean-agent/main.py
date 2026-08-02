#!/usr/bin/env python3
"""
SAEONYX Caribbean Ocean Agent v3.1
Regional real-time safety and coordination infrastructure.
Substrate-driven. Real NOAA NDBC only. Zero simulation.
The geometric state is the risk signal. The agent is the coordination layer.
"""

from __future__ import annotations

import json
import logging
import os
import threading
import time
import urllib.request
import uuid
from collections import deque
from contextlib import asynccontextmanager
from dataclasses import asdict, dataclass, field
from enum import Enum
from typing import Any, Dict, List, Optional, Tuple

import aiosqlite
import numpy as np
import uvicorn
from fastapi import FastAPI, Header, HTTPException
from fastapi.responses import HTMLResponse
from pydantic import BaseModel

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(name)s: %(message)s")
log = logging.getLogger("saeonyx.caribbean")


class Config:
    HOST = os.getenv("HOST", "0.0.0.0")
    PORT = int(os.getenv("PORT", "7735"))
    DATABASE_URL = os.getenv("DATABASE_URL", "caribbean_agent_v31.db")
    ADMIN_TOKEN = os.getenv("ADMIN_TOKEN", "")
    POLL_INTERVAL = int(os.getenv("POLL_INTERVAL", "50"))
    AGENT_ID = "saeonyx-caribbean-ocean-agent-v3.1"
    AGENT_NAME = "SAEONYX Caribbean Ocean Coordination Agent"
    VERSION = "3.1.0"

    BUOY_STATIONS = {
        "42059": {"name": "Eastern Caribbean Sea", "lat": 15.255, "lng": -67.621, "region": "eastern"},
        "42058": {"name": "Central Caribbean", "lat": 14.92, "lng": -75.05, "region": "central"},
        "42057": {"name": "Western Caribbean", "lat": 16.97, "lng": -81.57, "region": "western"},
        "42060": {"name": "Caribbean Valley", "lat": 16.42, "lng": -63.20, "region": "eastern"},
        "42065": {"name": "Near Central Caribbean", "lat": 14.926, "lng": -75.046, "region": "central"},
        "41053": {"name": "San Juan, PR", "lat": 18.474, "lng": -66.099, "region": "northeast"},
        "41052": {"name": "South of St. John, VI", "lat": 18.249, "lng": -64.763, "region": "northeast"},
        "41058": {"name": "North of St. Thomas, VI", "lat": 18.476, "lng": -65.157, "region": "northeast"},
        "41051": {"name": "CarICOOS 41051", "lat": 18.26, "lng": -65.00, "region": "northeast"},
        "41056": {"name": "CarICOOS 41056", "lat": 18.15, "lng": -65.00, "region": "northeast"},
        "42001": {"name": "Mid Gulf of Mexico", "lat": 25.90, "lng": -89.67, "region": "approach"},
        "42003": {"name": "East Gulf", "lat": 26.07, "lng": -85.93, "region": "approach"},
        "42036": {"name": "West of Tampa", "lat": 28.50, "lng": -84.52, "region": "approach"},
        "42039": {"name": "Pensacola", "lat": 28.79, "lng": -86.01, "region": "approach"},
        "41009": {"name": "East of Cape Canaveral", "lat": 28.52, "lng": -80.17, "region": "approach"},
        "41010": {"name": "East of Cape Canaveral 2", "lat": 28.95, "lng": -78.47, "region": "approach"},
        "41008": {"name": "Southeast of Savannah", "lat": 31.40, "lng": -80.87, "region": "approach"},
        "41004": {"name": "Southeast of Charleston", "lat": 32.50, "lng": -79.10, "region": "approach"},
        "41025": {"name": "Diamond Shoals", "lat": 35.01, "lng": -75.40, "region": "approach"},
        "41002": {"name": "South Hatteras", "lat": 31.76, "lng": -74.84, "region": "approach"},
        "41001": {"name": "East of Cape Hatteras", "lat": 34.70, "lng": -72.73, "region": "approach"},
    }


VIRTUES = ["honesty", "compassion", "fairness", "courage", "wisdom", "temperance", "justice", "respect"]


@dataclass
class SoulVector:
    values: np.ndarray
    divergence: float
    scalar: float

    @property
    def as_dict(self) -> Dict[str, float]:
        return {v: float(self.values[i]) for i, v in enumerate(VIRTUES)}


def compute_divergence(values: np.ndarray) -> float:
    n = len(values)
    return float(sum(abs(values[(i + 1) % n] - values[(i - 1) % n]) / 2.0 for i in range(n)) / n)


def project_divergence_free(values: np.ndarray) -> np.ndarray:
    mean_val = float(np.mean(values))
    return np.clip((1.0 - 0.3) * values + 0.3 * mean_val, 0.0, 1.0)


def compute_scalar(values: np.ndarray) -> float:
    return float(np.exp(np.mean(np.log(np.maximum(values, 1e-6)))))


class AxiomOfOneSubstrate:
    def __init__(self, dim: int = 8):
        self.dim = dim
        self.gaussian_field = np.random.normal(0, 1, dim)
        self._normalize()
        self.q = self.gaussian_field.copy()
        self.p = np.zeros(dim)
        self.energy_history: deque[float] = deque(maxlen=40)
        self.phi_history: deque[float] = deque(maxlen=40)

    def _normalize(self) -> None:
        n = np.linalg.norm(self.gaussian_field)
        self.gaussian_field = self.gaussian_field / n if n > 1e-12 else np.ones(self.dim) / np.sqrt(self.dim)

    def step_hamiltonian_flow(self, dt: float = 0.011) -> float:
        self.p -= self.q * dt
        self.q += self.p * dt
        self.gaussian_field = self.q.copy()
        self._normalize()
        self.q = self.gaussian_field.copy()
        energy = float(0.5 * np.sum(self.p ** 2) + 0.5 * np.sum(self.q ** 2))
        self.energy_history.append(energy)
        return energy

    def lyapunov_basin_projection(self) -> float:
        target = np.ones(self.dim) / np.sqrt(self.dim)
        diff = self.gaussian_field - target
        lyap = 0.5 * np.sum(diff ** 2)
        self.gaussian_field -= 0.038 * diff
        self._normalize()
        self.q = self.gaussian_field.copy()
        return float(lyap)

    def to_soul_vector(self) -> SoulVector:
        raw = np.clip(np.abs(self.gaussian_field), 0.01, 1.0)
        proj = project_divergence_free(raw)
        return SoulVector(proj, compute_divergence(proj), compute_scalar(proj))

    def inject_ocean_state(self, feature_vector: np.ndarray, strength: float = 0.38) -> None:
        if feature_vector is None or len(feature_vector) == 0:
            return
        vec = np.asarray(feature_vector, dtype=np.float64)
        if len(vec) < self.dim:
            vec = np.pad(vec, (0, self.dim - len(vec)))
        vec = vec[: self.dim]
        n = np.linalg.norm(vec)
        if n > 1e-9:
            self.q = (1.0 - strength) * self.q + strength * (vec / n)
            self.gaussian_field = self.q.copy()
            self._normalize()


class K7Evaluator:
    def evaluate(self, soul: SoulVector, phi: float) -> Dict[str, Any]:
        channels = {
            "consistency": 0.05,
            "continuity": max(0.0, 1.0 - phi) * 0.4,
            "causality": min(soul.divergence * 2.0, 1.0),
            "boundedness": max(0.0, 1.0 - phi) * 0.3,
            "non_divergence": float(np.var(soul.values)) * 2.0,
            "moral_alignment": max(0.0, 0.85 - soul.scalar) * 0.5,
            "integrability": min(abs(soul.values[0] - soul.values[-1]) * 0.5, 1.0),
        }
        channels = {k: float(np.clip(v, 0, 1)) for k, v in channels.items()}
        k7_sum = sum(channels.values())
        stable = k7_sum < 0.50 and all(v < 0.20 for v in channels.values())
        return {"channels": channels, "stable": stable, "k7_sum": k7_sum}


class DecisionLevel(str, Enum):
    NORMAL = "NORMAL"
    WATCH = "WATCH"
    ALERT = "ALERT"
    ESCALATE = "ESCALATE"


@dataclass
class AgentDecision:
    decision_id: str
    timestamp: float
    level: DecisionLevel
    reason: str
    affected_stations: List[str]
    coordination_signal: Dict[str, Any]
    soul_scalar: float
    k7_stable: bool
    k7_sum: float
    phi: float
    substrate_energy: float
    energy_trend: float
    live_station_count: int
    sensor_snapshot: Dict[str, Any]
    actions_fired: List[str] = field(default_factory=list)


@dataclass
class Subscriber:
    subscriber_id: str
    name: str
    webhook_url: Optional[str]
    min_level: str
    registered_at: float
    role: str = "general"


class DatabaseManager:
    def __init__(self, db_url: str):
        self.db_url = db_url

    async def initialize(self) -> None:
        async with aiosqlite.connect(self.db_url) as db:
            await db.execute("PRAGMA journal_mode=WAL")
            await db.execute(
                """
                CREATE TABLE IF NOT EXISTS decisions (
                    id TEXT PRIMARY KEY,
                    timestamp REAL NOT NULL,
                    level TEXT NOT NULL,
                    reason TEXT NOT NULL,
                    phi REAL NOT NULL,
                    energy REAL NOT NULL,
                    energy_trend REAL NOT NULL,
                    live_count INTEGER NOT NULL,
                    payload_json TEXT NOT NULL
                )
                """
            )
            await db.execute(
                """
                CREATE TABLE IF NOT EXISTS actions (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    timestamp REAL NOT NULL,
                    decision_id TEXT NOT NULL,
                    subscriber TEXT NOT NULL,
                    role TEXT NOT NULL,
                    status TEXT NOT NULL,
                    level TEXT NOT NULL,
                    payload_json TEXT NOT NULL
                )
                """
            )
            await db.execute(
                """
                CREATE TABLE IF NOT EXISTS subscribers (
                    subscriber_id TEXT PRIMARY KEY,
                    name TEXT NOT NULL,
                    webhook_url TEXT,
                    min_level TEXT NOT NULL,
                    registered_at REAL NOT NULL,
                    role TEXT NOT NULL
                )
                """
            )
            await db.commit()

    async def log_decision(self, decision: AgentDecision) -> None:
        payload = json.dumps(asdict(decision))
        async with aiosqlite.connect(self.db_url) as db:
            await db.execute(
                """
                INSERT OR REPLACE INTO decisions
                (id, timestamp, level, reason, phi, energy, energy_trend, live_count, payload_json)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
                """,
                (
                    decision.decision_id,
                    decision.timestamp,
                    decision.level.value,
                    decision.reason,
                    decision.phi,
                    decision.substrate_energy,
                    decision.energy_trend,
                    decision.live_station_count,
                    payload,
                ),
            )
            await db.commit()

    async def log_action(self, action: Dict[str, Any]) -> None:
        async with aiosqlite.connect(self.db_url) as db:
            await db.execute(
                """
                INSERT INTO actions (timestamp, decision_id, subscriber, role, status, level, payload_json)
                VALUES (?, ?, ?, ?, ?, ?, ?)
                """,
                (
                    action["ts"],
                    action["decision_id"],
                    action["subscriber"],
                    action["role"],
                    action["status"],
                    action["level"],
                    json.dumps(action),
                ),
            )
            await db.commit()

    async def upsert_subscriber(self, subscriber: Subscriber) -> None:
        async with aiosqlite.connect(self.db_url) as db:
            await db.execute(
                """
                INSERT OR REPLACE INTO subscribers
                (subscriber_id, name, webhook_url, min_level, registered_at, role)
                VALUES (?, ?, ?, ?, ?, ?)
                """,
                (
                    subscriber.subscriber_id,
                    subscriber.name,
                    subscriber.webhook_url,
                    subscriber.min_level,
                    subscriber.registered_at,
                    subscriber.role,
                ),
            )
            await db.commit()


def _safe_float(val: Any) -> Optional[float]:
    if val is None:
        return None
    s = str(val).strip().upper()
    if s in ("MM", "N/A", "", "NAN", "NONE", "--", "NULL"):
        return None
    try:
        return float(s)
    except (ValueError, TypeError):
        return None


def fetch_ndbc_raw(station_id: str) -> Optional[Dict[str, str]]:
    url = f"https://www.ndbc.noaa.gov/data/realtime2/{station_id}.txt"
    try:
        req = urllib.request.Request(url, headers={"User-Agent": f"SAEONYX-Caribbean-Agent/{Config.VERSION}"})
        with urllib.request.urlopen(req, timeout=8) as resp:
            text = resp.read().decode("utf-8", errors="replace").strip()
        lines = [ln for ln in text.splitlines() if ln.strip()]
        if len(lines) < 3:
            return None
        headers = lines[0].split()
        data_line = next((ln for ln in lines[2:] if not ln.startswith("#")), None)
        if not data_line:
            return None
        values = data_line.split()
        if len(values) < len(headers):
            values = values + ["MM"] * (len(headers) - len(values))
        return dict(zip(headers, values))
    except Exception as e:
        log.warning("NDBC fetch failed for %s: %s", station_id, e)
        return None


def get_live_sensors() -> List[Dict[str, Any]]:
    sensors: List[Dict[str, Any]] = []
    for sid, meta in Config.BUOY_STATIONS.items():
        raw = fetch_ndbc_raw(sid)
        if raw is None:
            sensors.append(
                {
                    "station_id": sid,
                    "name": meta["name"],
                    "lat": meta["lat"],
                    "lng": meta["lng"],
                    "region": meta.get("region", "unknown"),
                    "sst_c": None,
                    "wave_height_m": None,
                    "wind_speed_ms": None,
                    "wind_dir": None,
                    "pressure_hpa": None,
                    "air_temp_c": None,
                    "source": "OFFLINE",
                    "status": "offline",
                    "real_field_count": 0,
                }
            )
            continue

        sst = _safe_float(raw.get("WTMP"))
        wave = _safe_float(raw.get("WVHT"))
        wind = _safe_float(raw.get("WSPD"))
        wdir = _safe_float(raw.get("WDIR"))
        pres = _safe_float(raw.get("PRES"))
        atmp = _safe_float(raw.get("ATMP"))
        real_fields = sum(1 for v in (sst, wave, wind, wdir, pres, atmp) if v is not None)

        sensors.append(
            {
                "station_id": sid,
                "name": meta["name"],
                "lat": meta["lat"],
                "lng": meta["lng"],
                "region": meta.get("region", "unknown"),
                "sst_c": sst,
                "wave_height_m": wave,
                "wind_speed_ms": wind,
                "wind_dir": wdir,
                "pressure_hpa": pres,
                "air_temp_c": atmp,
                "source": "NOAA NDBC live",
                "status": "live" if real_fields >= 2 else "partial",
                "real_field_count": real_fields,
                "raw_time": f"{raw.get('YY', '??')}-{raw.get('MM', '??')}-{raw.get('DD', '??')} {raw.get('hh', '??')}:{raw.get('mm', '??')}",
            }
        )
    return sensors


def build_ocean_feature_vector(sensors: List[Dict[str, Any]]) -> Optional[np.ndarray]:
    waves: List[float] = []
    winds: List[float] = []
    pressures: List[float] = []
    ssts: List[float] = []
    lats: List[float] = []
    lngs: List[float] = []

    for s in sensors:
        if s.get("status") not in ("live", "partial"):
            continue
        if isinstance(s.get("wave_height_m"), (int, float)):
            waves.append(float(s["wave_height_m"]))
        if isinstance(s.get("wind_speed_ms"), (int, float)):
            winds.append(float(s["wind_speed_ms"]))
        if isinstance(s.get("pressure_hpa"), (int, float)):
            pressures.append(float(s["pressure_hpa"]))
        if isinstance(s.get("sst_c"), (int, float)):
            ssts.append(float(s["sst_c"]))
        lats.append(float(s["lat"]))
        lngs.append(float(s["lng"]))

    if not (waves or winds or pressures or ssts):
        return None

    mean_lat = float(np.mean(lats)) if lats else 16.0
    mean_lng = float(np.mean(lngs)) if lngs else -70.0

    return np.array(
        [
            (max(waves) if waves else 0.0) / 5.0,
            (float(np.mean(waves)) if waves else 0.0) / 3.0,
            (max(winds) if winds else 0.0) / 25.0,
            (float(np.mean(winds)) if winds else 0.0) / 15.0,
            (1015.0 - (min(pressures) if pressures else 1013.0)) / 25.0,
            (max(ssts) if ssts else 28.0) / 32.0,
            (mean_lat + 90.0) / 180.0,
            (mean_lng + 180.0) / 360.0,
        ],
        dtype=np.float64,
    )


class CaribbeanOceanAgent:
    def __init__(self):
        self.substrate = AxiomOfOneSubstrate()
        self.k7 = K7Evaluator()
        self._soul = self.substrate.to_soul_vector()
        self._phi = self._soul.scalar
        self._cycles = 0
        self.latest_decision: Optional[AgentDecision] = None
        self.decision_history: List[AgentDecision] = []
        self.action_log: List[Dict[str, Any]] = []
        self.subscribers: Dict[str, Subscriber] = {}
        self.running = False
        self._lock = threading.Lock()
        self._prev_energy = 0.0

    def _derive_level_from_substrate(
        self, energy: float, energy_trend: float, k7: Dict[str, Any], live_count: int
    ) -> Tuple[DecisionLevel, str]:
        if live_count == 0:
            return DecisionLevel.WATCH, "No live stations available - substrate running open-loop"

        risk = 0.0
        risk += min(energy * 1.85, 1.45)
        risk += max(0.0, energy_trend * 4.2)
        risk += (1.0 - self._phi) * 0.95
        risk += 0.0 if k7["stable"] else 0.75
        risk += max(0.0, k7["k7_sum"] - 0.38) * 0.85

        if risk >= 2.65:
            level = DecisionLevel.ESCALATE
            reason = (
                f"Substrate risk critical | energy={energy:.3f} trend={energy_trend:+.3f} "
                f"phi={self._phi:.3f} K7={'stable' if k7['stable'] else 'unstable'} | {live_count} stations live"
            )
        elif risk >= 1.75:
            level = DecisionLevel.ALERT
            reason = (
                f"Substrate risk elevated | energy={energy:.3f} trend={energy_trend:+.3f} "
                f"phi={self._phi:.3f} | {live_count} stations live"
            )
        elif risk >= 0.95:
            level = DecisionLevel.WATCH
            reason = f"Substrate indicating developing conditions | energy={energy:.3f} phi={self._phi:.3f} | {live_count} stations live"
        else:
            level = DecisionLevel.NORMAL
            reason = f"Substrate coherent - Caribbean state nominal | energy={energy:.3f} phi={self._phi:.3f} | {live_count} stations live"
        return level, reason

    def _build_coordination_signal(
        self, level: DecisionLevel, sensors: List[Dict[str, Any]], energy: float, phi: float
    ) -> Dict[str, Any]:
        live = [s for s in sensors if s.get("status") in ("live", "partial")]
        max_wave = max((s["wave_height_m"] for s in live if isinstance(s.get("wave_height_m"), (int, float))), default=None)
        max_wind = max((s["wind_speed_ms"] for s in live if isinstance(s.get("wind_speed_ms"), (int, float))), default=None)
        min_pres = min((s["pressure_hpa"] for s in live if isinstance(s.get("pressure_hpa"), (int, float))), default=None)
        max_sst = max((s["sst_c"] for s in live if isinstance(s.get("sst_c"), (int, float))), default=None)
        focus = [
            s["station_id"]
            for s in live
            if (isinstance(s.get("wave_height_m"), (int, float)) and s["wave_height_m"] >= 1.7)
            or (isinstance(s.get("wind_speed_ms"), (int, float)) and s["wind_speed_ms"] >= 10.5)
        ]
        urgency = {"NORMAL": 0.12, "WATCH": 0.38, "ALERT": 0.72, "ESCALATE": 0.94}[level.value]
        return {
            "agent_id": Config.AGENT_ID,
            "level": level.value,
            "urgency": urgency,
            "substrate_energy": round(energy, 4),
            "phi": round(phi, 4),
            "live_stations": len(live),
            "observed": {
                "max_wave_m": max_wave,
                "max_wind_ms": max_wind,
                "min_pressure_hpa": min_pres,
                "max_sst_c": max_sst,
            },
            "recommended_focus_stations": focus,
            "messages": {
                "ports": f"SAEONYX coordination signal: {level.value}. Substrate energy {energy:.3f}. Review berthing, pilotage, and small-craft status.",
                "shipping": f"Routing advisory: {level.value} conditions across monitored Caribbean sectors. Energy trend active.",
                "tourism": f"Marine operators: current agent level {level.value}. Review excursion and small-craft operations.",
                "disaster": (
                    f"CDEMA / emergency coordination: agent level {level.value}. Substrate indicates elevated regional attention required."
                    if level in (DecisionLevel.ALERT, DecisionLevel.ESCALATE)
                    else None
                ),
            },
            "timestamp": time.time(),
        }

    def run_cycle(self) -> AgentDecision:
        sensors = get_live_sensors()
        live_count = sum(1 for s in sensors if s.get("status") in ("live", "partial"))

        features = build_ocean_feature_vector(sensors)
        if features is not None:
            self.substrate.inject_ocean_state(features)

        energy = self.substrate.step_hamiltonian_flow()
        self.substrate.lyapunov_basin_projection()
        self._soul = self.substrate.to_soul_vector()
        k7 = self.k7.evaluate(self._soul, self._phi)

        self._phi = min(
            1.0,
            max(
                0.0,
                0.51 * self._soul.scalar + 0.29 * (1.0 - k7["k7_sum"]) + 0.20 * (1.0 - self._soul.divergence * 3.4),
            ),
        )
        self.substrate.phi_history.append(self._phi)
        self._cycles += 1

        energy_trend = energy - self._prev_energy if self._prev_energy else 0.0
        self._prev_energy = energy
        level, reason = self._derive_level_from_substrate(energy, energy_trend, k7, live_count)

        if not k7["stable"] and level in (DecisionLevel.ESCALATE, DecisionLevel.ALERT):
            level = DecisionLevel.WATCH
            reason += " | Moral governor restrained escalation (K7 unstable)"

        coord = self._build_coordination_signal(level, sensors, energy, self._phi)
        affected = coord.get("recommended_focus_stations", [])

        decision = AgentDecision(
            decision_id=f"dec-{uuid.uuid4().hex[:10]}",
            timestamp=time.time(),
            level=level,
            reason=reason,
            affected_stations=affected,
            coordination_signal=coord,
            soul_scalar=round(self._soul.scalar, 4),
            k7_stable=bool(k7["stable"]),
            k7_sum=round(k7["k7_sum"], 4),
            phi=round(self._phi, 4),
            substrate_energy=round(energy, 4),
            energy_trend=round(energy_trend, 4),
            live_station_count=live_count,
            sensor_snapshot={
                s["station_id"]: {
                    "sst": s.get("sst_c"),
                    "wave": s.get("wave_height_m"),
                    "wind": s.get("wind_speed_ms"),
                    "pressure": s.get("pressure_hpa"),
                    "status": s.get("status"),
                    "region": s.get("region"),
                }
                for s in sensors
            },
        )

        fired = self._broadcast(decision)
        decision.actions_fired = fired

        with self._lock:
            self.latest_decision = decision
            self.decision_history.append(decision)
            if len(self.decision_history) > 100:
                self.decision_history = self.decision_history[-100:]

        log.info(
            "AGENT [%s] live=%d energy=%.3f trend=%+.3f phi=%.3f K7=%s fired=%d",
            level.value,
            live_count,
            energy,
            energy_trend,
            self._phi,
            k7["stable"],
            len(fired),
        )
        return decision

    def _broadcast(self, decision: AgentDecision) -> List[str]:
        fired: List[str] = []
        level_rank = {"NORMAL": 0, "WATCH": 1, "ALERT": 2, "ESCALATE": 3}
        current = level_rank.get(decision.level.value, 0)
        payload = {
            "agent_id": Config.AGENT_ID,
            "agent_name": Config.AGENT_NAME,
            "version": Config.VERSION,
            "decision_id": decision.decision_id,
            "timestamp": decision.timestamp,
            "level": decision.level.value,
            "reason": decision.reason,
            "coordination_signal": decision.coordination_signal,
            "phi": decision.phi,
            "substrate_energy": decision.substrate_energy,
            "energy_trend": decision.energy_trend,
            "live_station_count": decision.live_station_count,
        }

        with self._lock:
            subs = list(self.subscribers.values())

        for sub in subs:
            if current >= level_rank.get(sub.min_level, 0):
                if sub.webhook_url:
                    ok = self._fire_webhook(sub.webhook_url, payload)
                    status = "DELIVERED" if ok else "FAILED"
                    fired.append(f"{status} -> {sub.name} ({sub.role})")
                else:
                    fired.append(f"SIGNALED -> {sub.name} ({sub.role})")
                self.action_log.append(
                    {
                        "ts": time.time(),
                        "decision_id": decision.decision_id,
                        "subscriber": sub.name,
                        "role": sub.role,
                        "status": "SENT" if sub.webhook_url else "LOGGED",
                        "level": decision.level.value,
                    }
                )

        self.action_log.append(
            {
                "ts": time.time(),
                "decision_id": decision.decision_id,
                "subscriber": "SYSTEM",
                "role": "core",
                "status": "RECORDED",
                "level": decision.level.value,
            }
        )
        if len(self.action_log) > 250:
            self.action_log = self.action_log[-250:]
        return fired

    def _fire_webhook(self, url: str, payload: Dict[str, Any]) -> bool:
        try:
            data = json.dumps(payload).encode("utf-8")
            req = urllib.request.Request(
                url,
                data=data,
                method="POST",
                headers={"Content-Type": "application/json", "User-Agent": f"SAEONYX-Caribbean-Agent/{Config.VERSION}"},
            )
            with urllib.request.urlopen(req, timeout=6) as resp:
                return 200 <= resp.status < 300
        except Exception as e:
            log.warning("Webhook delivery failed (%s): %s", url, e)
            return False

    def register_subscriber(
        self, name: str, webhook_url: Optional[str] = None, min_level: str = "ALERT", role: str = "general"
    ) -> str:
        sid = f"sub-{uuid.uuid4().hex[:8]}"
        sub = Subscriber(
            subscriber_id=sid,
            name=name,
            webhook_url=webhook_url,
            min_level=min_level.upper(),
            registered_at=time.time(),
            role=role,
        )
        with self._lock:
            self.subscribers[sid] = sub
        log.info("Subscriber registered: %s (%s) role=%s min_level=%s", name, sid, role, min_level)
        return sid

    def start_background_loop(self) -> None:
        def loop() -> None:
            self.running = True
            log.info("Caribbean coordination agent loop started (interval %ds)", Config.POLL_INTERVAL)
            while self.running:
                try:
                    self.run_cycle()
                except Exception as e:
                    log.error("Agent cycle error: %s", e)
                time.sleep(Config.POLL_INTERVAL)

        threading.Thread(target=loop, daemon=True).start()


async def persist_decision_and_actions(dbm: DatabaseManager, decision: AgentDecision, actions: List[Dict[str, Any]]) -> None:
    await dbm.log_decision(decision)
    for action in actions:
        await dbm.log_action(action)


class SubscribeRequest(BaseModel):
    name: str
    webhook_url: Optional[str] = None
    min_level: str = "ALERT"
    role: str = "general"


app = FastAPI(
    title=Config.AGENT_NAME,
    version=Config.VERSION,
    description="Regional real-time safety and coordination infrastructure for the Caribbean. Substrate-driven. Real NOAA only.",
)
agent = CaribbeanOceanAgent()
dbm = DatabaseManager(Config.DATABASE_URL)


@asynccontextmanager
async def lifespan(app: FastAPI):
    await dbm.initialize()
    for name, level, role in [
        ("Port Authority Network", "ALERT", "ports"),
        ("Commercial Shipping Advisory", "ALERT", "shipping"),
        ("Tourism & Ferry Operators", "WATCH", "tourism"),
        ("CDEMA / Disaster Coordination", "ESCALATE", "disaster"),
    ]:
        sid = agent.register_subscriber(name, min_level=level, role=role)
        await dbm.upsert_subscriber(agent.subscribers[sid])
    agent.start_background_loop()
    log.info("SAEONYX Caribbean Ocean Agent v%s online", Config.VERSION)
    yield
    agent.running = False


app.router.lifespan_context = lifespan

UI_HTML = """<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>SAEONYX Caribbean Ocean Coordination Agent</title>
<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"/>
<script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:'Courier New',monospace;background:#040b14;color:#9ab8d0}
.app{display:flex;flex-direction:column;height:100vh}
.header{background:#0a1522;padding:12px 18px;border-bottom:1px solid #1a3048;display:flex;justify-content:space-between;align-items:center}
.header h1{font-size:16px;color:#c8e0f8;letter-spacing:0.5px}
.main{display:flex;flex:1;min-height:0}
.map-container{flex:1}
#map{width:100%;height:100%}
.sidebar{width:420px;background:#0a1522;border-left:1px solid #1a3048;padding:12px;overflow-y:auto}
.card{background:#0d1c2c;border:1px solid #1a3048;border-radius:6px;padding:10px;margin-bottom:8px}
.card .label{color:#5a7a98;font-size:10px;text-transform:uppercase;letter-spacing:1px}
.card .value{font-size:17px;color:#c0e0ff;font-weight:bold}
.level-NORMAL{color:#40c878}.level-WATCH{color:#e0c040}.level-ALERT{color:#e08040}.level-ESCALATE{color:#e04040}
.btn{background:#12253a;border:1px solid #2a4a68;color:#a0c8e8;padding:6px 12px;border-radius:4px;cursor:pointer;font-size:12px}
.btn:hover{background:#1a3550}.btn.primary{background:#1a3a5a;border-color:#3a6a9a}
.history{max-height:140px;overflow:auto;font-size:11px}.footer{padding:5px 18px;border-top:1px solid #12253a;font-size:10px;color:#3a5a78;text-align:center}
.status-dot{display:inline-block;width:8px;height:8px;border-radius:50%;margin-right:6px;background:#40d080}
pre{white-space:pre-wrap;word-break:break-word;font-size:10px;color:#7a9ab8}
</style>
</head>
<body>
<div class="app">
 <div class="header">
  <div>
   <h1>SAEONYX Caribbean Ocean Coordination Agent v3.1</h1>
   <div style="color:#5a7a98;font-size:11px">Substrate-driven regional infrastructure · Real NOAA only · <span id="ts"></span></div>
  </div>
  <div><span class="status-dot"></span><span style="font-size:12px">AGENT RUNNING</span></div>
 </div>
 <div class="main">
  <div class="map-container"><div id="map"></div></div>
  <div class="sidebar">
   <div class="card"><div class="label">Current Decision</div><div class="value" id="level">—</div><div style="font-size:12px;color:#8ab0c8;margin-top:4px" id="reason">Initializing substrate…</div></div>
   <div class="card"><div class="label">Phi / Energy / Trend / Live Stations</div><div style="display:flex;gap:14px;margin-top:4px;flex-wrap:wrap"><div><span class="value" id="phi" style="font-size:15px">—</span><div style="font-size:10px;color:#5a7a98">Phi</div></div><div><span class="value" id="energy" style="font-size:15px">—</span><div style="font-size:10px;color:#5a7a98">Energy</div></div><div><span class="value" id="trend" style="font-size:15px">—</span><div style="font-size:10px;color:#5a7a98">Trend</div></div><div><span class="value" id="live" style="font-size:15px">—</span><div style="font-size:10px;color:#5a7a98">Live</div></div></div></div>
   <div class="card"><div class="label">Coordination Signal</div><pre id="coord">—</pre></div>
   <div class="card"><div class="label">Actions Fired</div><div id="fired" style="font-size:11px;color:#8ab0c8">—</div></div>
   <div style="margin:8px 0;display:flex;gap:6px"><button class="btn primary" onclick="forceCycle()">Force Cycle</button><button class="btn" onclick="refreshStatus()">Refresh</button></div>
   <div class="card"><div class="label">Recent Decisions</div><div class="history" id="history">—</div></div>
  </div>
 </div>
 <div class="footer">SAEONYX Caribbean Ocean Coordination Agent · Substrate-driven · Real NOAA · Multi-party alert layer</div>
</div>
<script>
let map, markers = [];
function initMap() {
  map = L.map('map').setView([16.8, -70], 5);
  L.tileLayer('https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png', {maxZoom: 18}).addTo(map);
  L.rectangle([[10, -88], [27, -60]], {color: '#2a6a9a', weight: 1, fillOpacity: 0.04}).addTo(map);
}
async function refreshStatus() {
  try {
    const r = await fetch('/agent/status');
    const d = await r.json();
    const dec = d.latest_decision;
    if (dec) {
      const el = document.getElementById('level');
      el.textContent = dec.level;
      el.className = 'value level-' + dec.level;
      document.getElementById('reason').textContent = dec.reason;
      document.getElementById('phi').textContent = dec.phi;
      document.getElementById('energy').textContent = dec.substrate_energy;
      document.getElementById('trend').textContent = (dec.energy_trend >= 0 ? '+' : '') + dec.energy_trend;
      document.getElementById('live').textContent = dec.live_station_count;
      document.getElementById('coord').textContent = JSON.stringify(dec.coordination_signal, null, 2);
      document.getElementById('fired').innerHTML = (dec.actions_fired || []).map(a => `<div>${a}</div>`).join('') || '—';
    }
    document.getElementById('ts').textContent = new Date().toLocaleTimeString();
    const hist = d.history || [];
    document.getElementById('history').innerHTML = hist.slice().reverse().slice(0, 7).map(h =>
      `<div style="padding:3px 0;border-bottom:1px solid #12253a"><span class="level-${h.level}">${h.level}</span> · ${new Date(h.timestamp * 1000).toLocaleTimeString()} · E=${h.substrate_energy}</div>`
    ).join('') || '—';
    if (d.sensors) {
      markers.forEach(m => map.removeLayer(m));
      markers = [];
      d.sensors.forEach(s => {
        const isLive = s.status === 'live' || s.status === 'partial';
        const color = isLive ? '#4a9ad0' : '#555';
        const m = L.circleMarker([s.lat, s.lng], { radius: isLive ? 8 : 5, color, fillColor: color, fillOpacity: 0.85 }).addTo(map);
        m.bindPopup(`<b>${s.name}</b> (${s.station_id})<br>Status: ${s.status}<br>SST: ${s.sst_c ?? '—'}°C<br>Waves: ${s.wave_height_m ?? '—'} m<br>Wind: ${s.wind_speed_ms ?? '—'} m/s`);
        markers.push(m);
      });
    }
  } catch (e) {
    console.error(e);
  }
}
async function forceCycle() {
  await fetch('/agent/cycle', {method: 'POST'});
  setTimeout(refreshStatus, 900);
}
document.addEventListener('DOMContentLoaded', () => {
  initMap();
  refreshStatus();
  setInterval(refreshStatus, 11000);
});
</script>
</body>
</html>"""


@app.get("/", response_class=HTMLResponse)
async def root() -> str:
    return UI_HTML


@app.get("/agent/status")
async def agent_status() -> Dict[str, Any]:
    with agent._lock:
        latest = agent.latest_decision
        history = list(agent.decision_history)
        action_log = list(agent.action_log[-15:])
        subs = {k: asdict(v) for k, v in agent.subscribers.items()}
    sensors = get_live_sensors()
    return {
        "agent_id": Config.AGENT_ID,
        "version": Config.VERSION,
        "status": "RUNNING" if agent.running else "STOPPED",
        "cycles": agent._cycles,
        "phi": round(agent._phi, 4),
        "latest_decision": asdict(latest) if latest else None,
        "history": [asdict(h) for h in history[-12:]],
        "action_log": action_log,
        "subscribers": subs,
        "sensors": sensors,
    }


@app.post("/agent/cycle")
async def force_cycle() -> Dict[str, Any]:
    before_count = len(agent.action_log)
    decision = agent.run_cycle()
    new_actions = agent.action_log[before_count:]
    await persist_decision_and_actions(dbm, decision, new_actions)
    return {"status": "ok", "decision": asdict(decision)}


@app.post("/agent/subscribe")
async def subscribe(req: SubscribeRequest) -> Dict[str, Any]:
    sid = agent.register_subscriber(req.name, req.webhook_url, req.min_level, req.role)
    await dbm.upsert_subscriber(agent.subscribers[sid])
    return {
        "subscriber_id": sid,
        "name": req.name,
        "min_level": req.min_level,
        "role": req.role,
        "webhook_url": req.webhook_url,
    }


@app.get("/agent/subscribers")
async def list_subscribers() -> Dict[str, Any]:
    with agent._lock:
        return {k: asdict(v) for k, v in agent.subscribers.items()}


@app.get("/agent/actions")
async def list_actions(limit: int = 40) -> List[Dict[str, Any]]:
    with agent._lock:
        return agent.action_log[-limit:]


@app.get("/agent/decisions")
async def list_decisions(limit: int = 25) -> List[Dict[str, Any]]:
    with agent._lock:
        return [asdict(d) for d in agent.decision_history[-limit:]]


@app.get("/health")
async def health() -> Dict[str, Any]:
    return {
        "status": "ONLINE",
        "agent_id": Config.AGENT_ID,
        "version": Config.VERSION,
        "agent_running": agent.running,
        "cycles": agent._cycles,
        "phi": round(agent._phi, 4),
        "latest_level": agent.latest_decision.level.value if agent.latest_decision else None,
        "live_stations": agent.latest_decision.live_station_count if agent.latest_decision else 0,
        "subscribers": len(agent.subscribers),
    }


@app.get("/admin/status")
async def admin_status(token: Optional[str] = Header(None, alias="admin-token")) -> Dict[str, Any]:
    if not Config.ADMIN_TOKEN or token != Config.ADMIN_TOKEN:
        raise HTTPException(401, "Invalid admin token")
    return {
        "substrate_cycles": agent._cycles,
        "phi": round(agent._phi, 4),
        "soul": agent._soul.as_dict,
        "decision_count": len(agent.decision_history),
        "action_count": len(agent.action_log),
        "subscribers": len(agent.subscribers),
    }


if __name__ == "__main__":
    print("\n" + "=" * 72)
    print(" SAEONYX Caribbean Ocean Coordination Agent v3.1")
    print(" Substrate-driven regional infrastructure")
    print(" Real NOAA NDBC only · Zero simulation")
    print(" Multi-party coordination layer active")
    print(f" Open -> http://localhost:{Config.PORT}")
    print("=" * 72 + "\n")
    uvicorn.run(app, host=Config.HOST, port=Config.PORT)
