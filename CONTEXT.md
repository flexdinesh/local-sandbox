# Local Sandbox Images

This context covers standalone Docker images used to run local harnesses against a host directory with filesystem boundaries.

## Language

**Sandbox Image**:
A standalone Docker image that contains one Harness and its runtime dependencies.
_Avoid_: Base image, shared image, dev image

**Harness**:
A named agent runtime profile backed by a Sandbox Image and Manual Docker Commands.
_Avoid_: Agent CLI, intractor CLI, tool, app

**Mounted Workspace**:
The host directory tree mounted into a Sandbox Image as the primary filesystem context for a Harness run.
_Avoid_: Project folder, working folder, bind target

**Workspace Mount**:
A caller-selected host directory tree mounted into a Sandbox Image as part of the filesystem context for a Harness run.
_Avoid_: Docker mount, volume, bind target

**Manual Docker Commands**:
The Docker build and run commands documented as the source of truth for Sandbox Image behavior.
_Avoid_: noCLI commands, examples, scripts
