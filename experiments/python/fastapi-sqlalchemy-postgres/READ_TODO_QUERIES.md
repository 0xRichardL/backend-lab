# Reading todos

`GET /todos` uses `select(Todo).order_by(Todo.id)` through the request's `AsyncSession`. `await session.scalars(statement)` returns mapped `Todo` objects directly; `.all()` collects them. SQL has no guaranteed row order without `ORDER BY`, so ordering by the primary key makes collection responses stable for this demonstration.

The route's list comprehension converts each mapped row to the public response shape:

```python
return [
    TodoRead(id=todo.id, title=todo.title, completed=todo.completed) for todo in todos
]
```

This is equivalent to starting an empty list, looping over `todos`, appending a `TodoRead` for each row, and returning the list. The query decides row order; the comprehension preserves that order while changing each item's representation.

`GET /todos/{todo_id}` uses `await session.get(Todo, todo_id)`. `get` looks up a mapped object by primary key and returns `None` when it is absent. In this app each request receives a new session, so a lookup that is not already in that session's identity map reads the database. The route will turn `None` into HTTP `404`.

For a found row, this implementation explicitly constructs `TodoRead(id=todo.id, title=todo.title, completed=todo.completed)` and returns that Pydantic object. FastAPI then applies the declared `response_model=TodoRead` to validate, filter, and serialize the response. The route does not hand the mapped `Todo` object directly to FastAPI.

If a handler instead returns a compatible dictionary, FastAPI can validate it against `TodoRead`; fields outside the response model are filtered out. If required response fields are missing or invalid, that is a server-side response validation failure, not a client input `422`. Returning a `Response` object directly bypasses this automatic response-model conversion. The current `-> TodoRead` annotation also tells type checkers that this handler should return a `TodoRead`.

Both routes return `TodoRead` shapes. Reads do not call `commit()` because they do not change rows. Closing the request session ends its read transaction and releases its connection.

See the [SQLAlchemy ORM querying guide](https://docs.sqlalchemy.org/en/20/orm/queryguide/select.html), [session primary-key lookup documentation](https://docs.sqlalchemy.org/en/20/orm/session_api.html#sqlalchemy.orm.Session.get), and [FastAPI response model documentation](https://fastapi.tiangolo.com/tutorial/response-model/).
