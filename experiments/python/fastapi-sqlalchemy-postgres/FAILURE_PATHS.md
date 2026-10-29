# Todo failure paths

FastAPI validates request bodies and path parameters before calling a route. Invalid create or update input produces HTTP `422` with validation details; no route write occurs. A valid integer ID with no matching row reaches the route. Read and PATCH use `session.get(Todo, id)` and raise `HTTPException(404)` when it returns `None`; DELETE uses `RETURNING` to detect the same missing-row case.

A database write can fail later, such as when a unique or non-null constraint is violated. SQLAlchemy rolls back the failed flush at the database level, but the session remains inactive until `await session.rollback()` if that same session will be reused. In this app, each request's session is closed after the request; closing discards the failed session and releases its resources. Internal database details should not be returned to API clients as a successful response.

The demonstration checked `422` and `404` responses through the API. A separate in-memory SQLite session check forced a constraint failure and showed the session's failed state and unchanged previously committed data after rollback. This does not add a public route whose purpose is to fail.

The DELETE route uses `delete(Todo).where(Todo.id == todo_id).returning(Todo.id)` so it can identify a matching row without a preceding `SELECT`. `result.scalar_one_or_none()` is the deleted ID when a row matched and `None` when none matched. The route raises `HTTPException(404)` for `None`; a matching delete needs `await session.commit()` to persist. This route has not yet been exercised against live Postgres in the current environment.

See the [FastAPI error-handling guide](https://fastapi.tiangolo.com/tutorial/handling-errors/) and [SQLAlchemy session rollback guidance](https://docs.sqlalchemy.org/en/20/orm/session_basics.html#rolling-back).
