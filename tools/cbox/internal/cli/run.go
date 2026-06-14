package cli

import (
	"fmt"
	containerpath "path"
	"path/filepath"
	"strings"

	"github.com/flexdinesh/cbox/tools/cbox/internal/harness"
	"github.com/spf13/cobra"
)

func newRunCommand(cfg config) *cobra.Command {
	var workspaceMounts []string
	cmd := &cobra.Command{
		Use:   "run <harness> [-- command...]",
		Short: "Run a Sandbox Image in the foreground",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, passThrough, err := parseRunArgs(cmd, args)
			if err != nil {
				return err
			}

			return runHarness(cmd, cfg, h, workspaceMounts, passThrough)
		},
	}
	cmd.Flags().StringArrayVar(&workspaceMounts, "workspace-mount", nil, "Mount a host directory tree into the Sandbox Image")

	return cmd
}

func newShorthandRunCommand(cfg config, name string) *cobra.Command {
	var workspaceMounts []string
	cmd := &cobra.Command{
		Use:   name + " [-- command...]",
		Short: "Run the " + name + " Sandbox Image in the foreground",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, ok := harness.Lookup(name)
			if !ok {
				return fmt.Errorf("invalid Harness %q (valid Harnesses: %s)", name, strings.Join(validHarnessNames(harness.All()), ", "))
			}
			if cmd.Flags().ArgsLenAtDash() < 0 && len(args) > 0 {
				return fmt.Errorf("container commands must be passed after --")
			}

			return runHarness(cmd, cfg, h, workspaceMounts, args)
		},
	}
	cmd.Flags().StringArrayVar(&workspaceMounts, "workspace-mount", nil, "Mount a host directory tree into the Sandbox Image")

	return cmd
}

func parseRunArgs(cmd *cobra.Command, args []string) (harness.Harness, []string, error) {
	dash := cmd.Flags().ArgsLenAtDash()
	if len(args) == 0 {
		return harness.Harness{}, nil, fmt.Errorf("missing Harness (valid Harnesses: %s)", strings.Join(validHarnessNames(harness.All()), ", "))
	}
	if dash < 0 && len(args) > 1 {
		return harness.Harness{}, nil, fmt.Errorf("container commands must be passed after --")
	}

	h, ok := harness.Lookup(args[0])
	if !ok {
		return harness.Harness{}, nil, fmt.Errorf("invalid Harness %q (valid Harnesses: %s)", args[0], strings.Join(validHarnessNames(harness.All()), ", "))
	}

	var passThrough []string
	if dash >= 0 {
		passThrough = args[1:]
	}

	return h, passThrough, nil
}

func runHarness(cmd *cobra.Command, cfg config, h harness.Harness, workspaceMountSpecs []string, passThrough []string) error {
	workdir, err := cfg.workingDir()
	if err != nil {
		return fmt.Errorf("failed to resolve current directory: %w", err)
	}

	homeDir, err := cfg.homeDir()
	if err != nil {
		return fmt.Errorf("failed to resolve user home directory: %w", err)
	}

	workspace, err := resolveWorkspace(workdir, homeDir, h, workspaceMountSpecs)
	if err != nil {
		return err
	}

	return cfg.runner.Run(cmd.Context(), h.RunArgvWithWorkspace(workspace, homeDir, passThrough))
}

