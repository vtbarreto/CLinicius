package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/vtbarreto/CLinicius/internal/rules"
)

// JSONReporter serializes violations as a structured JSON object.
type JSONReporter struct {
	w io.Writer
}

// NewJSONReporter creates a JSONReporter that writes to w.
func NewJSONReporter(w io.Writer) *JSONReporter {
	return &JSONReporter{w: w}
}

type jsonReport struct {
	Violations []rules.Violation `json:"violations"`
	Total      int               `json:"total"`
}

// Report marshals violations to JSON and writes them to the underlying writer.
func (r *JSONReporter) Report(violations []rules.Violation) error {
	if violations == nil {
		violations = []rules.Violation{}
	}

	report := jsonReport{
		Violations: violations,
		Total:      len(violations),
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling violations: %w", err)
	}

	_, err = fmt.Fprintln(r.w, string(data))
	return err
}
