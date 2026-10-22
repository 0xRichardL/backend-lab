.PHONY: sync dev lint format-check db-up db-down db-check

# Install the exact dependency versions already recorded in uv.lock.
sync:
	uv sync --locked

# Serve the app locally and reload it when Python source changes.
dev:
	uv run --locked uvicorn main:app --host 127.0.0.1 --port 8000 --reload

# Keep the small demonstration consistent with its configured Ruff rules.
lint:
	uv run --locked ruff check .

format-check:
	uv run --locked ruff format --check .

# Start and stop the local Postgres service without deleting its named volume.
db-up:
	docker compose up -d --wait

db-down:
	docker compose down

# Load local credentials into the Python process and run a real SELECT 1.
db-check:
	uv run --locked --env-file .env python database.py
