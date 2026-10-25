# Creating a todo

`POST /todos` crosses three representations: validated `TodoCreate` input, a SQLAlchemy `Todo` row, and a `TodoRead` response. FastAPI validates the body before calling the route.

```python
todo = Todo(title=payload.title, completed=payload.completed)
session.add(todo)  # Pending in this session; no durable row yet.
await session.commit()  # Flushes the INSERT and commits its transaction.
await session.refresh(todo)  # Reads database-generated fields and stored values.
return TodoRead(id=todo.id, title=todo.title, completed=todo.completed)
```

The route declares `status_code=201` and uses the request's `AsyncSession` from `Depends(get_session)`. The `id` is assigned by Postgres during insertion. `completed` is sent explicitly from the validated input; the table's server default matters for inserts that omit the column.

If validation fails, FastAPI returns `422` before the route runs. If the database write fails, the request does not return a successful `201`; the session must not be reused in its failed transaction state. Closing the request session releases its resources. A later item covers failed writes and rollback in detail.

See the [FastAPI response status documentation](https://fastapi.tiangolo.com/tutorial/response-status-code/) and [SQLAlchemy async session documentation](https://docs.sqlalchemy.org/en/20/orm/extensions/asyncio.html).
