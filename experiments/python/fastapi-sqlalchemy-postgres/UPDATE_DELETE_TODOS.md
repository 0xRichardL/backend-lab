# Updating and deleting todos

`PATCH /todos/{todo_id}` reads the mapped `Todo`, applies only fields supplied in `TodoUpdate`, commits, and returns a `TodoRead`. `payload.model_dump(exclude_unset=True)` distinguishes omission from a field explicitly sent by the client. For `{"completed": true}`, it produces `{"completed": true}`; `title` stays unchanged.

An explicit `null` is included by `exclude_unset=True`. Both stored columns are non-null, so the route should reject a supplied `null` instead of assigning it to the mapped object. An empty `{}` patch changes nothing.

`DELETE /todos/{todo_id}` reads the row, calls `await session.delete(todo)` to mark it for deletion, then `await session.commit()` to persist the delete. A successful delete returns HTTP `204 No Content`, so it has no response body. A missing ID is handled as `404`.

See the [FastAPI partial-update guide](https://fastapi.tiangolo.com/tutorial/body-updates/) and [SQLAlchemy async session documentation](https://docs.sqlalchemy.org/en/20/orm/extensions/asyncio.html).
