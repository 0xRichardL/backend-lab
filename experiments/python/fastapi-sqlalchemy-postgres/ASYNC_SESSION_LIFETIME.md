# Async session lifetime

## Core model

- The `AsyncEngine` owns the database connection pool and lives across requests. Creating it does not immediately open a connection.
- An `async_sessionmaker` holds the settings for making sessions. Calling it creates a new `AsyncSession`; the factory can be shared, but a session must not be shared across concurrent requests.
- A FastAPI dependency that yields a session provides it to a route for one request. Its cleanup resumes after the route finishes and closes that session.
- A session checks out a connection when database work needs one. A session object is not itself a permanently held connection.

## Two-request trace

Requests A and B use the same engine and session factory. The dependency creates session A and session B separately. Each route uses its own session. After each request's route work, its dependency closes that request's session. The engine remains available for later requests.

## App code in this experiment

`main.py` uses a FastAPI lifespan function: before its `yield`, it creates the engine and stores a session factory on `app.state`; after its `yield`, it disposes the engine at shutdown. `database.py` has a separate `get_session` dependency: it calls that factory inside `async with`, yields a fresh session to one request, then closes that session when FastAPI resumes the dependency. The two `yield` expressions serve different lifetimes.

## Transaction boundary

Closing a session releases its resources and rolls back any uncommitted transaction. It does not commit writes. A write route must explicitly choose when to commit. The create-todo item will show that boundary.

## Relevant failure modes

- Reusing one session for concurrent requests can mix mutable transaction state and is unsupported by SQLAlchemy.
- Letting the engine be created afresh for every request wastes the shared pool and complicates disposal.
- Omitting session cleanup can leave connections checked out longer than needed.

## Sources

- [SQLAlchemy 2.0 asyncio documentation](https://docs.sqlalchemy.org/en/20/orm/extensions/asyncio.html)
- [FastAPI dependencies with `yield`](https://fastapi.tiangolo.com/tutorial/dependencies/dependencies-with-yield/)
- [FastAPI lifespan events](https://fastapi.tiangolo.com/advanced/events/)
