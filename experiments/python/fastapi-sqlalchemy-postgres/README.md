# FastAPI–SQLAlchemy–Postgres experiment

## Research question

How does a small FastAPI todo request move through Pydantic validation, an async SQLAlchemy session, and Postgres, and where does the transaction end?

## Status

Active. The health route, Pydantic models, Postgres connection, todo table mapping, initial Alembic migration, request-scoped sessions, and create, list, and read routes are in place. Update and delete are still ahead.

## Run locally

Use Python 3.13 and Docker. This experiment uses Postgres 17; `uv.lock` records the exact FastAPI, SQLAlchemy, Alembic, and other Python dependency versions.

From this directory, run `make sync`. For the first database run, copy `.env.example` to `.env` and replace the placeholder password with a local-only password. Then:

```sh
make db-up
uv run --locked --env-file .env alembic upgrade head
make db-check
make dev
```

`make dev` serves `GET /health` at `http://127.0.0.1:8000/health`, `GET /db-check` for a query through a request-scoped async session, `POST /todos` to create a row, `GET /todos` to list rows by ascending ID, and `GET /todos/{todo_id}` to read one row. Interactive API docs are at `/docs`. Stop the server with `Ctrl+C`. Run `make db-down` to stop Postgres without deleting its data. The available code checks are `make lint` and `make format-check`.

## Documentation

- [Learning guide](./FASTAPI_SQLALCHEMY_POSTGRES_GUIDE.md): scenario, planned steps, scope, and primary sources.
- [Progress log](./FASTAPI_SQLALCHEMY_POSTGRES_PROGRESS.md): completed work, observed evidence, and next step.
- [Todo table mapping](./TODO_TABLE_MAPPING.md): model fields, database columns, and defaults.
- [Alembic notes](./ALEMBIC_NOTES.md): migration workflow and the initial revision example.
- [Async session lifetime](./ASYNC_SESSION_LIFETIME.md): shared engine and factory, per-request session, and cleanup.
- [Create todo flow](./CREATE_TODO_FLOW.md): validated input, insert and commit, refresh, and response.
- [Read todo queries](./READ_TODO_QUERIES.md): ordered collection query, primary-key lookup, and response conversion.

## Current findings

Local checks have confirmed the health response, Pydantic model behavior, Postgres connectivity, and the initial `todos` schema after migration. Two `/db-check` requests returned HTTP 200 with `{"connected": true}` through distinct async sessions. An in-process `POST /todos` against local Postgres returned HTTP 201 and `{"id":1,"title":"Buy milk","completed":false}`; a separate query found the matching stored row. The route commits before returning. `GET /todos` returned IDs `[1,2,3,4]` in ascending order, and `GET /todos/1` returned the matching object. The full CRUD path has not been demonstrated yet, so there is no final verdict.
