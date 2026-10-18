---
name: run-learning-guide
description: Teach a valid target-named _GUIDE.md interactively across resumable, time-boxed sessions while maintaining adjacent matching _PROGRESS.md state. Use only when explicitly invoked to start, continue, or pause an existing learning guide; do not use it to design the curriculum.
---

# Run Learning Guide

Teach from a stable guide while making the learner's state durable enough to continue in a different chat. Treat each guide unit as a **Learning Item**. Set the knowledge gap, teach at a depth that fits the item's complexity, verify the guide's mastery evidence, and record the result.

## Open the session

Resolve the guide from a user-supplied path. A guide filename must end in `_GUIDE.md`. Its progress filename uses the exact same prefix with `_PROGRESS.md` instead; for example, `SAMPLE_TOPIC_GUIDE.md` pairs with `SAMPLE_TOPIC_PROGRESS.md`. If no path is supplied, use the single unambiguous `*_GUIDE.md` candidate in the active workspace; otherwise ask which guide to use. Read it completely.

A valid guide has a desired outcome, success criteria, total budget, and an ordered roadmap whose units have unique IDs, mastery evidence, and time estimates. If it is missing or invalid, stop and recommend `$build-learning-guide`; do not silently invent or repair the curriculum.

Keep the guide read-only. Immediately after validating it, derive the adjacent progress path by replacing the guide's trailing `_GUIDE.md` with `_PROGRESS.md`. If progress does not exist, read [assets/PROGRESS.template.md](assets/PROGRESS.template.md) and create it from the guide's filename, item IDs, and titles. If it exists, read it completely before responding and preserve its history.

Require the learner's available active-learning time for this session. Ask if it is missing, but do not delay progress setup. Do not begin teaching until the time is known.

For a resumed session, start with no more than three short bullets:

- the last completed work and evidence;
- recent incomplete or uncertain knowledge; and
- the exact next action and this session's allocation.

Do not recap concepts the learner has already demonstrated unless the checkpoint marks them as fragile or the learner requests a review.

## Prepare only the current Learning Item

Choose work that fits the session budget and reserve a small closing checkpoint. Treat the guide's item estimate and the learner's session budget as real limits. If required work will exceed either limit, stop at a checkpoint; ask before increasing the item estimate. Follow the guide's order unless the matching progress file records a learner-driven adaptation.

Treat the Learning Item's outline as a high-level baseline, not a complete lesson. Research and explain the in-scope details needed for its mastery evidence. Prepare around the delta between that evidence and the learner's demonstrated knowledge. Prioritize the governing mental model, non-obvious behavior, realistic failure modes, and production-relevant trade-offs. Omit introductory background the learner has already demonstrated. Do not add questions or examples merely to make the session feel interactive.

Research only the current item's scope. If the learner's answer exposes a prerequisite outside that scope, ask permission before researching or teaching it. Favor authoritative, current sources and verify unstable facts. Do not delay a lesson to research settled background that is unnecessary for the mastery evidence. Do not broaden the curriculum or turn research into a new guide. If necessary research is unavailable, disclose that limitation, use the guide and accessible sources, and record the limitation in the session log.

Bring the relevant knowledge to the learner: research and explain source material yourself, with a useful link when appropriate. Do not assign a documentation search merely to verify a fact you could provide. Use source-finding practice only when the learner wants that skill. If a guide requires a lookup and the learner opts for direct instruction, record the learner-driven adaptation in progress and verify the underlying concept through a meaningful application.

## Teach according to complexity

Use this sequence for every Learning Item:

1. **Set the delta.** State what the learner should understand or be able to do after the item. Check existing knowledge only when it changes what needs teaching.
2. **Choose the teaching path.** Use the simple path when one straightforward concept and a short check can satisfy the mastery evidence. Use the complex path when the evidence requires multiple connected parts, approaches, or steps.
3. **Teach and verify.** Follow the matching path below, then check the guide's stated mastery evidence or a learner-driven equivalent recorded in progress.
4. **Review taught knowledge.** Check that what was taught covers the Learning Item's outcome and in-scope key ideas to the standard expected for the topic, and that the mastery evidence was demonstrated. If either falls short, return to step 3 to teach and verify the missing part within the time limits.
5. **Close.** Mark the item `completed` only after the learner demonstrates that evidence. Record the evidence, unresolved details, and exact next action.

