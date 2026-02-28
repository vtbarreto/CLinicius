package rules_test

import (
	"testing"

	"github.com/vtbarreto/CLinicius/internal/analyzer"
	"github.com/vtbarreto/CLinicius/internal/config"
	"github.com/vtbarreto/CLinicius/internal/rules"
)

func buildGraph(edges map[string][]string) *analyzer.DependencyGraph {
	g := analyzer.NewDependencyGraph()
	for from, imports := range edges {
		g.AddNode(from, nil)
		for _, to := range imports {
			g.AddEdge(from, to)
		}
	}
	return g
}

func TestLayerRule_NoViolations(t *testing.T) {
	g := buildGraph(map[string][]string{
		"myapp/internal/handler":    {"myapp/internal/domain"},
		"myapp/internal/domain":     {"myapp/internal/shared"},
		"myapp/internal/repository": {"myapp/internal/domain"},
	})
	cfg := &config.Config{
		Layers: []config.LayerConfig{
			{Name: "handler", Path: "internal/handler", Forbid: []string{"internal/repository"}},
		},
	}

	r := rules.NewLayerRule()
	violations := r.Validate(g, cfg)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestLayerRule_DetectsViolation(t *testing.T) {
	g := buildGraph(map[string][]string{
		"myapp/internal/handler": {
			"myapp/internal/domain",
			"myapp/internal/repository", // forbidden
		},
	})
	cfg := &config.Config{
		Layers: []config.LayerConfig{
			{Name: "handler", Path: "internal/handler", Forbid: []string{"internal/repository"}},
		},
	}

	r := rules.NewLayerRule()
	violations := r.Validate(g, cfg)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(violations), violations)
	}
	if violations[0].Layer != "handler" {
		t.Errorf("violation.Layer = %q, want %q", violations[0].Layer, "handler")
	}
	if violations[0].Rule != "layer-boundary" {
		t.Errorf("violation.Rule = %q, want %q", violations[0].Rule, "layer-boundary")
	}
}

func TestLayerRule_MultipleForbidden(t *testing.T) {
	g := buildGraph(map[string][]string{
		"myapp/internal/domain": {
			"myapp/internal/infra",
			"myapp/internal/repository",
		},
	})
	cfg := &config.Config{
		Layers: []config.LayerConfig{
			{Name: "domain", Path: "internal/domain", Forbid: []string{"internal/infra", "internal/repository"}},
		},
	}

	r := rules.NewLayerRule()
	violations := r.Validate(g, cfg)
	if len(violations) != 2 {
		t.Errorf("expected 2 violations, got %d", len(violations))
	}
}

func TestLayerRule_EmptyConfig(t *testing.T) {
	g := buildGraph(map[string][]string{
		"myapp/internal/handler": {"myapp/internal/repository"},
	})
	cfg := &config.Config{}

	r := rules.NewLayerRule()
	violations := r.Validate(g, cfg)
	if len(violations) != 0 {
		t.Errorf("expected no violations with empty config, got %v", violations)
	}
}
