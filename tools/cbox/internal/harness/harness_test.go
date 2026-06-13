package harness

import (
	"reflect"
	"testing"
)

func TestAllReturnsCanonicalHarnessesInDocumentedOrder(t *testing.T) {
	harnesses := All()

	got := make([]string, 0, len(harnesses))
	for _, h := range harnesses {
		got = append(got, h.Name)
	}

	want := []string{NameOpenCode, NamePI}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected harness names %v, got %v", want, got)
	}
}

func TestLookupReturnsExactDefinitions(t *testing.T) {
	tests := []struct {
		name       string
		imageTag   string
		dockerfile string
	}{
		{name: NameOpenCode, imageTag: "sandbox-opencode", dockerfile: "images/opencode/Dockerfile"},
		{name: NamePI, imageTag: "sandbox-pi", dockerfile: "images/pi/Dockerfile"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, ok := Lookup(tt.name)
			if !ok {
				t.Fatalf("expected lookup for %q to succeed", tt.name)
			}

			if h.Name != tt.name {
				t.Fatalf("expected name %q, got %q", tt.name, h.Name)
			}
			if h.ImageTag != tt.imageTag {
				t.Fatalf("expected image tag %q, got %q", tt.imageTag, h.ImageTag)
			}
			if h.Dockerfile != tt.dockerfile {
				t.Fatalf("expected Dockerfile %q, got %q", tt.dockerfile, h.Dockerfile)
			}
		})
	}
}

func TestLookupRejectsUnknownHarness(t *testing.T) {
	if _, ok := Lookup("unknown"); ok {
		t.Fatal("expected unknown harness lookup to fail")
	}
}

func TestBuildArgv(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: NameOpenCode,
			want: []string{"build", "-f", "images/opencode/Dockerfile", "-t", "sandbox-opencode", "."},
		},
		{
			name: NamePI,
			want: []string{"build", "-f", "images/pi/Dockerfile", "-t", "sandbox-pi", "."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, ok := Lookup(tt.name)
			if !ok {
				t.Fatalf("expected lookup for %q to succeed", tt.name)
			}

			if got := h.BuildArgv(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected build argv:\n%q\ngot:\n%q", tt.want, got)
			}
		})
	}
}

func TestRunArgv(t *testing.T) {
	tests := []struct {
		name        string
		passThrough []string
		want        []string
	}{
		{
			name:        NameOpenCode,
			passThrough: []string{"opencode", "debug"},
			want: []string{
				"run", "-it", "--rm",
				"-v", "/repo:/workdir",
				"-w", "/workdir",
				"-v", "opencode-config:/root/.config/opencode",
				"-v", "opencode-shared:/root/.local/share/opencode",
				"-v", "opencode-state:/root/.local/state/opencode",
				"-v", "/home/test/.config/opencode/opencode.jsonc:/root/.config/opencode/opencode.jsonc:ro",
				"-v", "/home/test/.config/opencode/tui.json:/root/.config/opencode/tui.json:ro",
				"-v", "/home/test/.config/opencode/plugins:/root/.config/opencode/plugins:ro",
				"-v", "/home/test/.config/opencode/prompts:/root/.config/opencode/prompts:ro",
				"-v", "/home/test/.local/share/opencode/auth.json:/root/.local/share/opencode/auth.json:ro",
				"sandbox-opencode",
				"opencode", "debug",
			},
		},
		{
			name:        NamePI,
			passThrough: []string{"pi", "--version"},
			want: []string{
				"run", "-it", "--rm",
				"-v", "/repo:/workdir",
				"-w", "/workdir",
				"-v", "shared-pi:/root/.pi",
				"-v", "/home/test/.pi/agent/extensions:/root/.pi/agent/extensions:ro",
				"-v", "/home/test/.pi/agent/auth.json:/root/.pi/agent/auth.json:ro",
				"-v", "/home/test/.pi/agent/keybindings.json:/root/.pi/agent/keybindings.json:ro",
				"-v", "/home/test/.pi/agent/settings.json:/root/.pi/agent/settings.json:ro",
				"sandbox-pi",
				"pi", "--version",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, ok := Lookup(tt.name)
			if !ok {
				t.Fatalf("expected lookup for %q to succeed", tt.name)
			}

			got := h.RunArgv("/repo", "/home/test", tt.passThrough)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected run argv:\n%q\ngot:\n%q", tt.want, got)
			}
		})
	}
}

func TestRunArgvWithoutPassThroughEndsAtImageName(t *testing.T) {
	h, ok := Lookup(NameOpenCode)
	if !ok {
		t.Fatalf("expected lookup for %q to succeed", NameOpenCode)
	}

	got := h.RunArgv("/repo", "/home/test", nil)
	if got[len(got)-1] != "sandbox-opencode" {
		t.Fatalf("expected run argv to end at image name, got %q", got)
	}
}

func TestDefinitionsAreCopied(t *testing.T) {
	harnesses := All()
	harnesses[0].Volumes[0].Name = "changed"

	h, ok := Lookup(NameOpenCode)
	if !ok {
		t.Fatalf("expected lookup for %q to succeed", NameOpenCode)
	}

	if h.Volumes[0].Name != "opencode-config" {
		t.Fatalf("expected definitions to be copied, got volume %q", h.Volumes[0].Name)
	}
}
