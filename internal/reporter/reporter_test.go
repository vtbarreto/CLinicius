package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/vtbarreto/CLinicius/internal/reporter"
	"github.com/vtbarreto/CLinicius/internal/rules"
)

// nil localizer → defaults to English inside NewConsoleReporter

// ---- ConsoleReporter -------------------------------------------------------

func TestConsoleReporter_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewConsoleReporter(&buf, nil)
	r.Report(nil)

	out := buf.String()
	if !strings.Contains(out, "No architectural violations") {
		t.Errorf("expected success message, got: %q", out)
	}
}

func TestConsoleReporter_WithViolations(t *testing.T) {
	violations := []rules.Violation{
		{
			Rule:     "layer-boundary",
			Layer:    "handler",
			Importer: "myapp/internal/handler",
			Imported: "myapp/internal/repository",
			Message:  "handler layer cannot depend on internal/repository",
		},
	}

	var buf bytes.Buffer
	r := reporter.NewConsoleReporter(&buf, nil)
	r.Report(violations)

	out := buf.String()
	if !strings.Contains(out, "Architectural Violation") {
		t.Errorf("expected violation header, got: %q", out)
	}
	if !strings.Contains(out, "handler") {
		t.Errorf("expected layer name in output, got: %q", out)
	}
	if !strings.Contains(out, "violation") {
		t.Errorf("expected violation count, got: %q", out)
	}
}

func TestConsoleReporter_CountMatchesInput(t *testing.T) {
	violations := []rules.Violation{
		{Rule: "cycle-detection", Message: "cyclic dependency: a → b → a"},
		{Rule: "layer-boundary", Layer: "domain", Message: "domain cannot depend on infra"},
	}

	var buf bytes.Buffer
	reporter.NewConsoleReporter(&buf, nil).Report(violations)

	if !strings.Contains(buf.String(), "2 violations found.") {
		t.Errorf("expected '2 violation(s)', got: %q", buf.String())
	}
}

// ---- JSONReporter ----------------------------------------------------------

func TestJSONReporter_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewJSONReporter(&buf)
	if err := r.Report(nil); err != nil {
		t.Fatalf("Report() error = %v", err)
	}

	var report struct {
		Total      int              `json:"total"`
		Violations []rules.Violation `json:"violations"`
	}
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON output: %v\nraw: %s", err, buf.String())
	}
	if report.Total != 0 {
		t.Errorf("total = %d, want 0", report.Total)
	}
	if len(report.Violations) != 0 {
		t.Errorf("violations = %v, want empty", report.Violations)
	}
}

func TestJSONReporter_WithViolations(t *testing.T) {
	violations := []rules.Violation{
		{Rule: "layer-boundary", Layer: "handler", Importer: "myapp/handler", Message: "violation msg"},
	}

	var buf bytes.Buffer
	r := reporter.NewJSONReporter(&buf)
	if err := r.Report(violations); err != nil {
		t.Fatalf("Report() error = %v", err)
	}

	var report struct {
		Total      int              `json:"total"`
		Violations []rules.Violation `json:"violations"`
	}
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	if report.Total != 1 {
		t.Errorf("total = %d, want 1", report.Total)
	}
	if report.Violations[0].Layer != "handler" {
		t.Errorf("violations[0].Layer = %q, want %q", report.Violations[0].Layer, "handler")
	}
}
