# FastAPI–SQLAlchemy–Postgres experiment

## Research question

How does a small FastAPI todo request move through Pydantic validation, an async SQLAlchemy session, and Postgres, and where does the transaction end?

## Status

Active: the uv project and health endpoint are implemented. Follow the [learning guide](./FASTAPI_SQLALCHEMY_POSTGRES_GUIDE.md) for the planned demonstration and its [progress file](./FASTAPI_SQLALCHEMY_POSTGRES_PROGRESS.md) for the current checkpoint.

## Versions and setup

- Python: 3.13 (selected by `.python-version`)
- uv: 0.11.15
- Exact dependency versions: `uv.lock`

From this directory:

```sh
make sync
make dev
```

`make dev` serves the app at `http://127.0.0.1:8000` and reloads it when Python source changes. Stop it with `Ctrl+C`. The other current shortcuts are `make lint` and `make format-check`.

For the Postgres demonstration, start Docker Desktop, copy `.env.example` to
`.env`, and replace the placeholder password with a local-only password. Then:

```sh
make db-up
make db-check
make db-down
```

Compose binds Postgres to `127.0.0.1:5433` and keeps its data in a named volume.
`make db-down` stops the container and preserves that volume. The ignored `.env`
supplies credentials to Compose and to `uv run --env-file .env`; no password is
stored in tracked files. `make db-check` runs `SELECT 1` through SQLAlchemy's
async engine and `asyncpg`. It reported `Postgres connectivity: SELECT 1 returned 1`.

The equivalent uv commands for the app are:

```sh
uv sync --locked
uv run --locked python -c 'import fastapi, uvicorn; print(f"FastAPI {fastapi.__version__}; Uvicorn {uvicorn.__version__}")'
uv run --locked uvicorn main:app --host 127.0.0.1 --port 8000 --reload
```

`uv add` records application dependencies under `[project].dependencies` and development tools under `[dependency-groups].dev` in `pyproject.toml`; it updates `uv.lock`. The `.venv` environment is local and ignored by Git.

With the server running, `GET http://127.0.0.1:8000/health` returns `{"status":"ok"}`. FastAPI serves interactive API documentation at `http://127.0.0.1:8000/docs` and its OpenAPI document at `http://127.0.0.1:8000/openapi.json`.

## Scenario, implementation notes, findings, and verdict

The todo API and its verification will be recorded here as the guide progresses. With Uvicorn 0.54.0 and FastAPI 0.159.1, a local `GET /health` returned HTTP 200 and `{"status":"ok"}`. `HEAD /docs` returned HTTP 200, and `/openapi.json` listed `GET /health`. The local server was stopped after these checks. The API contract is in `schemas.py`; direct Pydantic checks demonstrated create defaults, partial update fields, empty-title rejection, and read-model output filtering. The local Postgres container passed its health check, and SQLAlchemy 2.0.55 with asyncpg 0.32.0 returned `1` for `SELECT 1`.

## References

- [uv projects](https://docs.astral.sh/uv/guides/projects/)
- [uv project files](https://docs.astral.sh/uv/concepts/projects/layout/)
- [FastAPI first steps](https://fastapi.tiangolo.com/tutorial/first-steps/)
- [Postgres official image](https://hub.docker.com/_/postgres)
- [Docker Compose environment variables](https://docs.docker.com/compose/how-tos/environment-variables/variable-interpolation/)
- [SQLAlchemy asyncpg dialect](https://docs.sqlalchemy.org/en/20/dialects/postgresql.html#module-sqlalchemy.dialects.postgresql.asyncpg)
