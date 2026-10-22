# Learning Guide: FastAPI–SQLAlchemy–Postgres

> Created: 2026-10-22  
> Guide version: 1  
> Target slug: `fastapi-sqlalchemy-postgres`

## Learning contract

| Field | Value |
| --- | --- |
| Desired outcome | Follow a step-by-step demonstration that builds and verifies a todo API with Python, FastAPI, Pydantic, async SQLAlchemy, Alembic, Postgres, uv, Pytest, Pyright, and Ruff. |
| Current level | Senior Node.js and Go backend engineer; familiar with Python language features, new to this stack. |
| Success criteria | Explain the request-to-database-to-response path, the initial migration, and the transaction boundary; read the commented implementation; interpret the demonstrated API and quality-check results. |
| Total active-learning budget | 6 hours maximum; planned total is 4 hours 50 minutes. |
| Constraints | Demonstration only, with no learner coding or exercises. Teaching code requires comments. Cover one initial migration; omit migration changes and database test isolation. |

## Priority and scope

### Must learn

- uv project files, dependency locking, and command execution.
- FastAPI routes, request handling, dependencies, responses, and errors.
- Pydantic input and output models and runtime validation.
- Postgres access with SQLAlchemy models, async sessions, queries, and commits.
- The initial Alembic migration.
- Todo CRUD behavior.
- Pytest API checks and Pyright and Ruff checks.

### Supporting knowledge

- Configuration and database URLs.
- Engine lifetime versus request session lifetime.
- What a failed transaction does.
- Pydantic API models versus SQLAlchemy database models.

### Deliberately deferred

- Changing, reversing, or branching migrations: this guide covers the initial schema.
- Isolated test databases and transaction fixtures: excluded from this project's focus.
- Authentication, pagination, jobs, deployment, observability, and performance: outside the basic request-to-database path.
- Generic repository or service layers: unnecessary for this small demonstration.

## Roadmap

Each estimate includes a brief recap and inspection of the demonstrated result. Evidence is code, command output, an API response, or a read-only explanation—never a required hands-on exercise.

| ID | Unit | Observable outcome | Knowledge to learn | Demonstration and mastery evidence | Estimate |
| --- | --- | --- | --- | --- | ---: |
| LG-01 | Create the uv project | Identify what makes the project reproducible. | `pyproject.toml`, `uv.lock`, runtime and development dependencies, `uv add`, `uv run`. | Show the project files and an executed command; explain each file's role. | 15m |
| LG-02 | Run the first FastAPI endpoint | Trace a request to a route and response. | App instance, route decorator, `async def`, server command, `/docs`. | Run a commented health endpoint; inspect its response and API documentation. | 15m |
| LG-03 | Define the todo API contract | Distinguish accepted input from returned data. | Create, update, and read Pydantic models; required and optional fields; constraints; validation and serialization. | Compare valid and invalid JSON with the responses and the model used at each boundary. | 20m |
| LG-04 | Configure Postgres access | Explain how the app locates its database. | Database URL, async driver, configuration loading, connection errors. | Show configuration and a successful connectivity check; identify where credentials enter. | 20m |
| LG-05 | Map the todo table | Match a Python mapped class to database columns. | Declarative base, `Mapped`, `mapped_column`, ID, title, completed flag, defaults. | Read the commented model and match its fields to the expected table. | 20m |
| LG-06 | Apply the initial migration | Explain how the first table reaches Postgres. | Alembic environment, model metadata, generated revision, revision review, `upgrade head`. | Show the reviewed revision, migration output, and resulting table. | 20m |
| LG-07 | Give each request an async session | Trace database resource ownership. | Async engine, `async_sessionmaker`, FastAPI dependency with `yield`, session cleanup and lifetime. | Trace two requests; explain why they use distinct sessions and when each closes. | 25m |
| LG-08 | Create a todo | Trace a write from request JSON to stored row and response. | Pydantic input, mapped object, `add`, `commit`, `refresh`, `201 Created`, output model. | Compare the POST request, stored row, and returned JSON. | 25m |
| LG-09 | List and read todos | Trace queries and their results. | `select`, result extraction, list response, ID lookup, stable ordering. | Show collection and single-todo responses; explain the matching queries. | 25m |
| LG-10 | Update and delete a todo | Show how partial changes and deletion persist. | `PATCH`, omitted versus supplied fields, mapped-object changes, `delete`, commit, `204 No Content`. | Show before-and-after rows; confirm an omitted field remains unchanged. | 25m |
| LG-11 | Handle failure paths | Explain client errors and failed writes. | Validation `422`, missing-row `404`, `HTTPException`, rollback, public versus internal error. | Inspect both error responses and a failed write that leaves data unchanged. | 25m |
| LG-12 | Check the API with Pytest | Read tests for the public contract. | Assertions, FastAPI `TestClient`, focused dependency override, success and failure cases. | Show passing checks for representative create/read behavior, invalid input, missing todo, and deletion. No isolated database fixture is introduced. | 30m |
| LG-13 | Run quality checks and final demo | Reproduce the verification sequence. | Pyright configuration, Ruff lint and format checks, `uv run`, live Postgres smoke path. | Show Pytest, Pyright, and Ruff results, then create → read → update → delete against Postgres. | 25m |

