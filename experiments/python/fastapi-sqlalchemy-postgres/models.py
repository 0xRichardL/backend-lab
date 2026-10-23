"""SQLAlchemy mapping for the todo table."""

from sqlalchemy import Text, text
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column


class Base(DeclarativeBase):
    # All mapped tables inherit this shared metadata.
    pass


class Todo(Base):
    __tablename__ = "todos"

    # Postgres generates the integer primary key when a row is inserted.
    id: Mapped[int] = mapped_column(primary_key=True)
    # Mapped[str] makes this column non-null; Text keeps title unbounded.
    title: Mapped[str] = mapped_column(Text)
    # The database default applies even if an INSERT omits this column.
    completed: Mapped[bool] = mapped_column(server_default=text("false"))
