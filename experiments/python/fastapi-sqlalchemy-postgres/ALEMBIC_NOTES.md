# Alembic quick reference

Changing a SQLAlchemy model does not change the database. Alembic helps bridge that gap: it compares the model metadata with the current database schema and generates a migration. Review that migration, apply it, then check the result.

The commands below use this experiment's `uv` project and `.env` file. Run them from this directory with Postgres available.

## 1. Set up Alembic once

```sh
uv add --dev alembic
uv run --locked alembic init -t async alembic
```

The first command adds Alembic to the project dependencies. The second creates `alembic.ini`, `alembic/env.py`, and `alembic/versions/` using the async template. This setup already exists in this experiment, so do not repeat it for each migration.

## 2. Connect Alembic to the models and database

To generate a migration, Alembic needs to know what the schema should look like and what the database looks like now. In `alembic/env.py`, `target_metadata` points to `Base.metadata` from `models.py`; `database_url()` from `database.py` supplies the connection.

The online path creates an async engine from that URL, so it does not use the placeholder `sqlalchemy.url` in `alembic.ini`. This configuration is already in place for the experiment.

## 3. Generate a migration

After changing the models, generate a candidate revision:

```sh
uv run --locked --env-file .env alembic revision --autogenerate -m "describe the change"
```

Alembic compares the model metadata with the connected database and writes a new file under `alembic/versions/`. This does not apply the change. The message becomes part of the filename; the revision ID is generated separately.

## 4. Review the migration file

Autogenerate proposes a migration; the new version file is where you decide whether that proposal is right. Check that `upgrade()` makes the intended schema change, `downgrade()` reverses it, and `down_revision` points to the preceding version. For an initial migration, `down_revision` is `None`. Edit the generated operations before applying them if they do not match what you intended.

## 5. Apply the migration

```sh
uv run --locked --env-file .env alembic upgrade head
```

Alembic runs the unapplied `upgrade()` operations through the latest revision (`head`) and records the current revision in the database's `alembic_version` table.

## 6. Verify the result

Inspect the database after the upgrade. Its columns and constraints should match the reviewed migration, and `alembic_version.version_num` should match the applied revision. If the migration changes existing rows, check the resulting data as well.

## Example from this experiment

The first migration, `8d8b678efb5f_create_todos.py`, created the `todos` table; its downgrade dropped it. After applying the migration, inspection showed the expected primary key, non-null columns, and `false` default for `completed`. The database recorded revision `8d8b678efb5f`, and the table had zero rows.

For that demonstration only, the local database was reset with `docker compose down --volumes` because it contained tables without matching migration files. That command deleted the experiment's local database contents and is not part of the regular workflow.