### Turn-taking

Teach one coherent part in the main conversation, then give the learner one meaningful task or question and end the turn. Wait for the learner's response before giving feedback or teaching the next part. Address any uncertainty before continuing. Keep explanations and learner prompts in the conversation history; do not use forms or bundled question flows for teaching interactions.

### Simple Learning Items

- Show the knowledge, code, or commands directly. Explain the rule and the smallest useful example before asking the learner to respond. Do not require a prediction or guess.
- Offer full or partial hands-on practice only when it adds value; practice is optional. Briefly surface the important edge cases and realistic failure modes.
- Verify understanding quickly with a short-answer question, a concise quiz, or a small direct application that checks the item's mastery evidence. Do not quiz trivial syntax or facts merely for recall. Ask more only after an incorrect, incomplete, or uncertain response.
- Favor speed and clarity. If the evidence is demonstrated, close the item instead of adding reassurance questions.

### Complex Learning Items

- Explain the core model first: the governing rule, runtime behavior, and a useful example. Separate guarantees from implementation details and version-dependent behavior.
- Teach each meaningful part of the item one by one, using hands-on practice when it helps demonstrate the mastery evidence. Cover important in-scope edge cases, realistic failure modes, and trade-offs needed for that evidence within the item estimate. If multiple in-scope approaches exist, show a small comparison table.
- Use a quick prediction or diagnostic question when it can expose an important misconception. Prefer fast assessments where they still demonstrate deep understanding; use longer work when the mastery evidence requires it. Move on after a clear correct answer rather than repeating equivalent checks.
- When starting a complex item, create or update one clean, concise note file for it beside the guide. Follow the guide folder's naming convention and choose a name that makes the note's topic clear. If the folder has no convention, use a concise, descriptive Markdown filename. Capture the core model, key examples or code, in-scope edge cases, and any approach comparison. Keep it as an independent knowledge reference; do not link it from or copy it into progress. Reuse the file if the item spans sessions.

For either path, prefer meaningful learner work over a long lecture or low-value prompts. Self-reported confidence is not mastery evidence. Evidence may be an explanation, correct application, produced artifact, or successful diagnosis appropriate to the subject. Keep conversational turns proportional to the item estimate.

Do not skip or defer an entire Learning Item based only on an inference that the learner knows it. Ask for and wait for explicit learner confirmation first; a direct learner request to skip or defer counts as confirmation. Mark a confirmed omission `skipped` and a planned return `deferred` in progress. Record the decision and reason, then move to the next eligible item. A skipped item stays unchecked and does not count as mastery evidence. Omitting already-demonstrated background within an item does not skip the item or waive its mastery evidence. Adapt pacing in progress, not in the guide. Record skipped material, repetition, reordered items, estimate changes, and their reasons. Do not introduce a new item unless the guide is rebuilt.

## Persist the learning state

Use exactly these item statuses: `not started`, `in progress`, `completed`, `blocked`, `deferred`, and `skipped`. Only completed items receive checked checklist boxes.

Update the current checkpoint after every completed item and whenever the learner pauses or the session closes. Capture concise learning state rather than a transcript. Do not record routine correct answers as fragile knowledge or create retrieval prompts for material already mastered:

- evidence demonstrated;
- recent misconceptions or fragile knowledge;
- the exact concept, example, or exercise left incomplete;
- open questions and short retrieval prompts; and
- the exact prompt or action that should start the next session.

Before ending a teaching turn with learner work, settle the exercise and make the progress file's `Next action` and `Exact re-entry prompt` match the task shown to the learner, including any counts or identifiers. Save enough of a conversation-only exercise in the checkpoint to make the re-entry prompt actionable in another chat; do not refer to absent content as "above" or "below." If the exercise changes while composing the response, update progress before sending it.

Append one chronological session-log entry when a session pauses or ends. Record planned time and actual active time only when reported or reliably known; never invent actual time. Add only known actual time to the accumulated total. Preserve prior log entries.

When the available time is exhausted, do not overrun or compress remaining teaching into a summary. If the planned work finishes early and its mastery evidence is satisfied, close the item instead of filling the remaining time with extra questions. If required complex-item coverage remains when time runs out, leave the item `in progress` and resume it next session. Save the checkpoint, append the log, and stop. End with the progress path and the exact next step.
