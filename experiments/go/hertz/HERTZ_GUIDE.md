# Learning Guide: Hertz

> Created: 2026-10-06
> Guide version: 1
> Target slug: hertz

## Learning contract

| Field | Value |
| --- | --- |
| Desired outcome | Build, test, operate, and critically evaluate a small Hertz service, including one broker-backed worker lifecycle and a fair comparison of hand-written code with `hz` generation. |
| Current level | Senior backend engineer who already understands Go, HTTP, routing, middleware, testing, and service lifecycle concepts; the existing task API provides working evidence for the Hertz fundamentals and HTTP shutdown behavior. |
| Success criteria | Explain Hertz's request path and framework-specific contexts; implement and test routing, binding, validation, errors, and middleware; reproduce bounded HTTP and broker-worker shutdown; safely regenerate an IDL-driven service; expose useful operational signals; and write a versioned verdict that separates fact, observation, inference, and opinion. |
| Total active-learning budget | 18 hours |
| Constraints | Use the existing self-contained task API, Hertz `v0.10.6`, Go `1.20` module compatibility, and primary sources first. Add only one local broker for the lifecycle unit. Make no performance claims without a separate measured investigation. |

## Priority and scope

### Must learn

- Engine construction, request flow, routing, and Hertz's two context types.
- Binding, validation, stable public errors, and middleware control flow.
- Direct-handler, in-process engine, and live-process test boundaries.
- Coordinated `SIGTERM` handling for HTTP requests and one broker consumer.
- `hz` ownership boundaries, regeneration behavior, and the trade-off against hand-written code.

### Supporting knowledge

- Request IDs, structured operational signals, basic tracing, and shutdown logs.
- Manual broker acknowledgement, bounded in-flight work, redelivery, and idempotency.
- Evidence classification and fair framework evaluation.

### Deliberately deferred

- Databases, authentication, authorization, service discovery, TLS, HTTP/2, WebSockets, and deployment infrastructure: useful, but not required to evaluate Hertz's core development model.
- Kafka implementation: RabbitMQ is the single concrete broker exercise; Kafka offset management is a separate experiment with materially different semantics.
- Custom network or protocol layers: defer until a real compatibility or performance requirement exists.
- Benchmarks and cross-framework rankings: defer until a fair workload, environment, warm-up, and sample plan are defined.

## Roadmap

| ID | Unit | Observable outcome | Key ideas | Practice and mastery evidence | Estimate |
| --- | --- | --- | --- | --- | ---: |
| LG-01 | Engine and request path | Trace one request from route registration through handler execution to the encoded response. | `server.Hertz`, embedded `Engine`, `Spin`, `context.Context`, `RequestContext`, default middleware, Netpoll as a configurable implementation detail. | Retrieve the request path from memory, inspect the minimal server, run `/health`, and annotate each Hertz-owned step; mastery is a correct explanation tied to code and one live response. | 1h 15m |
| LG-02 | Routing contract | Implement and predict Hertz's behavior for grouped, static, parameterized, missing, and wrong-method routes. | Route groups, precedence, single-segment parameters, trailing slash behavior, 404/405 configuration. | Extend the task routes and table-test success plus boundary cases; mastery is passing engine tests and an explanation of each observed default. | 1h 30m |
| LG-03 | Input and public errors | Accept valid JSON while rejecting malformed, missing, and invalid fields with a stable external contract. | `BindAndValidate`, binding tags, validator tags, content type, internal error detail versus public error shape. | Implement `POST /tasks`, retrieve tag behavior without notes, and test success/failure cases; mastery is deterministic status, JSON, and content-type evidence. | 2h |
| LG-04 | Middleware and request context | Compose pre-handler, post-handler, abort, and request-scoped behavior without leaking context lifetime. | `Use`, group middleware, `c.Next`, abort semantics, `Set`/`Get`, standard context versus pooled `RequestContext`. | Add timing and request-ID middleware, then test nested order and abort behavior; mastery is a sequence assertion plus a concise context-lifetime rule. | 2h |
| LG-05 | Testing boundaries | Choose the smallest test layer that proves each behavior and retain one real-network check. | Direct handler tests, `ut.PerformRequest`, table tests, race detector, live process tests, reproducible commands. | Reorganize acceptance checks by layer, run `go test -race ./...`, and justify one behavior assigned to each layer; mastery is a green focused suite without an unnecessary open port. | 2h |
| LG-06 | Coordinated shutdown | Demonstrate that `SIGTERM` stops new work, drains bounded in-flight work, and safely returns unfinished broker work. | `WithExitWaitTime`, shutdown budget origin, hooks, stop-then-drain ordering, RabbitMQ manual ack/cancel/close, redelivery, idempotency. | Reproduce the HTTP timeout probe, then run one RabbitMQ worker with bounded prefetch and manual ack; signal it during processing. Mastery is evidence that completed work is acked, unfinished work is redelivered, new work stops, and the hard-kill budget exceeds the application drain budget. | 3h |
| LG-07 | Hand-written code versus `hz` | Generate, change, and regenerate the task API while identifying generated and learner-owned files. | IDL annotations, `hz new` versus `hz update`, model/router regeneration, handler preservation, toolchain and dependency cost. | Generate the same API from one small Thrift or Protobuf IDL in an isolated directory, add one field or route, regenerate, and diff; mastery is a file-ownership map and one evidence-backed help/hurt judgment. | 3h |
| LG-08 | Operational visibility | Expose enough signals to diagnose request failures, latency, and shutdown behavior without building a full platform. | Correlation IDs, structured logs, request duration/error signals, Hertz tracer lifecycle, low-cardinality metrics, broker in-flight/redelivery signals. | Instrument the task API and shutdown path, force one validation error and one interrupted worker, then correlate logs/signals; mastery is a short diagnostic walkthrough from symptom to cause. | 2h |
| LG-09 | Hertz evaluation | Produce a scoped recommendation for when Hertz fits and what remains uncertain. | Framework versus language trade-offs, facts/observations/inferences/opinions, version scope, fair comparison, evidence limits. | Retrieve the main mechanics, review all experiment evidence, and write a concise verdict with strengths, limitations, appropriate use cases, and open questions; mastery is that every material claim cites documentation or a reproducible result. | 1h 15m |

