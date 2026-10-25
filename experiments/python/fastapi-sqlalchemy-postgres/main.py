"""FastAPI app and its shared database resources."""

from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager
from typing import Annotated

from fastapi import Depends, FastAPI
from sqlalchemy import text
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

from database import database_url, get_session


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None]:
    # One engine and session factory serve all requests for this app lifetime.
    engine = create_async_engine(database_url())
    app.state.session_factory = async_sessionmaker(engine, expire_on_commit=False)
    try:
        yield
    finally:
        # Dispose of pooled connections when the app shuts down.
        await engine.dispose()


# The ASGI server imports this application object from ``main:app``.
app = FastAPI(title="Todo API Lab", lifespan=lifespan)


# FastAPI registers this function for GET requests to /health.
@app.get("/health")
async def health() -> dict[str, str]:
    # FastAPI serializes the returned dict as a JSON response with HTTP 200.
    return {"status": "ok"}


@app.get("/db-check")
async def db_check(
    session: Annotated[AsyncSession, Depends(get_session)],
) -> dict[str, bool]:
    # FastAPI gets a fresh session from get_session for this request.
    value = await session.scalar(text("SELECT 1"))
    return {"connected": value == 1}
