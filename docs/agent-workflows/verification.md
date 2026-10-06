# Verifying an experiment

Use this guidance when testing, reproducing behavior, changing run instructions, or reporting benchmark results.

- There is no assumed repository-wide test command. Discover and run checks for the affected experiment from its directory, unless its README specifies otherwise.
- Run the narrowest command that proves the behavior. Check meaningful success, failure, and edge cases for the current scenario. Run broader tests only when they address a remaining risk. Format and lint with ecosystem-standard tools when configured.
- Capture the command, relevant configuration, and result needed to reproduce each finding. Inspect generated code and defaults when they materially affect the result. Verify documented run commands after changing setup or dependencies.
- Use safe timeouts for servers and blocking examples, and stop background processes after verification. External services, containers, and paid resources require an explicit research need and user authorization.
- Measure before making performance claims. Record benchmark command, workload, environment, sample count, and limitations. Treat marketing claims and isolated microbenchmarks as claims, not conclusions; avoid conclusions from a single run.
