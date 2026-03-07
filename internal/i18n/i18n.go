package i18n

import (
	"embed"
	"encoding/json"
	"os"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	entries, _ := localeFS.ReadDir("locales")
	for _, entry := range entries {
		// errors loading a locale file are silently ignored so the
		// binary never panics — it just falls back to English.
		_, _ = bundle.LoadMessageFileFS(localeFS, "locales/"+entry.Name())
	}
}

// NewLocalizer creates a Localizer for the given language tag.
// If lang is empty, the language is detected from the environment
// following this priority:
//
//  1. CLINICIUS_LANG env var
//  2. LANG env var  (e.g. pt_BR.UTF-8 → pt-BR)
//  3. LANGUAGE env var
//  4. Fallback: en-US
func NewLocalizer(lang string) *i18n.Localizer {
	if lang == "" {
		lang = detectLanguage()
	}
	return i18n.NewLocalizer(bundle, lang, "en-US")
}

// T is a convenience wrapper that translates a message ID using the
// provided localizer, falling back to the ID itself if not found.
func T(loc *i18n.Localizer, id string) string {
	msg, err := loc.Localize(&i18n.LocalizeConfig{MessageID: id})
	if err != nil {
		return id
	}
	return msg
}

// TPlural translates a plural-aware message with a numeric count.
func TPlural(loc *i18n.Localizer, id string, count int) string {
	msg, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:   id,
		PluralCount: count,
		TemplateData: map[string]int{
			"Count": count,
		},
	})
	if err != nil {
		return id
	}
	return msg
}

func detectLanguage() string {
	for _, env := range []string{"CLINICIUS_LANG", "LANG", "LANGUAGE"} {
		if val := os.Getenv(env); val != "" {
			return normalizeLang(val)
		}
	}
	return "en-US"
}

// normalizeLang converts OS locale formats to BCP 47 language tags.
// Examples: "pt_BR.UTF-8" → "pt-BR", "en_US" → "en-US"
func normalizeLang(lang string) string {
	lang = strings.SplitN(lang, ".", 2)[0]
	lang = strings.ReplaceAll(lang, "_", "-")
	return lang
}
