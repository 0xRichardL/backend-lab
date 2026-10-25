# Learning Progress: FastAPI–SQLAlchemy–Postgres

> Guide: [FASTAPI_SQLALCHEMY_POSTGRES_GUIDE.md](./FASTAPI_SQLALCHEMY_POSTGRES_GUIDE.md)
> Last updated: 2026-10-26 +07:00

## Status

| Field | Value |
| --- | --- |
| Overall status | in progress |
| Accumulated reported active time | 0 minutes |
| Current unit | LG-11 — Handle failure paths (`not started`) |
| Next action | Begin LG-11 with validation `422`, missing-row `404`, and a failed write that leaves data unchanged; show every shell command before running it. |

## Unit checklist

Use only `not started`, `in progress`, `completed`, `blocked`, `deferred`, or `skipped`. Use `skipped` for a learner-confirmed omission and `deferred` for a planned return. Check a box only for `completed`.

- [x] LG-01 — Create the uv project — `completed`
- [x] LG-02 — Run the first FastAPI endpoint — `completed`
- [x] LG-03 — Define the todo API contract — `completed`
- [x] LG-04 — Configure Postgres access — `completed`
- [x] LG-05 — Map the todo table — `completed`
- [x] LG-06 — Apply the initial migration — `completed`
- [x] LG-07 — Give each request an async session — `completed`
- [x] LG-08 — Create a todo — `completed`
- [x] LG-09 — List and read todos — `completed`
- [x] LG-10 — Update and delete a todo — `completed`
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
- LG-06: the reviewed initial revision was applied to fresh local Postgres. `\d todos` showed `id` as an integer primary key, `title` as non-null text, and `completed` as non-null boolean with default `false`. `alembic_version` recorded `8d8b678efb5f`; `todos` had zero rows.
- LG-06: the migration, notes, and repeatable `alembic upgrade head` README command are committed. On 2026-10-25, Ruff lint and format checks passed after limiting three style-rule exceptions to generated Alembic revision files.
- LG-07: learner correctly explained that `Depends(get_session)` needs no NestJS-style provider registration and that two concurrent requests receive two distinct sessions.
- LG-07: learner correctly explained that A's session closes when A finishes while B's session and the shared engine remain available.
- LG-07: learner identified `session_factory()` inside `async with` as the expression that creates each request's distinct session. Two `/db-check` requests against local Postgres returned HTTP 200 and `{"connected": true}`; temporary instrumentation recorded different `AsyncSession` objects.
- LG-08: learner correctly identified `await session.commit()` as the durability boundary and `await session.refresh(todo)` as reloading stored row values into the object.
- LG-08: the in-process `POST /todos` with `{"title":"Buy milk"}` returned HTTP `201` and `{"id":1,"title":"Buy milk","completed":false}`; a separate Postgres query returned `(1, "Buy milk", false)`. Learner identified the Pydantic default as the source of `completed=false` and Postgres as the source of `id=1`.
- LG-09: learner correctly chose `.order_by(Todo.id)` for a stable ascending collection. Ascending is the default direction.
- LG-09: learner correctly explained that the query determines the collection order and the list comprehension maps each stored `Todo` to a new `TodoRead` item.
- LG-09: learner identified `session.get(Todo, 1)` as returning the mapped todo; the route explicitly constructs `TodoRead` and returned the matching JSON. Learner asked about other handler return types and chose to skip an extra response-model check after the explanation.
- LG-10: learner correctly distinguished omitted `title` from explicit `null` after a local Pydantic check. A live PATCH changed stored row 5 from `(5, "Keep title", false)` to `(5, "Keep title", true)`; explicit null returned `422`. DELETE returned `204` with an empty body, and a separate query found no row 5.

### Recent misconceptions or fragile knowledge

- LG-10: learner initially described `exclude_unset=True` as removing all null-value fields; clarified and verified that it removes omitted fields, while explicit `null` remains `None`. No LG-10 work remains.
- LG-09: learner initially treated `[3,1,2]` as the natural guaranteed order of `select(Todo)`. Clarified that without `ORDER BY`, SQL does not guarantee a row order.
- LG-09: learner requested clarification of the Python list comprehension used to convert mapped `Todo` objects to `TodoRead` objects.
- LG-09: learner correctly identified `session.get(Todo, 1)` as returning a mapped `Todo`, but guessed that the route returns that object directly. Clarified the explicit `TodoRead(...)` construction before FastAPI response processing.
- LG-09: learner asked what happens when a handler returns a different type; compatible dict filtering versus response validation failure was explained. Learner chose to skip the additional check.

### Incomplete knowledge or work