func resolveWorkspace(workdir, homeDir string, h harness.Harness, specs []string) (harness.Workspace, error) {
	workdir = filepath.Clean(workdir)
	workspace := harness.Workspace{
		Mounts:     make([]harness.WorkspaceMount, 0, len(specs)+1),
		WorkingDir: "/workdir",
	}

	seenHostPaths := map[string]struct{}{}
	seenContainerPaths := map[string]struct{}{}
	var best *harness.WorkspaceMount
	var bestRel string
	for _, spec := range specs {
		mount, err := parseWorkspaceMount(spec, workdir, homeDir)
		if err != nil {
			return harness.Workspace{}, err
		}
		if err := validateWorkspaceMount(mount, h, seenHostPaths, seenContainerPaths); err != nil {
			return harness.Workspace{}, err
		}
		workspace.Mounts = append(workspace.Mounts, mount)
		seenHostPaths[mount.HostPath] = struct{}{}
		seenContainerPaths[mount.ContainerPath] = struct{}{}

		rel, ok := relativeTo(workdir, mount.HostPath)
		if !ok {
			continue
		}
		if best == nil || len(mount.HostPath) > len(best.HostPath) {
			copy := mount
			best = &copy
			bestRel = rel
		}
	}

	if best != nil {
		workspace.WorkingDir = containerpath.Join(best.ContainerPath, filepath.ToSlash(bestRel))
		return workspace, nil
	}

	workspace.Mounts = append([]harness.WorkspaceMount{
		{HostPath: workdir, ContainerPath: "/workdir"},
	}, workspace.Mounts...)
	return workspace, nil
}

func validateWorkspaceMount(mount harness.WorkspaceMount, h harness.Harness, seenHostPaths, seenContainerPaths map[string]struct{}) error {
	if _, ok := seenHostPaths[mount.HostPath]; ok {
		return fmt.Errorf("duplicate Workspace Mount host path %q", mount.HostPath)
	}
	if _, ok := seenContainerPaths[mount.ContainerPath]; ok {
		return fmt.Errorf("duplicate Workspace Mount container path %q", mount.ContainerPath)
	}
	for containerPath := range seenContainerPaths {
		if containerPathsOverlap(mount.ContainerPath, containerPath) {
			return fmt.Errorf("Workspace Mount container path %q overlaps %q", mount.ContainerPath, containerPath)
		}
	}
	if containerPathsOverlap(mount.ContainerPath, "/workdir") {
		return fmt.Errorf("Workspace Mount container path %q conflicts with the fallback Mounted Workspace", mount.ContainerPath)
	}
	for _, managedPath := range h.ManagedContainerPaths() {
		if containerPathsOverlap(mount.ContainerPath, managedPath) {
			return fmt.Errorf("Workspace Mount container path %q overlaps Harness-managed path %q", mount.ContainerPath, managedPath)
		}
	}
	return nil
}

func parseWorkspaceMount(spec, workdir, homeDir string) (harness.WorkspaceMount, error) {
	if strings.Count(spec, ":") != 1 {
		return harness.WorkspaceMount{}, fmt.Errorf("Workspace Mount must use HOST_PATH:CONTAINER_PATH")
	}

	parts := strings.SplitN(spec, ":", 2)
	hostPath, containerPath := parts[0], parts[1]
	if hostPath == "" || containerPath == "" {
		return harness.WorkspaceMount{}, fmt.Errorf("Workspace Mount must use HOST_PATH:CONTAINER_PATH")
	}
	if !containerpath.IsAbs(containerPath) {
		return harness.WorkspaceMount{}, fmt.Errorf("Workspace Mount container path must be absolute")
	}

	hostPath = expandHome(hostPath, homeDir)
	if !filepath.IsAbs(hostPath) {
		hostPath = filepath.Join(workdir, hostPath)
	}

	return harness.WorkspaceMount{
		HostPath:      filepath.Clean(hostPath),
		ContainerPath: containerpath.Clean(containerPath),
	}, nil
}

func expandHome(path, homeDir string) string {
	if path == "~" {
		return homeDir
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

func relativeTo(path, base string) (string, bool) {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return "", false
	}
	if rel == "." {
		return "", true
	}
	if rel == ".." || strings.HasPrefix(rel, "../") {
		return "", false
	}
	return rel, true
}

func containerPathsOverlap(a, b string) bool {
	return sameOrDescendantContainerPath(a, b) || sameOrDescendantContainerPath(b, a)
}

func sameOrDescendantContainerPath(path, base string) bool {
	if path == base {
		return true
	}
	if base == "/" {
		return strings.HasPrefix(path, "/")
	}
	return strings.HasPrefix(path, base+"/")
}
