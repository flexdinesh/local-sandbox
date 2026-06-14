# Issue 005: Extend Nix Project Environment to All Run Forms

Type: AFK

## Parent

`.scratch/nix-project-environments/PRD.md`

## What to Build

Extend explicit Nix Project Environment behavior to all supported Harnesses and shorthand run commands while keeping long-form and shorthand behavior equivalent.

## Acceptance Criteria

- [ ] `cbox run opencode --project-env nix` starts `nix develop --command opencode`.
- [ ] `cbox run pi --project-env nix` starts `nix develop --command pi`.
- [ ] `cbox run codex --project-env nix` starts `nix develop --command codex`.
- [ ] `cbox opencode --project-env nix`, `cbox pi --project-env nix`, and `cbox codex --project-env nix` match their explicit run forms.
- [ ] Nix state volumes are mounted only when `--project-env nix` is selected.
- [ ] Tests cover every Harness and shorthand equivalence.

## Blocked By

- `.scratch/nix-project-environments/issues/04-wrap-command-overrides-in-nix-project-environment.md`
