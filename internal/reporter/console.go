package reporter

import (
	"fmt"
	"io"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/vtbarreto/CLinicius/internal/i18n"
	"github.com/vtbarreto/CLinicius/internal/rules"
)

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorBold  = "\033[1m"
	colorReset = "\033[0m"
)

// ConsoleReporter prints violations to a writer using human-readable,
// color-highlighted output in the configured language.
type ConsoleReporter struct {
	w   io.Writer
	loc *goi18n.Localizer
}

// NewConsoleReporter creates a ConsoleReporter. If localizer is nil,
// it defaults to English.
func NewConsoleReporter(w io.Writer, localizer *goi18n.Localizer) *ConsoleReporter {
	if localizer == nil {
		localizer = i18n.NewLocalizer("en-US")
	}
	return &ConsoleReporter{w: w, loc: localizer}
}

// Report writes all violations to the writer. When no violations exist,
// a success message is printed instead.
func (r *ConsoleReporter) Report(violations []rules.Violation) {
	t := func(id string) string { return i18n.T(r.loc, id) }

	if len(violations) == 0 {
		fmt.Fprintf(r.w, "%s%s%s\n", colorGreen, t("NoViolations"), colorReset)
		return
	}

	for _, v := range violations {
		fmt.Fprintf(r.w, "\n%s%s❌ %s%s\n", colorBold, colorRed, t("ViolationHeader"), colorReset)
		if v.Layer != "" {
			fmt.Fprintf(r.w, "  %-8s %s\n", t("LabelLayer")+":", v.Layer)
		}
		if v.Importer != "" {
			fmt.Fprintf(r.w, "  %-8s %s\n", t("LabelPackage")+":", v.Importer)
		}
		if v.File != "" {
			fmt.Fprintf(r.w, "  %-8s %s\n", t("LabelFile")+":", v.File)
		}
		if v.Imported != "" {
			fmt.Fprintf(r.w, "  %-8s %s\n", t("LabelImports")+":", v.Imported)
		}
		fmt.Fprintf(r.w, "  %-8s %s\n", t("LabelRule")+":", v.Rule)
		fmt.Fprintf(r.w, "  %-8s %s\n", t("LabelDetail")+":", v.Message)
	}

	fmt.Fprintf(r.w, "\n%s\n", i18n.TPlural(r.loc, "ViolationCount", len(violations)))
}
