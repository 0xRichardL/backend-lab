# Hertz learning syllabus

The [README teaching guide](README.md#teaching-guide-for-agents) controls lesson pacing and interaction. This file tracks topics, exercises, and acceptance criteria. Official documentation is linked for optional depth, not assigned as prerequisite reading.

## Stage 0 — Baseline: run and inspect — complete

**Question:** What is the minimum Hertz server, and what do `server.Hertz`, `Engine`, `context.Context`, and `RequestContext` each do?

**Exercise:** Run the health endpoint and its in-process test, then inspect startup and graceful interrupt behavior.

**Acceptance:** The learner can explain the request path from route registration to JSON response and reproduce the test and live request.

**Reference:** [Overview](https://www.cloudwego.io/docs/hertz/overview/) and [getting started](https://www.cloudwego.io/docs/hertz/getting-started/).

## Stage 1 — Routing and request input — complete

**Question:** How does Hertz match grouped, static, and parameterized routes?

**Exercise:** Implement the `/tasks` group with static and `/:id` routes, then cover valid, missing, extra-segment, and wrong-method requests.

**Acceptance:** Tests prove parameter capture, static-route precedence, single-segment matching, and Hertz's default missing-route and wrong-method status codes.

**Reference:** [Routing](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/route/).

## Stage 2 — JSON binding, validation, and errors — complete

**Question:** What does Hertz bind automatically, and what error information is safe and useful to return?

**Exercise:** Implement `POST /tasks` with JSON presence validation, non-empty validation, and a stable public error response.

**Acceptance:** Tests cover valid JSON, malformed JSON, missing title, empty title, status codes, response bodies, and content type.

**Reference:** [Binding and validation](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/binding-and-validate/) and [error handling](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/error-handle/).

## Stage 3 — Middleware and request context — active

**Question:** How do pre-handler work, post-handler work, aborts, and request-scoped values behave?

**Exercise:** Add request timing middleware and propagate a request ID through `RequestContext`; test nested `c.Next(ctx)` ordering and an early JSON abort response.

**Acceptance:** Tests assert middleware ordering, expose the request ID in the response, and prove an aborted request does not execute its route handler.

**Reference:** [Middleware](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/middleware/).

## Stage 4 — Testing workflow — planned

**Question:** Which behavior belongs in direct handler tests, engine tests, and live network tests?

**Exercise:** Organize the task API's acceptance checks as focused table-driven tests and retain one manual live smoke check.

**Acceptance:** Current success and failure paths run without an open port or external service, and `go test -race ./...` passes.

**Reference:** [Unit testing](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/unit-test/).

## Stage 5 — Lifecycle and graceful shutdown — planned

**Question:** What happens to an in-flight request when the process receives a termination signal?

**Exercise:** Use a temporary slow endpoint to compare request completion against two `server.WithExitWaitTime` values under `SIGTERM`.

**Acceptance:** Record the exact command, timeout, request duration, signal, and observed result; remove the temporary endpoint afterward unless it remains useful evidence.

**Reference:** [Graceful shutdown](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/graceful-shutdown/).

## Stage 6 — Hand-written code versus `hz` — planned

**Question:** When does IDL-driven generation help, and which code is generated versus learner-owned?

**Exercise:** Generate the same task API from one small Thrift or Protobuf definition in a temporary directory and compare it with this hand-written implementation.

**Acceptance:** Identify generated files, explain safe regeneration, compare dependencies and testability, and state one concrete case where `hz` helps or hurts the workflow.

**Reference:** [`hz` installation](https://www.cloudwego.io/docs/hertz/tutorials/toolkit/install/) and [basic usage](https://www.cloudwego.io/docs/hertz/tutorials/toolkit/usage/).

## Stage 7 — Optional focused investigation

Choose one explicit research question involving the Hertz client, HTTP/2, observability, an alternative network layer, or a measured framework comparison. Define a fair scenario and acceptance checks before adding dependencies or infrastructure.

Do not make performance claims from marketing material or a single local run. Record the workload, tool, environment, warm-up, sample count, and limitations when performance is the research question.

## Optional evidence notes

Record only non-obvious findings that will be useful later, such as a surprising default, a meaningful trade-off, or a framework-comparison insight. Tests remain the primary executable evidence; do not duplicate obvious test results in prose.
