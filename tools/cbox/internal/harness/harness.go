package harness

import "path/filepath"

const (
	NameOpenCode = "opencode"
	NamePI       = "pi"
	NameCodex    = "codex"

	workdirContainerPath = "/workdir"
)

type Harness struct {
	Name       string
	ImageTag   string
	Dockerfile string
	Volumes    []Volume
	HomeMounts []HomeMount
}

type Volume struct {
	Name          string
	ContainerPath string
}

type HomeMount struct {
	RelativePath  string
	ContainerPath string
	ReadOnly      bool
}

type Workspace struct {
	Mounts     []WorkspaceMount
	WorkingDir string
}

type WorkspaceMount struct {
	HostPath      string
	ContainerPath string
}

var definitions = []Harness{
	{
		Name:       NameOpenCode,
		ImageTag:   "sandbox-opencode",
		Dockerfile: "images/opencode/Dockerfile",
		Volumes: []Volume{
			{Name: "opencode-config", ContainerPath: "/root/.config/opencode"},
			{Name: "opencode-shared", ContainerPath: "/root/.local/share/opencode"},
			{Name: "opencode-state", ContainerPath: "/root/.local/state/opencode"},
		},
		HomeMounts: []HomeMount{
			{RelativePath: ".config/opencode/opencode.jsonc", ContainerPath: "/root/.config/opencode/opencode.jsonc", ReadOnly: true},
			{RelativePath: ".config/opencode/tui.json", ContainerPath: "/root/.config/opencode/tui.json", ReadOnly: true},
			{RelativePath: ".config/opencode/plugins", ContainerPath: "/root/.config/opencode/plugins", ReadOnly: true},
			{RelativePath: ".config/opencode/prompts", ContainerPath: "/root/.config/opencode/prompts", ReadOnly: true},
			{RelativePath: ".local/share/opencode/auth.json", ContainerPath: "/root/.local/share/opencode/auth.json", ReadOnly: true},
		},
	},
	{
		Name:       NamePI,
		ImageTag:   "sandbox-pi",
		Dockerfile: "images/pi/Dockerfile",
		Volumes: []Volume{
			{Name: "shared-pi", ContainerPath: "/root/.pi"},
		},
		HomeMounts: []HomeMount{
			{RelativePath: ".pi/agent/extensions", ContainerPath: "/root/.pi/agent/extensions", ReadOnly: true},
			{RelativePath: ".pi/agent/auth.json", ContainerPath: "/root/.pi/agent/auth.json", ReadOnly: true},
			{RelativePath: ".pi/agent/keybindings.json", ContainerPath: "/root/.pi/agent/keybindings.json", ReadOnly: true},
			{RelativePath: ".pi/agent/settings.json", ContainerPath: "/root/.pi/agent/settings.json", ReadOnly: true},
		},
	},
	{
		Name:       NameCodex,
		ImageTag:   "sandbox-codex",
		Dockerfile: "images/codex/Dockerfile",
		HomeMounts: []HomeMount{
			{RelativePath: ".codex", ContainerPath: "/root/.codex"},
		},
	},
}

func All() []Harness {
	harnesses := make([]Harness, len(definitions))
	for i, definition := range definitions {
		harnesses[i] = clone(definition)
	}

	return harnesses
}

func Lookup(name string) (Harness, bool) {
	for _, definition := range definitions {
		if definition.Name == name {
			return clone(definition), true
		}
	}

	return Harness{}, false
}

func (h Harness) BuildArgv() []string {
	return []string{"build", "-f", h.Dockerfile, "-t", h.ImageTag, "."}
}

func (h Harness) RunArgv(workdir, homeDir string, passThrough []string) []string {
	return h.RunArgvWithWorkspace(Workspace{
		Mounts: []WorkspaceMount{
			{HostPath: workdir, ContainerPath: workdirContainerPath},
		},
		WorkingDir: workdirContainerPath,
	}, homeDir, passThrough)
}

func (h Harness) RunArgvWithWorkspace(workspace Workspace, homeDir string, passThrough []string) []string {
	argv := []string{
		"run",
		"-it",
		"--rm",
	}

	for _, mount := range workspace.Mounts {
		argv = append(argv, "-v", mount.HostPath+":"+mount.ContainerPath)
	}

	argv = append(argv, "-w", workspace.WorkingDir)

	for _, volume := range h.Volumes {
		argv = append(argv, "-v", volume.Name+":"+volume.ContainerPath)
	}

	for _, mount := range h.HomeMounts {
		value := filepath.Join(homeDir, mount.RelativePath) + ":" + mount.ContainerPath
		if mount.ReadOnly {
			value += ":ro"
		}
		argv = append(argv, "-v", value)
	}

	argv = append(argv, h.ImageTag)
	argv = append(argv, passThrough...)

	return argv
}

func (h Harness) ManagedContainerPaths() []string {
	paths := make([]string, 0, len(h.Volumes)+len(h.HomeMounts))
	for _, volume := range h.Volumes {
		paths = append(paths, volume.ContainerPath)
	}
	for _, mount := range h.HomeMounts {
		paths = append(paths, mount.ContainerPath)
	}
	return paths
}

func clone(h Harness) Harness {
	h.Volumes = append([]Volume(nil), h.Volumes...)
	h.HomeMounts = append([]HomeMount(nil), h.HomeMounts...)
	return h
}
