# gotpl

A layout-based HTML template renderer for Go, backed by `embed.FS`.

## Install

```sh
go get github.com/rettenwander/gotpl
```

## Directory Structure

gotpl expects templates to be organized as follows:

```
templates/
  layout.html              # layout files in the root
  app.html
  partials/                # shared partials included in every view (optional)
    header.html
    footer.html
  views/
    layout/                # views for "layout.html"
      home.html
    app/                   # views for "app.html"
      dashboard.html
```

- **Layouts** are `.html` files in the template root. Each layout is a full page skeleton that references named templates (e.g. `{{template "content" .}}`).
- **Views** live under `views/<layout>/` where `<layout>` matches the layout filename without its extension. Each view defines the named templates the layout expects.
- **Partials** are optional shared snippets in `partials/`. They are included in every view and can be referenced by any layout or view.

## Usage

```go
package main

import (
	"embed"
	"log"
	"os"

	"github.com/rettenwander/gotpl"
)

//go:embed templates/*
var templateFS embed.FS

func main() {
	tmpl := gotpl.NewTemplate(templateFS)

	if err := tmpl.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := tmpl.Render(os.Stdout, "app/dashboard.html", nil); err != nil {
		log.Fatal(err)
	}
}
```

## Options

### WithTemplateRoot

Override the root directory inside the `embed.FS`. The default is `"templates"`.

```go
tmpl := gotpl.NewTemplate(fs, gotpl.WithTemplateRoot("web/templates"))
```

## Errors

- `gotpl.ErrTemplateNotFound` is returned by `Render` when the requested view has not been parsed by `Validate`.

## Testing

```sh
go test ./...
```
