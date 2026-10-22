# Learning Progress: FastAPI–SQLAlchemy–Postgres

> Guide: [FASTAPI_SQLALCHEMY_POSTGRES_GUIDE.md](./FASTAPI_SQLALCHEMY_POSTGRES_GUIDE.md)
> Last updated: 2026-10-22 16:03 +07:00

## Status

| Field | Value |
| --- | --- |
| Overall status | in progress |
| Accumulated reported active time | 0 minutes |
| Current unit | LG-04 — Configure Postgres access |
| Next action | Recheck that `postgresql` names the database dialect, `asyncpg` names the async Python driver, and `engine.connect()` is the first real connection attempt; then explain which line fails when Postgres is down. |

## Unit checklist

Use only `not started`, `in progress`, `completed`, `blocked`, `deferred`, or `skipped`. Use `skipped` for a learner-confirmed omission and `deferred` for a planned return. Check a box only for `completed`.

- [x] LG-01 — Create the uv project — `completed`
- [x] LG-02 — Run the first FastAPI endpoint — `completed`
- [x] LG-03 — Define the todo API contract — `completed`
- [ ] LG-04 — Configure Postgres access — `in progress`
- [ ] LG-05 — Map the todo table — `not started`
- [ ] LG-06 — Apply the initial migration — `not started`
- [ ] LG-07 — Give each request an async session — `not started`
- [ ] LG-08 — Create a todo — `not started`
- [ ] LG-09 — List and read todos — `not started`
- [ ] LG-10 — Update and delete a todo — `not started`
- [ ] LG-11 — Handle failure paths — `not started`
- [ ] LG-12 — Check the API with Pytest — `not started`
- [ ] LG-13 — Run quality checks and final demo — `not started`

## Current checkpoint

### Last demonstrated evidence

- LG-01: learner explained the roles of `pyproject.toml`, `uv.lock`, and `.python-version`, then correctly explained that `uv run --locked` checks the lockfile, prepares `.venv`, and runs the command. The project command printed FastAPI 0.159.1 and Uvicorn 0.54.0.
- LG-02: learner correctly identified that `@app.get("/health")` maps `health()` to `GET /health`.
- LG-02: learner correctly identified the `{"status":"ok"}` JSON body and HTTP 200 status.
- LG-02: learner saw Swagger UI list `GET /health` and explained that the decorator registers the route at startup while `health()` runs on a request.
- LG-03: learner correctly identified create and update input models, their defaults and title constraint, and explained that `TodoRead` omits an extra field because it is not declared in the response model.

### Recent misconceptions or fragile knowledge

- None observed.

### Incomplete knowledge or work

- LG-04: learner authorized Docker Compose. Added `compose.yaml`, ignored local `.env`, `.env.example`, and `database.py`; installed SQLAlchemy 2.0.55 with its asyncio extra and asyncpg 0.32.0. Docker Desktop and the Postgres 17 container started; Compose health check passed. `make db-check` ran `SELECT 1` through the async engine and returned 1. Learner identified `.env` as the credential source and named `check_connection()` as the failure area, but still needs to identify `asyncpg` as the async Python driver and `engine.connect()` as the precise first connection attempt. The container was stopped with `docker compose down`; its named volume remains.

### Open questions

- None.

### Retrieval prompts

- Can the learner distinguish engine construction from `engine.connect()` and identify the async driver?

### Exact re-entry prompt

> In `database.py`, what does each part of `postgresql+asyncpg` name? If Postgres is down, does `create_async_engine(...)` or `engine.connect()` fail first, and why?

## Adaptations

- None.

## Session log

### Session 1 — 2026-10-22 (paused)

- Planned time: 2 hours.
- Actual active time: not reported.
- Units worked: LG-01 completed; LG-02 in progress.
- Learner evidence: explained uv project files and `uv run --locked`; identified the route decorator mapping and the `{"status":"ok"}` HTTP 200 response.
- Misconceptions or uncertainty: Python decorator timing and the generated `/docs` entry remain unverified.
- Unresolved work: inspect `/docs`, identify `GET /health`, and explain definition-time registration versus request-time execution.
- Research notes: checked current official uv, FastAPI, Python, and TypeScript documentation. Local server and Makefile targets were verified; server stopped at pause.
- Next prompt: From `experiments/python/fastapi-sqlalchemy-postgres`, run `make dev` and open http://127.0.0.1:8000/docs. What method and path appear there, and does `@app.get("/health")` register `health()` when the module loads or call `health()` immediately?

### Session 2 — 2026-10-22 (paused)

- Planned time: 45 minutes.
- Actual active time: not reported.
- Units worked: LG-02 and LG-03 completed; LG-04 in progress.
- Learner evidence: identified `GET /health` in Swagger UI and decorator timing; explained create/update/read model roles and why undeclared response fields are omitted; identified `.env` as the credential source.
- Misconceptions or uncertainty: still needs to name `asyncpg` as the async Python driver and distinguish engine construction from the first connection attempt.
- Unresolved work: explain the URL's dialect and driver, and point to `engine.connect()` as the first operation that fails when Postgres is down.
- Demonstration: Postgres 17 via Docker Compose was healthy; `make db-check` returned `SELECT 1 returned 1`; Ruff passed. The container was stopped while preserving its data volume.
- Next prompt: In `database.py`, what does each part of `postgresql+asyncpg` name? If Postgres is down, does `create_async_engine(...)` or `engine.connect()` fail first, and why?
