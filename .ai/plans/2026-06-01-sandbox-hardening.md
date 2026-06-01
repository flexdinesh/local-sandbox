# Plan: Wrapper Refactor, Version Pins, Filesystem Hardening

## Summary

Keep the current image hierarchy for now, preserve the current config-sharing model, and improve maintainability by extracting shared wrapper logic, pinning package versions, and hardening container filesystem defaults.

## Decisions Made

- Keep current image hierarchy for now. No `sandbox-node` in this pass.
- Keep current config model: writable Docker volumes plus selected host config/auth overlays mounted read-only.
- Keep host pnpm store bind mount so it can be shared with host OS and other containers.
- No hardening escape hatch.
- Do not switch CLI execution to non-root in this pass.
- Add root `versions.env` as version source of truth.
- Do not implement strict network enforcement yet.

## Key Changes

1. Add `versions.env` with pinned base image, node image, OpenCode version, and PI version.
2. Build scripts source `versions.env` and pass build args to Docker.
3. Dockerfiles use pinned package installs.
4. Add `scripts/lib/run-tool.sh` for common wrapper behavior:
   - argument parsing
   - `HOST_DIR`, `CONTAINER_WORKDIR`, `NETWORK_ACCESS`
   - TTY detection
   - proxy mode env
   - required mount validation
   - symlink resolution
   - common Docker hardening args
5. Reduce `run-opencode.sh` and `run-pi.sh` to tool-specific config plus the shared runner.
6. Add wrapper hardening defaults:
   - `--read-only`
   - `--security-opt no-new-privileges`
   - `--tmpfs /tmp:exec`
   - `--tmpfs /run`
   - `--tmpfs /var/log`
   - `--tmpfs /root/.cache`
7. Keep writable:
   - `$HOST_DIR:$CONTAINER_WORKDIR`
   - host pnpm store mount
   - existing tool Docker named volumes
   - ephemeral `/root/.cache` tmpfs for CLI runtime cache
   - ephemeral `/tmp` tmpfs with `exec` so OpenCode/OpenTUI can load extracted native libraries
   - selected current host config/auth overlays read-only
8. Update the entrypoint so it copies Tinyproxy config to `/run/tinyproxy/tinyproxy.conf`, mutates that runtime copy, and starts Tinyproxy with it.
9. Add a Tinyproxy readiness loop after starting `supervisord`.
10. Keep README concise and explicit that network sandboxing is proxy allowlisting, not strict packet-level egress.

## Network Note

Strict no-sidecar egress enforcement would require in-container firewall rules with `NET_ADMIN`, separate users for Tinyproxy and the CLI, and owner-based outbound rules. That is intentionally out of scope for this pass.

## Verification

Run:

```bash
bash -n scripts/*.sh scripts/lib/*.sh images/base/entrypoint.sh
./scripts/build.sh
docker run --rm sandbox-base cat /etc/debian_version
docker run --rm sandbox-opencode node --version
docker run --rm sandbox-opencode opencode --version
docker run --rm sandbox-pi pi --version
docker run --rm sandbox-base sh -c "sleep 3; supervisorctl -c /etc/supervisor/conf.d/tinyproxy.conf status tinyproxy"
./scripts/run-opencode.sh -- opencode --version
./scripts/run-pi.sh -- pi --version
```

Also verify restricted mode blocks a non-allowlisted domain through the proxy, and full mode allows it.

## Execution Guidance

If implementation discovers a needed deviation, update this plan before proceeding and call out the deviation clearly.

## Deviation: OpenTUI native library loading

OpenCode/OpenTUI extracts a native render library to `/tmp` and loads it with executable mappings. Docker tmpfs mounts are commonly `noexec` unless `exec` is specified, which causes errors like `failed to map segment from shared object` when loading the extracted `.so`.

Keep `/tmp` as tmpfs, but mount it with `exec`: `--tmpfs /tmp:exec`.
