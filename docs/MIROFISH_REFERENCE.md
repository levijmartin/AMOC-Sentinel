# MiroFish Reference

## Local source reviewed
`workspace/MicroFish/MiroFish-main`

## What MiroFish appears to be
MiroFish presents itself as a multi-agent prediction and simulation engine. Based on the repository materials reviewed, it emphasizes:
- graph building from seed materials
- environment and persona setup
- simulation with many interacting agents
- report generation after simulation
- deep interaction with agents and outputs

It also uses a practical full-stack pattern:
- frontend and backend separation
- Python backend
- Node-based frontend
- environment-variable-driven model configuration
- Docker and source deployment paths

## Useful patterns for Future Caribbean

### 1. Workflow decomposition
MiroFish breaks the system into distinct stages:
- graph building
- environment setup
- simulation
- report generation
- deep interaction

For StormShield Caribbean, this supports a similar staged mindset:
- ingestion and normalization
- signal analysis
- climate-risk reasoning
- report generation
- operator interaction

### 2. Multi-agent framing
MiroFish is a reminder that the value is not just a single model call, but a structured system of roles and stages. This supports StormShield Caribbean's agentic approach where different responsibilities are separated instead of collapsing everything into one prompt.

### 3. Report-first output design
MiroFish highlights the importance of producing a tangible output artifact after reasoning. For StormShield Caribbean, the equivalent is a resilience report, alert summary, or operator-ready action plan.

### 4. Practical deployment shape
The repository shows a simple, practical developer experience:
- clear setup scripts
- separate backend and frontend
- Docker support
- environment-based model configuration

This is a useful implementation pattern to keep in mind for a buildathon MVP.

## What not to copy directly
StormShield Caribbean should not inherit all of MiroFish's product assumptions.

Avoid blindly copying:
- social-simulation-heavy framing
- persona-generation logic that does not fit climate-risk workflows
- unnecessary complexity in the first MVP

## Recommendation
Treat MiroFish as a reference for:
- staged agent workflow design
- simulation/report pipeline thinking
- practical repo/deployment organization

Do not treat it as the core domain architecture. StormShield Caribbean should remain grounded in environmental data, climate/ocean reasoning, and operator-facing resilience decisions.