**Planned total: 4 hours 50 minutes.**

## Milestones

| Milestone | Units | Cumulative time | Evidence |
| --- | --- | ---: | --- |
| Project and API contract | LG-01–LG-03 | 50m | The app runs and its data shapes are visible. |
| Persistence path | LG-04–LG-07 | 2h 15m | The table exists and request-scoped sessions reach Postgres. |
| Complete todo behavior | LG-08–LG-11 | 3h 55m | CRUD and selected failures have observable API and database results. |
| Verified demonstration | LG-12–LG-13 | 4h 50m | Checks pass and the API sequence works against Postgres. |

## Compact reference

### Three representations of a todo

- **Pydantic input model:** what the client may send.
- **SQLAlchemy mapped class:** how the row is stored.
- **Pydantic output model:** what the API returns.

### Database lifetime

- The async engine and session factory are shared application resources.
- Each request's work uses its own `AsyncSession`; one session is not shared across concurrent requests.

### Schema versus rows

- Alembic applies the initial **schema**.
- SQLAlchemy sessions read and change **rows**.
- A generated migration is a candidate to review before applying.

### Verification

- Pytest checks observable API behavior.
- Pyright checks type usage.
- Ruff checks lint and formatting.
- The Postgres demonstration checks the integrated path.

## Sources

1. [uv project guide](https://docs.astral.sh/uv/guides/projects/) — project files, dependencies, locking, and `uv run`.
2. [FastAPI tutorial](https://fastapi.tiangolo.com/tutorial/) — routes, requests, responses, and generated documentation.
3. [Pydantic models](https://pydantic.dev/docs/validation/latest/concepts/models/) — model validation and serialization.
4. [FastAPI dependencies with `yield`](https://fastapi.tiangolo.com/tutorial/dependencies/dependencies-with-yield/) — request dependency setup and cleanup.
5. [SQLAlchemy async documentation](https://docs.sqlalchemy.org/en/20/orm/extensions/asyncio.html) — async engines, sessions, and session concurrency rules.
6. [Alembic tutorial](https://alembic.sqlalchemy.org/en/latest/tutorial.html) — migration environments, revisions, and upgrades.
7. [Alembic asyncio cookbook](https://alembic.sqlalchemy.org/en/latest/cookbook.html#using-asyncio-with-alembic) — async migration setup.
8. [FastAPI testing](https://fastapi.tiangolo.com/tutorial/testing/) — API checks with `TestClient`.
9. [Pytest getting started](https://docs.pytest.org/en/stable/getting-started.html) — test discovery and assertions.
10. [Pyright configuration](https://github.com/microsoft/pyright/blob/main/docs/configuration.md) and [Ruff tutorial](https://docs.astral.sh/ruff/tutorial/) — type and lint checks.

## Research limits

- This is a plan, not a verified implementation. The demonstration will record exact compatible package versions and commands.
- Pyright's official CLI is distributed through npm; its Python package is a community-maintained wrapper. The demonstration will identify which route it uses.
- Documentation was checked on 2026-10-22. Version-sensitive behavior should be checked against the installed versions.
