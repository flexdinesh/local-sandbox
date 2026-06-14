# Nix Project Environments

CBox will support Project Environments as workspace-owned toolchains that a Harness session can enter, with Nix flakes as the first supported implementation. Nix activation is explicit through `--project-env nix` on run commands, uses the Mounted Workspace's `flake.nix`, wraps both default Harness commands and command overrides with `nix develop --command`, enters only the default dev shell, and persists Nix state in Docker-managed CBox volumes rather than mounting the host's Nix store. CBox requires `flake.nix` when Nix mode is requested, recommends but does not require `flake.lock`, enables `nix-command flakes` in Sandbox Image Nix config, and keeps the current root container model for the initial slice. Nix state volumes are mounted only for explicit Nix Project Environment runs. CBox will not pass `--impure` by default, and Nix mode accepts standard `nix develop` behavior, including network access under the current container network model and project-defined shell hooks. This keeps Harness images responsible for Harness runtimes while letting projects define their own tools without exposing more of the base machine filesystem.

**Considered Options**

- Put project tools in Sandbox Images: rejected because it couples Harness images to each project's toolchain.
- Mount the host Nix store: rejected because it expands the host filesystem exposed to Harness containers.
- Auto-detect Nix files: rejected for the initial design because entering a Project Environment can fetch dependencies, run hooks, and alter command behavior.
