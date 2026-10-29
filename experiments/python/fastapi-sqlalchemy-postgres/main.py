"""FastAPI app and its shared database resources."""

from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager
from typing import Annotated

from fastapi import Depends, FastAPI, HTTPException
from sqlalchemy import delete, select, text
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

from database import database_url, get_session
from models import Todo
from schemas import TodoCreate, TodoRead, TodoUpdate


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


@app.post("/todos", status_code=201, response_model=TodoRead)
async def create_todo(
    payload: TodoCreate,
    session: Annotated[AsyncSession, Depends(get_session)],
) -> TodoRead:
    # FastAPI validates the request body before entering this function.
    todo = Todo(title=payload.title, completed=payload.completed)
    session.add(todo)
    # Commit flushes the INSERT and ends the write transaction.
    await session.commit()
    # Reload the stored row, including its database-generated ID.
    await session.refresh(todo)
    return TodoRead(id=todo.id, title=todo.title, completed=todo.completed)


@app.get("/todos", response_model=list[TodoRead])
async def list_todos(
    session: Annotated[AsyncSession, Depends(get_session)],
) -> list[TodoRead]:
    # ORDER BY makes the collection response stable across requests.
    todos = (await session.scalars(select(Todo).order_by(Todo.id))).all()
    return [
        TodoRead(id=todo.id, title=todo.title, completed=todo.completed)
        for todo in todos
    ]


@app.get("/todos/{todo_id}", response_model=TodoRead)
async def read_todo(
    todo_id: int,
    session: Annotated[AsyncSession, Depends(get_session)],
) -> TodoRead:
    # get() uses the mapped primary key and returns None for a missing row.
    todo = await session.get(Todo, todo_id)
    if todo is None:
        raise HTTPException(status_code=404, detail="Todo not found")
    return TodoRead(id=todo.id, title=todo.title, completed=todo.completed)


@app.patch("/todos/{todo_id}", response_model=TodoRead)
async def update_todo(
    todo_id: int,
    payload: TodoUpdate,
    session: Annotated[AsyncSession, Depends(get_session)],
) -> TodoRead:
    # Exclude omitted fields, but keep explicit nulls so they can be rejected.
    changes = payload.model_dump(exclude_unset=True)
    if any(value is None for value in changes.values()):
        raise HTTPException(status_code=422, detail="Todo fields cannot be null")

    todo = await session.get(Todo, todo_id)
    if todo is None:
        raise HTTPException(status_code=404, detail="Todo not found")

    # TodoUpdate exposes only the two mutable mapped fields.
    for field, value in changes.items():
        setattr(todo, field, value)
    await session.commit()
    await session.refresh(todo)
    return TodoRead(id=todo.id, title=todo.title, completed=todo.completed)


@app.delete("/todos/{todo_id}", status_code=204)
async def delete_todo(
    todo_id: int,
    session: Annotated[AsyncSession, Depends(get_session)],
) -> None:
    # RETURNING identifies a deleted row without a preceding SELECT.
    result = await session.execute(
        delete(Todo).where(Todo.id == todo_id).returning(Todo.id)
    )
    if result.scalar_one_or_none() is None:
        raise HTTPException(status_code=404, detail="Todo not found")
    # The DELETE ran in this session's transaction; commit persists it.
    await session.commit()
