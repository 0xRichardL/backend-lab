"""First FastAPI route for the request-to-response demonstration."""

from fastapi import FastAPI

# The ASGI server imports this application object from ``main:app``.
app = FastAPI(title="Todo API Lab")


# FastAPI registers this function for GET requests to /health.
@app.get("/health")
async def health() -> dict[str, str]:
    # FastAPI serializes the returned dict as a JSON response with HTTP 200.
    return {"status": "ok"}
