// Package gotpl provides a layout-based HTML template renderer backed by [embed.FS].
//
// Templates are organized into layouts, views, and partials:
//
//	templates/
//	  layout.html          # layout files in the root
//	  app.html
//	  partials/            # shared partials included in every view
//	    header.html
//	  views/
//	    layout/            # views for "layout.html"
//	      home.html
//	    app/               # views for "app.html"
//	      dashboard.html
//
// Call [Template.Validate] to parse all templates, then [Template.Render] to
// execute a view by its "[layout]/[page.html]" name (e.g. "app/dashboard.html").
package gotpl

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"path"
	"path/filepath"
	"strings"
)

// PageData is a convenience wrapper for passing a page title and
// arbitrary data to a template.
type PageData struct {
	Title string
	Data  any
}

// Template holds parsed HTML templates organized by layout and view.
type Template struct {
	fs embed.FS

	config *option

	views map[string]*template.Template
}

// NewTemplate creates a new [Template] from the given [embed.FS].
// Use [OptionFunc] values such as [WithTemplateRoot] to customize behaviour.
// Call [Template.Validate] before rendering.
func NewTemplate(fs embed.FS, opts ...OptionFunc) *Template {
	config := defaultOption()
	for _, opt := range opts {
		opt(config)
	}

	return &Template{
		fs:     fs,
		config: config,

		views: make(map[string]*template.Template),
	}
}

// Render executes the named view template and writes the result to w.
//
// The view name follows the pattern "[layout]/[page.html]", where layout is the
// layout filename without its extension. For example, given layouts "layout.html"
// and "app.html", a view "dashboard.html" under the "app" layout is rendered as:
//
//	templ.Render(w, "app/dashboard.html", data)
func (templ *Template) Render(w io.Writer, view string, data any) error {
	v, ok := templ.views[view]
	if !ok {
		return ErrTemplateNotFound
	}

	return v.Execute(w, data)
}

// Validate parses all templates contained in the [embed.FS].
//
// It must be called before [Template.Render]; rendering without prior
// validation will always return [ErrTemplateNotFound].
func (t *Template) Validate() error {
	// partials are optional, ignore error if directory does not exist
	partials, _ := readDir(t.fs, t.config.templateRootName, "partials")

	layouts, err := readDir(t.fs, t.config.templateRootName)
	if err != nil {
		return err
	}

	viewsDir := path.Join(t.config.templateRootName, "views")
	views := make(map[string]*template.Template)
	for _, layout := range layouts {
		layoutView := strings.TrimSuffix(layout.name, filepath.Ext(layout.name))

		pages, err := readDir(t.fs, viewsDir, layoutView)
		if err != nil {
			return err
		}

		for _, view := range pages {
			viewName := fmt.Sprintf(layoutView+"/%s", view.name)

			patterns := []string{
				layout.path,
				view.path,
			}
			patterns = append(patterns, getPaths(partials)...)

			tmpl := template.New(layout.name)
			parsed, err := tmpl.ParseFS(t.fs, patterns...)
			if err != nil {
				return err
			}

			views[viewName] = parsed
		}
	}

	t.views = views

	return nil
}
