# Backend lab

This is a personal R&D lab for learning backend frameworks and platform-agnostic topics across programming languages. Each experiment follows a small cycle: read relevant primary sources, try a concrete backend scenario, evaluate the behavior and developer experience, record evidence and provisional opinions, then compare implementations only when their tested requirements are sufficiently similar. The goal is better technology judgment; finishing sample applications is secondary.

## Current layout

```text
backend-lab/
├── AGENTS.md
├── README.md
├── docs/agent-workflows/
│   ├── planning.md
│   ├── building.md
│   ├── review.md
│   └── verification.md
└── experiments/
    └── go/
        └── hertz/
```

- [Hertz experiment](experiments/go/hertz/README.md) is active; its README records its question, commands, evidence, and remaining work.
- [Agent guidance](AGENTS.md) contains standing rules and links to the task-specific workflows.

Future experiments belong under `experiments/<language>/<experiment>/`, with their own dependency files, source, tests, and README. Use lowercase kebab-case directory names unless the language ecosystem strongly prefers another convention.

Platform-agnostic learning and research notes belong under `topics/<topic>.md`, such as `topics/oauth2.md`. Keep runnable implementations in `experiments/`. Add `scenarios/<scenario>.md` when a language-neutral behavior contract is needed, and `docs/comparisons.md` when completed experiments provide comparable evidence. These paths do not exist yet. A scenario defines inputs, outputs, failure behavior, and acceptance checks without prescribing framework architecture.

## Experiment README contract

Each experiment README should record:

1. **Research question:** What the experiment is intended to learn.
2. **Status:** Planned, active, complete, or revisiting.
3. **Versions:** Framework, language, runtime, and important tools.
4. **Scenario:** Shared scenario and any deviations.
5. **Setup and commands:** Exact install, run, and test steps.
6. **Implementation notes:** Significant defaults, generated code, and design choices.
7. **Findings:** Verified behavior and direct observations.
8. **Strengths and limitations:** Contextual trade-offs supported by findings.
9. **Verdict:** Appropriate uses and remaining uncertainty.
10. **References:** Primary sources first.

Keep useful raw learning notes, but distinguish them from the current conclusion.

## Definition of done

An experiment is complete for its stated scope when the research question is explicit, the selected scenario works, relevant checks pass, versions and commands are reproducible, findings separate evidence from opinion, strengths and limitations have context, and remaining uncertainty is recorded. Update a cross-framework comparison when there is enough equivalent evidence. Stop there; extend the experiment only for a new research question.
