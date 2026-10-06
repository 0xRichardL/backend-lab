# Hertz experiment

## Research question

How does Hertz structure a small HTTP service, and what are the developer-experience trade-offs of its routing, request context, validation, middleware, testing, lifecycle, and `hz` code generation?

## Status

**Active.** The hand-written HTTP API and live-process shutdown investigation are implemented. Broker-coordinated shutdown, `hz` generation, operational visibility, and the final verdict remain open.

Learning is organized separately in [HERTZ_GUIDE.md](HERTZ_GUIDE.md), with durable state in [HERTZ_PROGRESS.md](HERTZ_PROGRESS.md).

## Versions

- Hertz: `v0.10.6`
- Module language version: Go `1.20`, the minimum declared by Hertz `v0.10.6`
- Locally verified toolchain: Go `1.25.6` on macOS/arm64
- Official documentation checked: 2026-10-06

## Scenario

The experiment implements a small in-memory task API:

| Method | Path | Purpose | Important behavior |
| --- | --- | --- | --- |
| `GET` | `/health` | Liveness response | Does not require a request ID. |
| `GET` | `/tasks/recent` | Static route in the task group | Demonstrates precedence over `/:id`. |
| `GET` | `/tasks/:id` | Parameterized route | Captures exactly one path segment. |
| `POST` | `/tasks` | Bind and validate JSON | Requires `X-Request-ID`; returns a stable public validation error. |
| `GET` | `/grateful/slow` | Shutdown diagnostic | Simulates 10 seconds of in-flight work. |
| `GET` | `/grateful/exceed` | Shutdown diagnostic | Simulates 35 seconds of in-flight work. |

The task API has no database or authentication. The diagnostic routes use `time.Sleep` to isolate HTTP lifecycle behavior; they are not production endpoints or background workers.

## Setup and commands

Run all commands from this directory.

```sh
go mod download
go test ./...
go test -race ./...
go run .
```

The server listens on `127.0.0.1:8888`.

### Live smoke checks

Health:

```sh
curl -i http://127.0.0.1:8888/health
```

Create a valid task:

```sh
curl -i \
  -H 'Content-Type: application/json' \
  -H 'X-Request-ID: smoke-01' \
  -d '{"title":"Learn Hertz"}' \
  http://127.0.0.1:8888/tasks
```

Expected task response:

```text
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
X-Request-ID: smoke-01

{"title":"Learn Hertz"}
```

### Graceful-shutdown reproduction

Build a binary so the signal targets Hertz rather than a `go run` wrapper:

```sh
go build -o /tmp/hertz-lifecycle .
/tmp/hertz-lifecycle
```

In another terminal, start one diagnostic request and send `SIGTERM` about one second later. Restart the binary before testing the other path.

```sh
server_pid=$(pgrep -f '^/tmp/hertz-lifecycle$')
curl -sS -i --max-time 45 -w '\nrequest_duration=%{time_total}s\n' \
  http://127.0.0.1:8888/grateful/slow &
request_pid=$!
sleep 1
kill -TERM "$server_pid"
wait "$request_pid"
```

Replace `/grateful/slow` with `/grateful/exceed` for the over-budget case. While shutdown is in progress, this command checks whether the listener still accepts new connections:

```sh
curl -i --max-time 2 http://127.0.0.1:8888/health
```

## Implementation notes

- `server.Default` creates the engine with Hertz's recovery middleware.
- `server.WithHostPorts` binds the experiment to loopback port `8888`.
- `server.WithExitWaitTime(30*time.Second)` bounds graceful shutdown. The budget begins when Hertz receives the termination signal.
- `newServer` separates engine construction from `main`, allowing route and middleware tests without opening a port.
- Handlers receive Go's `context.Context` and Hertz's pooled `*app.RequestContext`.
- `POST /tasks` uses `json:"title,required"` for field presence and `vd:"len($)>0"` for non-empty validation.
- Global timing middleware wraps all matched routes. Task-group middleware requires and stores `X-Request-ID` in the request context.
- `ut.CreateUtRequestContext` supports direct handler tests; `ut.PerformRequest` exercises the engine, routes, middleware, and handlers in process.
- Lifecycle behavior is verified with a compiled process and real signal because in-process engine tests do not exercise OS signal handling or sockets.

