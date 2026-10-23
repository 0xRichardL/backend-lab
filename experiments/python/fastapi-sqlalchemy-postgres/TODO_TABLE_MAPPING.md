# Todo table mapping

`Todo` in `models.py` is the ORM class for the `todos` table. `Base.metadata`
holds its table definition.

| Python attribute | Postgres column | Rule |
| --- | --- | --- |
| `id: Mapped[int]` | `id SERIAL` | Integer primary key; Postgres generates a value on insert. |
| `title: Mapped[str]` | `title TEXT NOT NULL` | Required stored title. |
| `completed: Mapped[bool]` | `completed BOOLEAN NOT NULL DEFAULT false` | Database default if an insert omits it. |

`Mapped[T]` tells SQLAlchemy and type checkers the Python attribute type.
`mapped_column(...)` supplies column details such as primary key, SQL type, and
default. Without `None` in the `Mapped` type, these columns are non-null.

The database default is separate from Pydantic's `TodoCreate.completed = False`:
Pydantic fills an omitted API field before a route writes the row, while the
database default protects inserts that omit the column altogether. Declaring
the mapping builds metadata in Python; it does not create a table in Postgres.