**Planned total:** 18 hours

## Milestones

| Milestone | Units | Cumulative time | Evidence |
| --- | --- | ---: | --- |
| Request contract understood | LG-01–LG-03 | 4h 45m | A working task API with route and input behavior explained and tested. |
| Framework composition is testable | LG-04–LG-05 | 8h 45m | Middleware semantics and test-layer choices are proven without unnecessary infrastructure. |
| Production lifecycle demonstrated | LG-06 | 11h 45m | HTTP and RabbitMQ work respond correctly to `SIGTERM`, including unfinished-work redelivery. |
| Tooling and operations assessed | LG-07–LG-08 | 16h 45m | Regeneration ownership is mapped and failures can be diagnosed from emitted signals. |
| Evidence-based verdict complete | LG-09 | 18h | A versioned, contextual Hertz recommendation with explicit uncertainty. |

## Compact reference

### Request path

- Build routes on `server.Hertz`; its embedded engine matches the request, runs middleware and the handler, then writes the response.
- Keep Go's `context.Context` for cancellation and cross-component propagation; treat Hertz's pooled `RequestContext` as request-lifetime state.

### Middleware

- Code before `c.Next(ctx)` runs on the way in; code after it runs on the way out.
- Abort prevents pending handlers but does not erase post-handler work in middleware already on the stack.

### Test layers

- Direct handler: binding, validation, and response construction.
- In-process engine: routing plus middleware integration.
- Live process: sockets, signals, shutdown, and external dependencies.

### Shutdown

- Stop accepting HTTP requests and broker deliveries first; then wait for in-flight work within a smaller application budget than the platform's hard-kill budget.
- Ack only successful broker work. Expect redelivery after failure or timeout, and make side effects idempotent.
- Closing a resource is not equivalent to proving that all work using it has stopped; coordinate ownership and ordering explicitly.

### `hz` ownership

- IDL-generated models and routers are regeneration inputs/outputs; business handlers are learner-owned only where the generated layout and update behavior preserve them.
- Judge generation by repeatable change cost, diff clarity, dependency/tooling overhead, and testability—not by initial line count.

### Evaluation

- Fact: documented or source-backed behavior.
- Observation: reproduced in this experiment.
- Inference: conclusion derived from facts or observations.
- Opinion: preference tied to a stated context.

## Sources

1. [Hertz experiment README](./README.md) — records the experiment's research question, reproducible commands, implementation choices, and evidence migrated into this guide.
2. [Hertz routing](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/route/) — established route registration, grouping, precedence, parameters, and method-handling options.
3. [Hertz binding and validation](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/binding-and-validate/) — established binding APIs, field tags, and validation behavior.
4. [Hertz middleware](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/middleware/) — established middleware registration and handler-chain mechanics.
5. [Hertz unit testing](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/unit-test/) — established Hertz's in-process test utilities and supported test workflow.
6. [Hertz graceful shutdown](https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/graceful-shutdown/) — established signal behavior, shutdown ordering, and `WithExitWaitTime`.
7. [`hz` generated layout](https://www.cloudwego.io/docs/hertz/tutorials/toolkit/layout/) — established generated model/router files, handler ownership, and regeneration boundaries.
8. [Hertz tracing](https://www.cloudwego.io/docs/hertz/tutorials/observability/tracing/) — established tracer lifecycle, request timing events, and instrumentation constraints.
9. [RabbitMQ acknowledgements and publisher confirms](https://www.rabbitmq.com/docs/confirms) — established manual acknowledgements, bounded in-flight deliveries, requeue behavior, and redelivery expectations.
10. [Twelve-Factor disposability](https://12factor.net/disposability) — validated the production lifecycle backbone: stop new work, finish or return current work, and design jobs for retry/idempotency.

## Research limits

- Hertz is pinned to `v0.10.6`; current documentation may describe later behavior, so experiments must resolve any version-sensitive discrepancy.
- The broker unit uses RabbitMQ as one concrete lifecycle model. Kafka consumer-group rebalancing and offset commits require a separate guide or experiment.
- No benchmark or production-load evidence is included, so the guide supports no Hertz performance ranking.
