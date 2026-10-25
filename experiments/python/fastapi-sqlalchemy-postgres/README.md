# FastAPI–SQLAlchemy–Postgres experiment

## Research question

How does a small FastAPI todo request move through Pydantic validation, an async SQLAlchemy session, and Postgres, and where does the transaction end?

## Status

Active. The health route, Pydantic models, Postgres connection, todo table mapping, and initial Alembic migration are in place. Request-scoped sessions and todo CRUD are still ahead.

## Run locally

Use Python 3.13 and Docker. This experiment uses Postgres 17; `uv.lock` records the exact FastAPI, SQLAlchemy, Alembic, and other Python dependency versions.

From this directory, run `make sync`. For the first database run, copy `.env.example` to `.env` and replace the placeholder password with a local-only password. Then:

```sh
make db-up
uv run --locked --env-file .env alembic upgrade head
make db-check
make dev
```

`make dev` serves `GET /health` at `http://127.0.0.1:8000/health`; interactive API docs are at `/docs`. Stop the server with `Ctrl+C`. Run `make db-down` to stop Postgres without deleting its data. The available code checks are `make lint` and `make format-check`.

## Documentation

- [Learning guide](./FASTAPI_SQLALCHEMY_POSTGRES_GUIDE.md): scenario, planned steps, scope, and primary sources.
- [Progress log](./FASTAPI_SQLALCHEMY_POSTGRES_PROGRESS.md): completed work, observed evidence, and next step.
- [Todo table mapping](./TODO_TABLE_MAPPING.md): model fields, database columns, and defaults.
- [Alembic notes](./ALEMBIC_NOTES.md): migration workflow and the initial revision example.

## Current findings

Local checks have confirmed the health response, Pydantic model behavior, Postgres connectivity, and the initial `todos` schema after migration. The request-to-database path and transaction boundary have not been demonstrated yet, so there is no final verdict.
