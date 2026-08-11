// Package handlers implements the Alfred entry point (the 우리말샘/opendict
// autocomplete search) and dispatches to it by subcommand name.
package handlers

import (
	"fmt"

	"github.com/inchanS/AlfKoreanSearch/internal/alfred"
	"github.com/inchanS/AlfKoreanSearch/internal/update"
	"github.com/inchanS/AlfKoreanSearch/internal/urlx"
)

// Icon file names bundled in the workflow directory.
const iconNoResults = "noresults.png"

// Dispatch runs the handler named cmd with the remaining CLI args.
func Dispatch(cmd string, args []string) {
	switch cmd {
	case "search":
		run(func(fb *alfred.Feedback) error { return search(fb, at(args, 0)) })
	default:
		fb := alfred.New()
		fb.Add(alfred.ItemOpts{Title: "Unknown command: " + cmd, Valid: false})
		fb.Send()
	}
}

// run wraps a handler body with the update notice, panic/error recovery, and
// final feedback emission. On error it shows a single error row, mirroring the
// old wf.run behaviour.
func run(body func(fb *alfred.Feedback) error) {
	fb := alfred.New()
	update.MaybeShow(fb)
	if err := safe(fb, body); err != nil {
		fb.Reset()
		fb.Add(alfred.ItemOpts{
			Title:    "Error",
			Subtitle: err.Error(),
			Icon:     iconNoResults,
			Valid:    false,
		})
	}
	fb.Send()
}

// safe invokes body, converting a panic into an error.
func safe(fb *alfred.Feedback, body func(*alfred.Feedback) error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	return body(fb)
}

// at returns args[i] or "" when out of range.
func at(args []string, i int) string {
	if i >= 0 && i < len(args) {
		return args[i]
	}
	return ""
}

// quick renders a query into a URL template of the form "…query=%s".
func quick(tmpl, word string) string {
	return fmt.Sprintf(tmpl, urlx.Quote(word))
}
