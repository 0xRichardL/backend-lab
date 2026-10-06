# Learning Progress: Hertz

> Guide: [HERTZ_GUIDE.md](./HERTZ_GUIDE.md)
> Last updated: 2026-10-06 09:11 +07

## Status

| Field | Value |
| --- | --- |
| Overall status | in progress |
| Accumulated reported active time | 0 minutes; historical active-learning time was not reported |
| Current unit | LG-06 — Coordinated shutdown |
| Next action | Design and run the RabbitMQ half of LG-06: stop deliveries on `SIGTERM`, drain bounded in-flight work, ack completed work, and prove unfinished work is redelivered. |

## Unit checklist

Use only `not started`, `in progress`, `completed`, `blocked`, or `deferred`. Check a box only for `completed`.

- [x] LG-01 — Engine and request path — `completed`
- [x] LG-02 — Routing contract — `completed`
- [x] LG-03 — Input and public errors — `completed`
- [x] LG-04 — Middleware and request context — `completed`
- [x] LG-05 — Testing boundaries — `completed`
- [ ] LG-06 — Coordinated shutdown — `in progress`
- [ ] LG-07 — Hand-written code versus `hz` — `not started`
- [ ] LG-08 — Operational visibility — `not started`
- [ ] LG-09 — Hertz evaluation — `not started`

## Current checkpoint

### Last demonstrated evidence

- LG-01: The working `/health` route and live request establish the path from route registration through the Hertz engine and handler to a JSON response.
- LG-02: Engine tests prove static-route precedence, path-parameter capture, single-segment matching, missing group-root behavior, and the default wrong-method result.
- LG-03: Direct handler tests cover valid, malformed, missing-title, and empty-title JSON, including stable status, body, and content type.
- LG-04: Tests prove nested `c.Next(ctx)` ordering, abort behavior, request-ID propagation, and response-time middleware behavior.
- LG-05: Behavior is separated between direct handler tests, in-process engine tests, and a live-network check; `go test -race ./...` passed during the completed stage.
- LG-06: The HTTP half is demonstrated. With a 30-second budget and `SIGTERM` about one second after request start, a 10-second handler returned `200` with `Connection: close`; a 35-second handler ended near the shutdown deadline with `curl: (52) Empty reply from server`; a new connection during shutdown failed.

### Recent misconceptions or fragile knowledge

- A long-running HTTP handler was initially described as a worker. The corrected distinction is important: Hertz drains HTTP connections, but a broker consumer needs separate stop, drain, acknowledgement, and client-close coordination.
- The shutdown budget begins when `SIGTERM` is received, not when an in-flight request began.

### Incomplete knowledge or work

- LG-06 still lacks the RabbitMQ exercise: bounded prefetch, manual acknowledgement after successful processing, consumer cancellation to stop new deliveries, bounded draining, and evidence that unfinished work is redelivered.
- Shutdown ownership and ordering between Hertz handlers, broker consumers, and shared broker publishers have not been applied in code.

### Open questions

- What component should own the broker consumer lifecycle when it runs beside Hertz?
- How should shutdown ordering prevent a shared publisher or connection from closing while HTTP or worker code still uses it?

### Retrieval prompts

- When does `WithExitWaitTime` begin counting, and what does it bound?
- State the safe worker shutdown order from stopping deliveries through closing the broker connection.
- Why can an unacknowledged redelivery duplicate a business side effect?

### Exact re-entry prompt

> Design the smallest RabbitMQ shutdown experiment that proves four outcomes: no new deliveries after shutdown starts, completed work is acknowledged, unfinished work is redelivered, and the application exits within its shutdown budget.

## Adaptations

- 2026-10-06: Migrated historical evidence from the former syllabus. Its first five completed slices map to LG-01–LG-05. Its HTTP lifecycle slice supplies only part of LG-06, so LG-06 remains in progress until the guide's RabbitMQ mastery evidence is demonstrated.
- 2026-10-06: Session allocation was 0 hours, so this session synchronized durable state only and performed no teaching.

## Session log

### Session 1 — 2026-10-06

- Planned time: 0 minutes
- Actual active time: 0 minutes; state synchronization only
- Units worked: LG-01–LG-06 historical evidence migration
- Learner evidence: Imported the completed legacy acceptance evidence; LG-01–LG-05 satisfy their guide evidence, while LG-06 satisfies only its HTTP shutdown portion.
- Misconceptions or uncertainty: Preserve the distinction between an in-flight HTTP handler and a broker worker; broker lifecycle coordination is not yet demonstrated.
- Unresolved work: Complete the RabbitMQ stop, drain, ack/redelivery, and close experiment for LG-06.
- Research notes: Synchronization used the former syllabus and the recorded observations in `README.md`; no new unit research was needed for this zero-time session.
- Next prompt: Design the smallest RabbitMQ shutdown experiment that proves four outcomes: no new deliveries after shutdown starts, completed work is acknowledged, unfinished work is redelivered, and the application exits within its shutdown budget.