## Findings

### Facts

- Hertz provides grouped, static, parameterized, and wildcard routes. Static routes take precedence over parameterized routes.
- Binding and validation are available through `RequestContext`, including `BindAndValidate`.
- Middleware can run work before and after `c.Next(ctx)` and can abort the remaining handler chain.
- In Hertz `v0.10.6`, `Spin()` handles `SIGINT`, `SIGHUP`, and `SIGTERM` with graceful shutdown.
- `WithExitWaitTime` defines the maximum graceful-shutdown wait; it does not make handler work cancellable or guarantee completion.
- `hz` can generate Hertz projects from Thrift or Protobuf IDL. This experiment has not evaluated its generation and regeneration workflow yet.

### Observations

- The minimum service needs an engine, a registered route, and `Spin()`.
- A separately constructed engine makes route and middleware integration tests small and avoids a network listener.
- Missing, empty, and malformed task titles share the public response `{"error":"invalid request"}` while `BindAndValidate` handles their different internal failure paths.
- Middleware around `c.Next(ctx)` executes pre-handler work in registration order and post-handler work in reverse order.
- Aborting inner middleware skips pending handlers but returns control to post-handler work in already-running outer middleware.
- Direct handler tests isolate binding and response construction; engine tests add routing and middleware without network transport.
- The implemented success and failure paths pass under `go test -race ./...`.

The lifecycle observations below were reproduced on 2026-09-21 with a 30-second shutdown budget and `SIGTERM` sent about one second after each request began:

- `/grateful/slow` completed in `10.003499s` with `200`, `{"status":"ok"}`, and `Connection: close`.
- `/grateful/exceed` ended in `31.003743s` with `curl: (52) Empty reply from server`; its handler-exit log did not appear before process exit.
- A new `/health` connection attempted during shutdown failed with `curl: (7)`, showing that the listener was no longer accepting connections.
- A separate `/grateful/exceed` run returned `200` when `SIGTERM` arrived about six seconds into the handler because the remaining work fit inside the 30-second shutdown budget.

### Inferences

- Hertz's core HTTP API is concise and familiar, but its `RequestContext`, test utilities, Netpoll default, and generated-code workflow require framework-specific evaluation.
- The shutdown deadline is measured from signal receipt, not request start; deployment budgets must account for the maximum remaining work at that moment.
- An empty HTTP reply does not prove that a real business side effect failed. Persistence, acknowledgements, retries, and idempotency require a separate experiment.

### Opinions

- Starting with a hand-written API is the clearest way to understand Hertz before evaluating `hz`.
- Keeping one direct-handler layer, one in-process engine layer, and a small number of live-process checks gives useful evidence without unnecessary test infrastructure.

## Strengths and limitations

For this small service, Hertz provides concise routing, explicit binding helpers, composable middleware, first-party in-process test utilities, and bounded graceful shutdown.

The evidence remains intentionally narrow. The API is in memory, the lifecycle probes simulate work with sleeps, and the experiment has not yet evaluated persistence, broker consumers, `hz` regeneration, observability integrations, deployment configuration, or performance. The current results do not support cross-framework ranking.

## Verdict

**Provisional.** Hertz is straightforward for a small Go HTTP API and exposes useful framework-level testing and lifecycle controls. A recommendation for production use still depends on the unfinished broker-lifecycle, `hz`, and observability investigations. Performance claims are explicitly out of scope until measured with a fair workload.

## References

- [Hertz overview](https://www.cloudwego.io/docs/hertz/overview/)
- [Getting started](https://www.cloudwego.io/docs/hertz/getting-started/)
- [Routing](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/route/)
- [Binding and validation](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/binding-and-validate/)
- [Middleware](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/middleware/)
- [Unit testing](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/unit-test/)
- [Graceful shutdown](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/graceful-shutdown/)
- [`hz` code generation](https://www.cloudwego.io/docs/hertz/tutorials/toolkit/)
- [Hertz releases](https://github.com/cloudwego/hertz/releases)
