# Initial Alembic migration notes

This note tracks LG-06's demonstrated steps. Commands run from this experiment directory unless stated otherwise.

## Starting state

`models.py` defines `Todo` with `id`, `title`, and `completed`. Its `Base.metadata` describes the intended table but does not create it. Alembic will compare that metadata with the live database to generate a candidate revision. Review the revision before applying it.

The prior demonstration had left `todos` and `alembic_version` in local Postgres without matching migration files. At the learner's request, the experiment database was reset.

## Commands and observed results

1. `docker compose down --volumes` removed this experiment's Postgres container and `fastapi-sqlalchemy-postgres_postgres_data` volume. This deleted the prior local database contents.
2. `make db-up` created a new volume and started a healthy Postgres 17 container.
3. A read-only query of `information_schema.tables` for the `public` schema returned zero tables.
4. `uv add --dev alembic` added Alembic 1.20.1 to `pyproject.toml` and `uv.lock`.
5. `uv run --locked alembic init -t async alembic` generated `alembic.ini`, `alembic/env.py`, `alembic/script.py.mako`, `alembic/README`, and an empty `alembic/versions` directory.
6. `cat alembic.ini alembic/env.py alembic/README alembic/script.py.mako` showed the unedited starter files. The generated `env.py` has `target_metadata = None`, and `alembic.ini` has a placeholder `sqlalchemy.url`.
7. Edited `alembic/env.py` to use `Base.metadata` for autogenerate and the app's `database_url()` for both offline and online migrations. The online path creates an async engine directly, so the placeholder URL in `alembic.ini` is not used. `cat alembic/env.py` showed the resulting file.
8. `uv run --locked --env-file .env alembic revision --autogenerate -m "create todos"` detected one added table and wrote `alembic/versions/8d8b678efb5f_create_todos.py`. Inspection showed an initial revision (`down_revision = None`) whose `upgrade()` creates `todos` with non-null `id`, `title`, and `completed`, a primary key on `id`, and a server default of `false` for `completed`. `downgrade()` drops the table. This is a generated candidate; it has not been applied.
9. `.venv/bin/ruff format alembic/versions/8d8b678efb5f_create_todos.py` reformatted the generated candidate. Inspection confirmed the operations remained the same.
10. `uv run --locked --env-file .env alembic upgrade head` applied revision `8d8b678efb5f` to the fresh local Postgres database.
11. A read-only `psql` inspection showed `todos(id integer primary key, title text not null, completed boolean not null default false)`, `alembic_version.version_num = 8d8b678efb5f`, and zero todo rows.

## Next step

Document the repeatable migration command in the experiment README, run the relevant local checks, and commit LG-06's revision and documentation.
