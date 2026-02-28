package rules

import (
	"fmt"
	"strings"

	"github.com/vtbarreto/CLinicius/internal/analyzer"
	"github.com/vtbarreto/CLinicius/internal/config"
)

// LayerRule enforces the layer boundary constraints declared in the config.
// For every package whose import path contains a layer's path prefix,
// it checks that none of its imports contain a forbidden path.
type LayerRule struct{}

// NewLayerRule creates a new LayerRule.
func NewLayerRule() *LayerRule {
	return &LayerRule{}
}

func (r *LayerRule) Name() string {
	return "layer-boundary"
}

// Validate iterates over all graph nodes and reports any import that
// violates the layer dependency constraints declared in the config.
func (r *LayerRule) Validate(graph *analyzer.DependencyGraph, cfg *config.Config) []Violation {
	var violations []Violation

	for pkgPath, node := range graph.Nodes() {
		for _, layer := range cfg.Layers {
			if !strings.Contains(pkgPath, layer.Path) {
				continue
			}
			for _, imp := range node.Imports {
				for _, forbidden := range layer.Forbid {
					if strings.Contains(imp, forbidden) {
						violations = append(violations, Violation{
							Rule:     r.Name(),
							Layer:    layer.Name,
							Importer: pkgPath,
							Imported: imp,
							Message:  fmt.Sprintf("%s layer cannot depend on %s", layer.Name, forbidden),
						})
					}
				}
			}
		}
	}

	return violations
}
