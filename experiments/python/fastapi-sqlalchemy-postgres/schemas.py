"""Public todo shapes used at the API boundary."""

from pydantic import BaseModel, Field


class TodoCreate(BaseModel):
    # A new todo must have a nonempty title; completion starts as false.
    title: str = Field(min_length=1)
    completed: bool = False


class TodoUpdate(BaseModel):
    # Optional fields let PATCH callers send only the values they want to change.
    title: str | None = Field(default=None, min_length=1)
    completed: bool | None = None


class TodoRead(BaseModel):
    # Responses include the server-assigned ID and the stored values.
    id: int
    title: str
    completed: bool
