# Planning an experiment

Use this guidance when starting a framework experiment, defining a research question, or designing a shared scenario.

## Define the question

- State the specific behavior or trade-off to investigate before adding substantial code. Examples include routing and middleware organization, request validation, graceful shutdown, testing, or performance.
- Avoid open-ended goals such as learning an entire framework. Choose the smallest experiment that can answer the question.
- Ask for a prediction before testing when the expected behavior is genuinely uncertain or surprising.

## Ground the plan

- Prefer current official documentation, source code, release notes, and maintained official examples. Identify the framework's intended design and use cases. Use third-party sources only if primary sources do not answer the question, and label them.
- Record framework, language, runtime, and important tool versions. Record access dates for web sources when behavior may change.
- Treat tutorial claims as hypotheses until documented or reproduced when practical.
- Start with the framework's documented, idiomatic path. Defer custom architecture until the default behavior is understood.

## Scope the scenario

A shared scenario is a behavior contract, not a prescribed architecture. Specify inputs, outputs, failure behavior, and acceptance checks while leaving room for each framework's normal idioms. Begin with small backend concerns such as health and routing; path, query, and JSON input; validation and errors; middleware and request context; handler or integration testing; and graceful shutdown. Add authentication, persistence, observability, concurrency, or deployment only when they become explicit research targets.
