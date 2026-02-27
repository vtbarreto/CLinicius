package analyzer

import (
	"fmt"

	"golang.org/x/tools/go/packages"
)

// LoadPackages loads Go packages matching the given patterns using
// module-aware resolution and builds a DependencyGraph from their imports.
//
// Only packages returned directly by the patterns are added as graph nodes.
// Their import paths are recorded as edges, enabling layer and cycle analysis.
func LoadPackages(patterns []string) (*DependencyGraph, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("loading packages: %w", err)
	}

	graph := NewDependencyGraph()

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			continue
		}
		node := graph.AddNode(pkg.PkgPath, pkg.GoFiles)
		for importPath := range pkg.Imports {
			node.addImport(importPath)
		}
	}

	return graph, nil
}
