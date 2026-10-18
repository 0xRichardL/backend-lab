---
name: build-learning-guide
description: Research any learning target and create a concise, time-boxed guide named with an uppercase target prefix and _GUIDE.md suffix. Use only when explicitly invoked to build a learning guide; do not use it to teach lessons or track progress.
---

# Build Learning Guide

Create the stable learning compass. Optimize for complete coverage of the high-value backbone, not exhaustive detail.

## Required inputs

Before researching or writing, establish:

- the learning target;
- the learner's desired outcome;
- the learner's current level; and
- the total active-learning time budget.

Ask one concise question for any missing required input. Do not create the guide until all four are known. Treat the budget as learner study time, not agent research time. Capture optional deadlines, source restrictions, accessibility needs, or required tools only when the user supplies them or they materially affect the plan.

Derive a lowercase hyphenated target slug and convert it to an uppercase snake-case filename prefix. Name the guide `<TARGET>_GUIDE.md`; for example, `sample-topic` becomes `SAMPLE_TOPIC_GUIDE.md`. Use the user's destination directory when specified. Otherwise use `learning/<target-slug>/<TARGET>_GUIDE.md` in the active workspace. If the user supplies an explicit filename, require the `_GUIDE.md` suffix and preserve its prefix exactly. Inspect the destination before writing. Never overwrite an existing guide; suggest a sibling destination such as `<target-slug>-v2` and wait for the user's confirmation.

## Research with a stopping rule

Inspect relevant user-provided and local material first, then use available research tools unless the user prohibited external research.

1. **Coverage pass:** identify the outcomes, core concepts, dependencies, and representative practice from primary or authoritative sources.
2. **Scope checkpoint:** present a concise proposal that explicitly separates the outcome, must-learn areas, and deliberately deferred areas in light of the learner's level and time budget. Surface consequential assumptions and ask the learner or human guide builder to confirm or correct the scope. End the turn and wait for their reply; your own scope proposal is not confirmation. The human judges fit with their goal, not the subject's curriculum. Skip this checkpoint only if the human already explicitly agreed on the scope in this conversation.
3. **Validation pass:** after scope confirmation, compare the agreed backbone with established curricula, syllabi, standards, or respected practitioner guidance and close material gaps. If validation reveals a material scope change, discuss it with the human before proceeding.

Prefer current primary sources. Use secondary sources mainly to validate coverage or provide practical perspective. Keep roughly 5–10 high-value sources when the subject supports that range; do not pad the list. Stop after the two passes, or earlier when another pass adds no essential outcome. If network research is unavailable or the target is private/local, use the accessible corpus and state the limitation in the guide.

## Design the time-boxed backbone

- Classify material as must learn, supporting, or deferred. Exclude low-value detail explicitly rather than hiding it.
- Create ordered, independently resumable units with stable IDs `LG-01`, `LG-02`, and so on.
- Give every unit one observable outcome, a concise list of the specific knowledge needed for it, suitable practice, mastery evidence, and an active-learning estimate. Name the in-scope concepts, components, operations, or procedures where relevant; a broad topic label alone does not define what the learner will cover. Check that the evidence for each unit would reveal whether the learner understands or can use every essential element listed in that unit; one integrated task can cover several elements.
- Include retrieval, application, and review time. The sum of all unit estimates must not exceed the user's total budget.
- If the target cannot fit honestly, narrow the promised outcome and list the rest as deferred. Do not use unrealistically small estimates to force a fit.
- Keep explanations skeletal. The teaching skill researches and expands one unit at a time.

## Review and write only the guide

Read [assets/GUIDE.template.md](assets/GUIDE.template.md) and draft the full guide in the conversation, replacing every placeholder. Keep the compact reference useful after learning is complete. Cite sources with direct links and a short note about what each source established.

Review the whole draft against the agreed outcome and learner level: check coverage, sequence and prerequisites, evidence for essential knowledge, credible estimates, sources, and visible deferrals. Fix gaps before presenting it. Show the full draft with concise strengths and meaningful limitations or trade-offs; do not invent drawbacks. Ask the learner or human guide builder for edits or explicit approval, then end the turn and wait. If they request changes, revise and present the draft again. Scope confirmation and your own review are not approval of the guide.

Before requesting approval, verify that required inputs and success criteria are explicit, every unit lists concrete in-scope knowledge with matching mastery evidence, unit IDs are unique, unit estimates add up within the budget, deferred material is visible, and no placeholder remains. After the human approves the full draft, inspect the destination again and create only `<TARGET>_GUIDE.md`. Do not create the matching `<TARGET>_PROGRESS.md`, lesson transcripts, exercises, or source archives. Never overwrite a pre-existing guide. Report the created path, planned total, and first unit concisely.
