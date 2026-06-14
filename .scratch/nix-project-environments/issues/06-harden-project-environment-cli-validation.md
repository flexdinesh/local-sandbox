# Issue 006: Harden Project Environment CLI Validation

Type: AFK

## Parent

`.scratch/nix-project-environments/PRD.md`

## What to Build

Finish the CLI validation contract for Project Environment selection so unsupported values fail strictly, run-only scope is preserved, and future Project Environment backends have a clear extension point.

## Acceptance Criteria

- [ ] Unsupported Project Environment values fail before invoking Docker.
- [ ] Error output lists `nix` as the supported Project Environment.
- [ ] Omitting `--project-env` keeps the current plain Harness behavior.
- [ ] `--project-env` is not accepted by build commands.
- [ ] `--project-env` must be provided before `--`; arguments after `--` remain container command arguments.
- [ ] Tests cover unsupported values, omitted Project Environment behavior, and build command rejection.

## Blocked By

- `.scratch/nix-project-environments/issues/05-extend-nix-project-environment-to-all-run-forms.md`
