# AMOC Sentinel — Full System Mermaid Diagram

## Purpose
This document preserves the fuller system diagram in Mermaid form so it can be reused in the repo, application materials, or future planning discussions.

This is best treated as a broader future-state / expanded-system diagram rather than the smallest buildathon MVP topology.

## Mermaid diagram

```mermaid
flowchart TD
    A[AMOC Data/APIs<br/>RAPID • Copernicus • Argo]
    --> B[Copernicus Python/API Subset + Download]

    B --> C[IoT / PicoClaw Edge Packets]
    C --> D[Four-Operator Field-to-Bit Filter]
    D --> E[QARTOD QC<br/>pass • suspect • fail • missing]
    E --> F{Usable + Trusted?}

    F -- No --> G[Quarantine + Expert Review]
    F -- Yes --> H[MWP-Audited LangGraph Agents]

    H --> I[AMOC / Cold Blob Current-State Model]
    I --> J[Caribbean Hazard Translation<br/>sargassum • heat • sea level]
    J --> K[MiroFish Response Simulation]
    K --> L{Risk / Confidence Threshold?}

    L -- Low --> M[Weekly Ocean Risk Brief]
    L -- High --> N[Action Agent Drafts Advisory]
    N --> O[Human-in-the-Loop Expert Approval]

    O -- Reject --> P[Revise / Hold]
    O -- Approve --> Q[Stakeholder Hub + API Release]
    Q --> R[x402 Data / Commerce Rails<br/>datasets • simulations • API calls • premium briefs]
```

## Scope note
This diagram is intentionally broader than the MVP topology. It is useful for showing the long-term system direction, including:
- edge sensing
- packet filtering
- quality control and trust gates
- simulation
- expert review
- monetized distribution

For the buildathon MVP, the cleaner 3-agent topology should remain the primary architecture visual.
