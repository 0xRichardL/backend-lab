# Hertz experiment

## Teaching guide for agents

Use this protocol whenever teaching, reviewing, or modifying this experiment.

### Audience and pace

- Assume the learner is a senior backend engineer who already understands Go, HTTP, routing, middleware, testing, and service lifecycle concepts.
- Focus on Hertz-specific APIs, behavior, defaults, and trade-offs. Do not reteach basic language or backend concepts.
- Teach one cohesive vertical slice per iteration. Pause for the learner once, after giving them a meaningful implementation task; do not turn every API call into a separate checkpoint.
- Do not require predictions or learning-log entries unless they would reveal a genuinely uncertain or surprising framework behavior.

### Teaching sequence

For each vertical slice:

1. State the capability and why it matters in the current experiment.
2. Extract at most five essential Hertz mechanics from the relevant official documentation. Link the exact documentation as an optional reference instead of asking the learner to read it end to end.
3. Show a complete adjacent example that includes its types, containing function, route or middleware registration, and required imports. Never send an orphan code fragment without explaining where it belongs.
4. Give the learner one focused implementation task with explicit behavior and acceptance checks.
5. Review the learner's actual diff and run the narrowest relevant checks.
6. Report findings first and propose exact corrections, but do not edit the learner's code unless they explicitly request implementation.
7. After the implementation passes, give only relevant improvements, advanced options, and production notes.

### Comparisons and sources

- Hertz is the focus. Mention Gin or Echo only when one short analogy materially accelerates understanding.
- Use whichever framework has the closest semantics; do not compare against both by default and do not force a mapping where behavior differs.
- Separate verified Hertz behavior, defaults observed in this experiment, inferences, and production recommendations.
- Prefer current official documentation and source code. Keep links available for depth, but present the useful technique directly in the lesson.

## Research question

How does Hertz structure a small HTTP service, and what are the developer-experience trade-offs of its routing, request context, validation, middleware, testing, shutdown, and `hz` code generation?

Use [LEARNING_PLAN.md](LEARNING_PLAN.md) as the topic syllabus and acceptance checklist. This README controls how those topics are taught.

## Status

Active. Stages 0–3 are complete and Stage 4, testing workflow, is planned.

## Versions

- Hertz: `v0.10.6`
- Module language version: Go `1.20` (the minimum declared by Hertz `v0.10.6`)
- Locally verified toolchain: Go `1.25.6` on macOS/arm64
- Documentation checked: 2026-09-19

## Scenario

Build a small in-memory task API in vertical steps:

- `GET /health`
- `POST /tasks` with JSON binding and validation
- `GET /tasks/:id` with a path parameter
- consistent JSON errors
- request timing middleware
- in-process handler tests
- graceful shutdown behavior

No database, authentication, container, or external service is needed for this research question.

## Setup and commands

Run commands from this directory.

```sh
go mod download
go test ./...
go run .
```

In another terminal:

```sh
curl -i http://127.0.0.1:8888/health
```

Expected response:

```text
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{"status":"ok"}
```

Stop the server with `Ctrl-C`. Hertz `v0.9.6` and later handles `SIGINT`, `SIGHUP`, and `SIGTERM` with graceful shutdown.

## Implementation notes

- `server.Default` provides Hertz's default engine and middleware.
- `server.WithHostPorts` binds the experiment to loopback port `8888`.
- A handler receives Go's `context.Context` plus Hertz's `*app.RequestContext`.
- `newServer` keeps route construction separate from `main`, so tests can exercise the routing engine without opening a network port.
- `ut.PerformRequest` runs a request through the engine in process, similar in purpose to the standard library's `httptest` utilities.
- The create-task request uses `json:"title,required"` for presence during binding and `vd:"len($)>0"` for non-empty validation.
- Server middleware registered with `h.Use` wraps routes registered after it; group middleware applies only to matched routes in that group.
- `RequestContext.Set` and `GetString` propagate the request ID within the request lifecycle.

## Findings

### Facts

- Hertz is an HTTP framework for Go. It uses Netpoll by default and can switch to the Go standard network implementation.
- Routes are registered with methods such as `GET`, `POST`, and `PUT`; route groups and path parameters are built in.
- Binding and validation are exposed through `RequestContext`, including `BindAndValidate`, `BindJSON`, `BindQuery`, and related methods.
- `hz` can generate project scaffolding from Thrift or Protobuf IDL. It is intentionally deferred until the hand-written API is understood.

### Observations

- The smallest server needs one engine, one route, and `Spin()`.
- Separating engine construction makes a route test small and avoids listening on a real port.
- Missing, empty, and malformed title inputs can share a stable public error response while `BindAndValidate` handles their different internal failure paths.
- Middleware around `c.Next(ctx)` runs pre-handler work in registration order and post-handler work in reverse order.
- Aborting an inner middleware skips pending handlers but still returns control to post-handler work in already-running outer middleware.

### Inferences

- Hertz will feel familiar to developers who have used Gin or Echo, but its own `RequestContext`, test utilities, network layer, and generated-code workflow deserve direct study.
- A tiny hand-written API should make it easier to judge what `hz` adds or obscures later.

### Opinions

- Starting without `hz` is the clearest learning path because it exposes the framework's core API before generated structure is introduced.

## Strengths and limitations

Current evidence suggests a concise routing API, explicit binding helpers, composable middleware, and first-party in-process test support. The experiment remains an in-memory API, so it does not yet support conclusions about lifecycle behavior, generated code, production operations, or performance.

## Verdict

Not enough evidence yet. Complete stages 4–6 in the learning plan before deciding where Hertz fits compared with other Go HTTP frameworks.

## References

- [Hertz overview](https://www.cloudwego.io/docs/hertz/overview/)
- [Getting started](https://www.cloudwego.io/docs/hertz/getting-started/)
- [Routing](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/route/)
- [Binding and validation](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/binding-and-validate/)
- [Middleware](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/middleware/)
- [Unit testing](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/unit-test/)
- [Graceful shutdown](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/graceful-shutdown/)
- [`hz` installation and operation](https://www.cloudwego.io/docs/hertz/tutorials/toolkit/install/)
- [Hertz releases](https://github.com/cloudwego/hertz/releases)
