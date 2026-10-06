# Building an experiment

Use this guidance when implementing an experiment or changing its setup, dependencies, or generated code.

- Implement only enough to answer the current research question. Follow the framework's documented, idiomatic path and significant defaults before introducing custom architecture. Record defaults that affect the result.
- Prefer the standard library and officially maintained integrations. Add third-party dependencies only when the scenario requires them or the dependency itself is under evaluation. Add databases, containers, queues, or observability only when the question requires them.
- Keep each framework project self-contained. Similar duplication is acceptable when it makes implementations independent and comparable; do not share cross-framework application code.
- Each experiment owns its manifest and lock file. Pin dependencies through the ecosystem's normal mechanism. Keep Go `go.mod` inside its framework directory; add no root `go.mod` or `go.work` unless multiple Go modules need local imports for a current experiment. Keep JavaScript, TypeScript, and Rust package files inside their experiment directories. Commit lock files where the ecosystem expects them.
- Run build, test, lint, and generation commands from the experiment directory unless its README says otherwise. Preserve generated code when it is part of the normal workflow, record its generation command, and do not hand-edit it.
- Document example environment variables when configuration is needed, without committing secrets. Add a language- or framework-specific `AGENTS.md` only when substantial local rules justify it. Do not introduce a repository-wide build system until repeated manual work demonstrates a concrete need.
