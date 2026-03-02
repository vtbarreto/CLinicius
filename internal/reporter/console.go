package reporter

import (
	"fmt"
	"io"

	"github.com/vtbarreto/CLinicius/internal/rules"
)

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorBold  = "\033[1m"
	colorReset = "\033[0m"
)

// ConsoleReporter prints violations to a writer using human-readable,
// color-highlighted output.
type ConsoleReporter struct {
	w io.Writer
}

// NewConsoleReporter creates a ConsoleReporter that writes to w.
func NewConsoleReporter(w io.Writer) *ConsoleReporter {
	return &ConsoleReporter{w: w}
}

// Report writes all violations to the writer. When no violations exist,
// a success message is printed instead.
func (r *ConsoleReporter) Report(violations []rules.Violation) {
	if len(violations) == 0 {
		fmt.Fprintf(r.w, "%s✅ No architectural violations found.%s\n", colorGreen, colorReset)
		return
	}

	for _, v := range violations {
		fmt.Fprintf(r.w, "\n%s%s❌ Architectural Violation%s\n", colorBold, colorRed, colorReset)
		if v.Layer != "" {
			fmt.Fprintf(r.w, "  Layer:   %s\n", v.Layer)
		}
		if v.Importer != "" {
			fmt.Fprintf(r.w, "  Package: %s\n", v.Importer)
		}
		if v.File != "" {
			fmt.Fprintf(r.w, "  File:    %s\n", v.File)
		}
		if v.Imported != "" {
			fmt.Fprintf(r.w, "  Imports: %s\n", v.Imported)
		}
		fmt.Fprintf(r.w, "  Rule:    %s\n", v.Rule)
		fmt.Fprintf(r.w, "  Detail:  %s\n", v.Message)
	}

	fmt.Fprintf(r.w, "\n%d violation(s) found.\n", len(violations))
}
