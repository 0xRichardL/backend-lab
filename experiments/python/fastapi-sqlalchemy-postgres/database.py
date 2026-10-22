"""Database configuration and a focused connectivity check."""

import asyncio
import os

from sqlalchemy import text
from sqlalchemy.engine import URL
from sqlalchemy.ext.asyncio import create_async_engine


def database_url() -> URL:
    # URL.create handles escaping the password without hand-building a URL string.
    return URL.create(
        drivername="postgresql+asyncpg",
        username=os.environ["POSTGRES_USER"],
        password=os.environ["POSTGRES_PASSWORD"],
        host="127.0.0.1",
        port=int(os.environ.get("POSTGRES_PORT", "5433")),
        database=os.environ["POSTGRES_DB"],
    )


async def check_connection() -> None:
    # Creating an engine does not connect; the first query checks the real path.
    engine = create_async_engine(database_url())
    try:
        async with engine.connect() as connection:
            value = await connection.scalar(text("SELECT 1"))
            print(f"Postgres connectivity: SELECT 1 returned {value}")
    finally:
        await engine.dispose()


if __name__ == "__main__":
    asyncio.run(check_connection())
