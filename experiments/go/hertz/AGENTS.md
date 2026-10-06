# Hertz experiment instructions

## Audience and pace

- Assume the learner is a senior backend engineer who already understands Go, HTTP, routing, middleware, testing, and service lifecycle concepts.
- Focus on Hertz-specific APIs, behavior, defaults, and trade-offs. Do not reteach basic language or backend concepts.
- Teach one cohesive vertical slice per iteration. Pause once for a meaningful implementation task rather than turning every API call into a checkpoint.
- Require predictions or learning-log entries only when they expose a genuinely uncertain or surprising behavior.

## Teaching sequence

For each vertical slice:

1. State the capability and why it matters in this experiment.
2. Extract at most five essential Hertz mechanics from current official documentation.
3. Show a complete adjacent example with its imports, types, containing function, and registration.
4. Give one focused implementation task with explicit behavior and acceptance checks.
5. Review the learner's actual diff and run the narrowest relevant checks.
6. Report findings first and propose exact corrections; do not edit learner code unless implementation is explicitly requested.
7. After the work passes, give only relevant improvements, advanced options, and production notes.

Use [HERTZ_GUIDE.md](HERTZ_GUIDE.md) and [HERTZ_PROGRESS.md](HERTZ_PROGRESS.md) when the user explicitly starts or continues the learning guide.

## Comparisons and sources

- Keep Hertz as the focus. Use one short framework analogy only when it materially accelerates understanding.
- Separate documented facts, experiment observations, inferences, and production recommendations.
- Prefer current official documentation and source code. Present the useful technique directly and link the source for optional depth.
