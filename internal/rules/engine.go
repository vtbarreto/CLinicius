package rules

import (
	"github.com/vtbarreto/CLinicius/internal/analyzer"
	"github.com/vtbarreto/CLinicius/internal/config"
)

// Violation describes a single detected architectural rule breach.
type Violation struct {
	Rule     string `json:"rule"`
	Layer    string `json:"layer,omitempty"`
	File     string `json:"file,omitempty"`
	Importer string `json:"importer,omitempty"`
	Imported string `json:"imported,omitempty"`
	Message  string `json:"message"`
}

// Rule is the interface that all architectural rules must implement.
type Rule interface {
	Name() string
	Validate(graph *analyzer.DependencyGraph, cfg *config.Config) []Violation
}

// Engine runs a set of Rules against a DependencyGraph and aggregates
// the resulting Violations.
type Engine struct {
	rules []Rule
}

// NewEngine creates an Engine pre-loaded with the given rules.
func NewEngine(rules ...Rule) *Engine {
	return &Engine{rules: rules}
}

// Run executes every registered rule and returns all violations found.
func (e *Engine) Run(graph *analyzer.DependencyGraph, cfg *config.Config) []Violation {
	var violations []Violation
	for _, r := range e.rules {
		violations = append(violations, r.Validate(graph, cfg)...)
	}
	return violations
}
