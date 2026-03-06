package discovery_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vtbarreto/CLinicius/internal/discovery"
)

func mkdirs(t *testing.T, root string, paths ...string) {
	t.Helper()
	for _, p := range paths {
		if err := os.MkdirAll(filepath.Join(root, p), 0o755); err != nil {
			t.Fatalf("mkdirs: %v", err)
		}
	}
}

func TestDiscover_NoLayersFound(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, root, "cmd", "pkg", "scripts")

	result, err := discovery.Discover(root)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	if len(result.Layers) != 0 {
		t.Errorf("expected 0 layers, got %d: %v", len(result.Layers), result.Layers)
	}
}

func TestDiscover_SingleLayer(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, root, "internal/handler")

	result, err := discovery.Discover(root)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	if len(result.Layers) != 1 {
		t.Fatalf("expected 1 layer, got %d", len(result.Layers))
	}
	if result.Layers[0].Type != discovery.LayerTypeHandler {
		t.Errorf("Type = %q, want %q", result.Layers[0].Type, discovery.LayerTypeHandler)
	}
}

func TestDiscover_HandlerForbidsRepositoryAndInfra(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, root,
		"internal/handler",
		"internal/repository",
		"internal/infra",
	)

	result, err := discovery.Discover(root)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	var handler *discovery.DiscoveredLayer
	for i := range result.Layers {
		if result.Layers[i].Type == discovery.LayerTypeHandler {
			handler = &result.Layers[i]
			break
		}
	}
	if handler == nil {
		t.Fatal("handler layer not found in result")
	}

	forbidSet := make(map[string]bool)
	for _, f := range handler.Forbid {
		forbidSet[f] = true
	}
	if !forbidSet["internal/repository"] {
		t.Error("handler should forbid internal/repository")
	}
	if !forbidSet["internal/infra"] {
		t.Error("handler should forbid internal/infra")
	}
}

func TestDiscover_DomainForbidsInfraRepositoryAndHandler(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, root,
		"internal/domain",
		"internal/handler",
		"internal/repository",
		"internal/infra",
	)

	result, err := discovery.Discover(root)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	var domain *discovery.DiscoveredLayer
	for i := range result.Layers {
		if result.Layers[i].Type == discovery.LayerTypeDomain {
			domain = &result.Layers[i]
			break
		}
	}
	if domain == nil {
		t.Fatal("domain layer not found in result")
	}

	forbidSet := make(map[string]bool)
	for _, f := range domain.Forbid {
		forbidSet[f] = true
	}
	if !forbidSet["internal/infra"] {
		t.Error("domain should forbid internal/infra")
	}
	if !forbidSet["internal/repository"] {
		t.Error("domain should forbid internal/repository")
	}
	if !forbidSet["internal/handler"] {
		t.Error("domain should forbid internal/handler")
	}
}

func TestDiscover_AlternativeNames(t *testing.T) {
	tests := []struct {
		dir      string
		wantType discovery.LayerType
	}{
		{"internal/controllers", discovery.LayerTypeHandler},
		{"internal/handlers", discovery.LayerTypeHandler},
		{"internal/services", discovery.LayerTypeUsecase},
		{"internal/repositories", discovery.LayerTypeRepository},
		{"internal/infrastructure", discovery.LayerTypeInfra},
	}

	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			root := t.TempDir()
			mkdirs(t, root, tt.dir)

			result, err := discovery.Discover(root)
			if err != nil {
				t.Fatalf("Discover() error = %v", err)
			}
			if len(result.Layers) != 1 {
				t.Fatalf("expected 1 layer, got %d", len(result.Layers))
			}
			if result.Layers[0].Type != tt.wantType {
				t.Errorf("Type = %q, want %q", result.Layers[0].Type, tt.wantType)
			}
		})
	}
}

func TestDiscover_SkipsVendorAndHidden(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, root,
		"vendor/internal/handler",
		".git/internal/repository",
		"internal/domain",
	)

	result, err := discovery.Discover(root)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	for _, l := range result.Layers {
		if l.Type == discovery.LayerTypeHandler || l.Type == discovery.LayerTypeRepository {
			t.Errorf("should have skipped vendor/.git, but found layer: %+v", l)
		}
	}
}

func TestDiscover_MultipleSameType(t *testing.T) {
	root := t.TempDir()
	// Two handler directories at different paths
	mkdirs(t, root,
		"internal/handler",
		"pkg/handlers",
		"internal/repository",
	)

	result, err := discovery.Discover(root)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	handlerCount := 0
	for _, l := range result.Layers {
		if l.Type == discovery.LayerTypeHandler {
			handlerCount++
			// Each handler should forbid the repository
			found := false
			for _, f := range l.Forbid {
				if f == "internal/repository" {
					found = true
				}
			}
			if !found {
				t.Errorf("handler at %q should forbid internal/repository", l.Path)
			}
		}
	}
	if handlerCount != 2 {
		t.Errorf("expected 2 handler layers, got %d", handlerCount)
	}
}
