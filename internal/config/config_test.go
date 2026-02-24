package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vtbarreto/CLinicius/internal/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
		want    *config.Config
	}{
		{
			name: "valid config with two layers",
			yaml: `
layers:
  - name: domain
    path: internal/domain
    forbid:
      - internal/infra
      - internal/repository
  - name: handler
    path: internal/handler
    forbid:
      - internal/repository
`,
			want: &config.Config{
				Layers: []config.LayerConfig{
					{Name: "domain", Path: "internal/domain", Forbid: []string{"internal/infra", "internal/repository"}},
					{Name: "handler", Path: "internal/handler", Forbid: []string{"internal/repository"}},
				},
			},
		},
		{
			name:    "invalid yaml",
			yaml:    "layers: [invalid: : yaml",
			wantErr: true,
		},
		{
			name: "empty config",
			yaml: `layers: []`,
			want: &config.Config{Layers: []config.LayerConfig{}},
		},
		{
			name: "layer without forbid list",
			yaml: `
layers:
  - name: domain
    path: internal/domain
`,
			want: &config.Config{
				Layers: []config.LayerConfig{
					{Name: "domain", Path: "internal/domain", Forbid: nil},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			cfgPath := filepath.Join(dir, "clinicius.yaml")
			if err := os.WriteFile(cfgPath, []byte(tt.yaml), 0o644); err != nil {
				t.Fatalf("writing temp config: %v", err)
			}

			got, err := config.Load(cfgPath)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if len(got.Layers) != len(tt.want.Layers) {
				t.Fatalf("got %d layers, want %d", len(got.Layers), len(tt.want.Layers))
			}
			for i, layer := range got.Layers {
				want := tt.want.Layers[i]
				if layer.Name != want.Name {
					t.Errorf("layer[%d].Name = %q, want %q", i, layer.Name, want.Name)
				}
				if layer.Path != want.Path {
					t.Errorf("layer[%d].Path = %q, want %q", i, layer.Path, want.Path)
				}
				if len(layer.Forbid) != len(want.Forbid) {
					t.Errorf("layer[%d].Forbid = %v, want %v", i, layer.Forbid, want.Forbid)
					continue
				}
				for j, f := range layer.Forbid {
					if f != want.Forbid[j] {
						t.Errorf("layer[%d].Forbid[%d] = %q, want %q", i, j, f, want.Forbid[j])
					}
				}
			}
		})
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/clinicius.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
