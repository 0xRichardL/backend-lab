# Learning Progress: FastAPI–SQLAlchemy–Postgres

> Guide: [FASTAPI_SQLALCHEMY_POSTGRES_GUIDE.md](./FASTAPI_SQLALCHEMY_POSTGRES_GUIDE.md)
> Last updated: 2026-10-23 09:37 +07:00

## Status

| Field | Value |
| --- | --- |
| Overall status | in progress |
| Accumulated reported active time | 0 minutes |
| Current unit | LG-06 — Apply the initial migration (`in progress`) |
| Next action | Run the shown autogenerate command with `.env`, then show the generated revision before editing or applying it. |

## Unit checklist

Use only `not started`, `in progress`, `completed`, `blocked`, `deferred`, or `skipped`. Use `skipped` for a learner-confirmed omission and `deferred` for a planned return. Check a box only for `completed`.

- [x] LG-01 — Create the uv project — `completed`
- [x] LG-02 — Run the first FastAPI endpoint — `completed`
- [x] LG-03 — Define the todo API contract — `completed`
- [x] LG-04 — Configure Postgres access — `completed`
- [x] LG-05 — Map the todo table — `completed`
- [ ] LG-06 — Apply the initial migration — `in progress`
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
- LG-04: learner correctly explained that `create_async_engine(...)` does not open a connection and `engine.connect()` fails if Postgres is unavailable.
- LG-04: learner identified `asyncpg` as the driver and, after clarification, `postgresql` as the PostgreSQL dialect in the URL. The successful `SELECT 1` connection check and credential source were shown in the prior session.
- LG-05: learner correctly identified `server_default=text("false")` as the source of Postgres `DEFAULT false` for a direct insert that omits `completed`, distinct from the Pydantic input default.

### Recent misconceptions or fragile knowledge

- None currently observed.

### Incomplete knowledge or work

- No LG-05 work remains. Read-only inspection found `todos` and `alembic_version` in local Postgres. `alembic_version` contained `2628f7aaa40d`; `todos` had `id`, `title`, and `completed` columns and zero rows. The migration files were removed in the discarded demonstration. The learner requested a fresh start, and `docker compose down --volumes` removed this experiment's container and named Postgres volume. `make db-up` then created a fresh healthy container and volume; read-only inspection of `public` returned zero tables. `uv add --dev alembic` added Alembic 1.20.1 to the experiment and updated the lockfile. `uv run --locked alembic init -t async alembic` generated `alembic.ini`, `alembic/env.py`, `alembic/script.py.mako`, `alembic/README`, and an empty `alembic/versions` directory; all generated file contents were shown before editing.

### Open questions

- None.

### Retrieval prompts

- None for LG-05; its mastery evidence was demonstrated.

### Exact re-entry prompt

> The fresh Postgres `public` schema has zero tables. Alembic 1.20.1 is a development dependency. The async starter files were generated and shown. `alembic/env.py` now uses `database_url()` and `Base.metadata`, and its full contents were shown. Next run `uv run --locked --env-file .env alembic revision --autogenerate -m "create todos"`, then show the generated revision before editing or applying it. Continue one step per turn and show every command before running it.

## Adaptations

- For LG-06, show one file or command change and its observed result per turn, then stop for discussion before the next change. Show generated files before editing them. Skip generated lockfile contents during the file walkthrough; its role is already understood.
- After each Learning Item's mastery evidence is demonstrated and its relevant checks pass, commit that item's files and progress update before starting the next Learning Item. Keep commits scoped to the completed item.
- The learner knows the general migration-tool pattern and asked to skip generic concept-check questions. Focus LG-06 on the specific Alembic files, commands, and observed Postgres result; verify mastery through the demonstrated artifact and a concise explanation only where it adds value.
- Show every shell command to the learner before running it, then show its relevant output.

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

### Session 3 — 2026-10-23 (in progress)

- Planned time: available until 10:00 AM, starting around 8:24 AM (about 1 hour 36 minutes).
- Actual active time: not reported.
- Units worked: LG-04 and LG-05 completed; LG-06 remains not started in the current checkpoint.
- Learner evidence: distinguished `postgresql` dialect from `asyncpg` driver and engine creation from first connection; identified `server_default=text("false")` as the Postgres default for a direct insert that omits `completed`.
- Demonstration and file state: `models.py` and `TODO_TABLE_MAPPING.md` remain. An initial Alembic revision was applied to local Postgres, but its setup files, dependency, Makefile target, and completed-migration documentation were removed. The local `todos` and `alembic_version` tables remain to be inspected before LG-06 proceeds.
- Teaching preference: work through each file and each change separately, showing the observed result before continuing. Generated lockfile contents do not need a walkthrough.
- Next prompt: Inspect the existing local Postgres `todos` and `alembic_version` tables; show the result and stop before changing the database or files.
