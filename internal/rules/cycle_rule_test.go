package rules_test

import (
	"testing"

	"github.com/vtbarreto/CLinicius/internal/config"
	"github.com/vtbarreto/CLinicius/internal/rules"
)

func TestCycleRule_NoCycle(t *testing.T) {
	g := buildGraph(map[string][]string{
		"a": {"b"},
		"b": {"c"},
	})

	r := rules.NewCycleRule()
	violations := r.Validate(g, &config.Config{})
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestCycleRule_DirectCycle(t *testing.T) {
	g := buildGraph(map[string][]string{
		"a": {"b"},
		"b": {"a"},
	})

	r := rules.NewCycleRule()
	violations := r.Validate(g, &config.Config{})
	if len(violations) == 0 {
		t.Fatal("expected at least one violation for a↔b cycle")
	}
	if violations[0].Rule != "cycle-detection" {
		t.Errorf("Rule = %q, want %q", violations[0].Rule, "cycle-detection")
	}
}

func TestCycleRule_LongCycle(t *testing.T) {
	g := buildGraph(map[string][]string{
		"a": {"b"},
		"b": {"c"},
		"c": {"d"},
		"d": {"a"},
	})

	r := rules.NewCycleRule()
	violations := r.Validate(g, &config.Config{})
	if len(violations) == 0 {
		t.Fatal("expected violation for a→b→c→d→a cycle")
	}
}

func TestCycleRule_SelfLoop(t *testing.T) {
	g := buildGraph(map[string][]string{
		"a": {"a"},
	})

	r := rules.NewCycleRule()
	violations := r.Validate(g, &config.Config{})
	if len(violations) == 0 {
		t.Fatal("expected violation for self-loop a→a")
	}
}

func TestCycleRule_MessageContainsCyclePath(t *testing.T) {
	g := buildGraph(map[string][]string{
		"pkg/x": {"pkg/y"},
		"pkg/y": {"pkg/x"},
	})

	r := rules.NewCycleRule()
	violations := r.Validate(g, &config.Config{})
	if len(violations) == 0 {
		t.Fatal("expected violation")
	}
	msg := violations[0].Message
	if msg == "" {
		t.Error("violation message should not be empty")
	}
}