- No LG-10 work remains. The live before-and-after row checks confirmed the partial PATCH and deletion behavior.
- No LG-09 work remains. In-process reads against local Postgres returned HTTP `200`: the list contained IDs `[1,2,3,4]` in ascending order; `/todos/1` returned `{"id":1,"title":"Buy milk","completed":false}`. Row `2` already existed when this demonstration ran.
- No LG-08 work remains. The route and reference note are in place; the create request and stored row matched. `make dev` could not bind a port in this sandbox, so the request was exercised through FastAPI's in-process `TestClient` with host-permission access to local Postgres.
- No LG-07 work remains. The shared engine/factory, per-request dependency, `/db-check` route, and `make dev` environment loading are in place. The lifespan return annotation is `AsyncGenerator[None]`, correcting a reported deprecation warning for `AsyncIterator` with `@asynccontextmanager`.
- No LG-05 work remains. Read-only inspection found `todos` and `alembic_version` in local Postgres. `alembic_version` contained `2628f7aaa40d`; `todos` had `id`, `title`, and `completed` columns and zero rows. The migration files were removed in the discarded demonstration. The learner requested a fresh start, and `docker compose down --volumes` removed this experiment's container and named Postgres volume. `make db-up` then created a fresh healthy container and volume; read-only inspection of `public` returned zero tables. `uv add --dev alembic` added Alembic 1.20.1 to the experiment and updated the lockfile. `uv run --locked alembic init -t async alembic` generated `alembic.ini`, `alembic/env.py`, `alembic/script.py.mako`, `alembic/README`, and an empty `alembic/versions` directory; all generated file contents were shown before editing.

### Open questions

- None.

### Retrieval prompts

- None for LG-05; its mastery evidence was demonstrated.

### Exact re-entry prompt

> LG-10 is complete. The learner distinguished omitted title from explicit null, and live Postgres inspection confirmed PATCH preserved title while changing completed, then DELETE removed the row and returned `204` with no body. Begin LG-11 with request validation `422`, missing-row `404`, and a failed write that leaves data unchanged. Show every shell command before running it.

## Adaptations

- For LG-06, show one file or command change and its observed result per turn, then stop for discussion before the next change. Show generated files before editing them. Skip generated lockfile contents during the file walkthrough; its role is already understood.
- After each Learning Item's mastery evidence is demonstrated and its relevant checks pass, commit that item's files and progress update before starting the next Learning Item. Keep commits scoped to the completed item.
- The learner knows the general migration-tool pattern and asked to skip generic concept-check questions. Focus LG-06 on the specific Alembic files, commands, and observed Postgres result; verify mastery through the demonstrated artifact and a concise explanation only where it adds value.
- Show every shell command to the learner before running it, then show its relevant output.
- On 2026-10-25, the learner chose a narrow Ruff per-file exception for Alembic-generated revision style instead of editing the generated revision.
- In LG-07, the learner asked for detailed explanations of `yield`, FastAPI dependency resolution, and the relevant NestJS comparison before code; the lesson used those explanations without changing the guide's order.
- In LG-09, the learner chose to skip an extra response-model quiz after receiving the explanation. The item was completed based on its required read-query evidence and live responses.

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

- Planned time: started around 8:24 AM with an initial 10:00 AM cutoff; learner extended the same session by 90 minutes at 10:22 AM, giving a new cutoff around 11:52 AM.
- Actual active time: not reported.
- Units worked: LG-04 and LG-05 completed; LG-06 remains not started in the current checkpoint.
- Learner evidence: distinguished `postgresql` dialect from `asyncpg` driver and engine creation from first connection; identified `server_default=text("false")` as the Postgres default for a direct insert that omits `completed`.
- Demonstration and file state: `models.py` and `TODO_TABLE_MAPPING.md` remain. The discarded migration's Postgres volume was removed at the learner's request, and fresh Postgres has an empty `public` schema. Alembic was added, its async starter files were generated and configured, and the setup was committed as `d2270bf`. No new revision has been generated or applied yet.
- Teaching preference: work through each file and each change separately, showing the observed result before continuing. Generated lockfile contents do not need a walkthrough.
- Next prompt: Show and run `uv run --locked --env-file .env alembic revision --autogenerate -m "create todos"`; inspect the generated revision before applying it. Show every command before running it.

### Session 4 — 2026-10-25 (ended)

- Planned time: 45 minutes.
- Actual active time: not reported.
- Unit worked: LG-06 completed.
- Learner evidence: prior reviewed migration and Postgres inspection showed the expected initial `todos` schema, revision row, and zero todo rows. The learner chose to ignore generated revision style findings through a scoped Ruff config.
- Verification: `UV_CACHE_DIR=/tmp/backend-lab-uv-cache make lint` passed; `UV_CACHE_DIR=/tmp/backend-lab-uv-cache make format-check` passed (11 files already formatted). The temporary uv cache path was needed because the sandbox could not access the default cache.
- Next prompt: Begin LG-07 by tracing the existing engine and a per-request async session. Show every command before running it and continue one step per turn.

### Session 5 — 2026-10-25 (ended)

- Planned time: 25 minutes.
- Actual active time: not reported.
- Unit worked: LG-07 completed.
- Learner evidence: explained automatic dependency resolution, distinct sessions for two requests, A's cleanup while B remains active, and `session_factory()` as the session creation point.
- Demonstration: two `/db-check` requests against local Postgres returned HTTP 200 with `{"connected": true}`; temporary instrumentation confirmed distinct session objects. Ruff lint and format checks passed. The local Docker container was stopped afterward, preserving its data volume.
- Next prompt: Begin LG-08 by tracing POST input into the mapped todo object. Show every command before running it.
