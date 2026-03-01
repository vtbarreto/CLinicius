package rules

import (
	"fmt"
	"strings"

	"github.com/vtbarreto/CLinicius/internal/analyzer"
	"github.com/vtbarreto/CLinicius/internal/config"
)

// CycleRule detects cyclic import dependencies between packages in the graph.
type CycleRule struct{}

// NewCycleRule creates a new CycleRule.
func NewCycleRule() *CycleRule {
	return &CycleRule{}
}

func (r *CycleRule) Name() string {
	return "cycle-detection"
}

// Validate finds all cycles in the dependency graph and reports each one
// as a Violation with a human-readable path.
func (r *CycleRule) Validate(graph *analyzer.DependencyGraph, _ *config.Config) []Violation {
	var violations []Violation

	for _, cycle := range graph.FindCycles() {
		violations = append(violations, Violation{
			Rule:    r.Name(),
			Message: fmt.Sprintf("cyclic dependency detected: %s", strings.Join(cycle, " → ")),
		})
	}

	return violations
}
