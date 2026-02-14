package gotpl

import "net/http"

// OptionFunc is a functional option for configuring a [Template].
type OptionFunc = func(c *option)

type csrfTokenGenerator = func(*http.Request) string

type option struct {
	templateRootName string
	csrfFieldName    string
	csrfTokenGenerator
}

func defaultOption() *option {
	return &option{
		templateRootName: "templates",
		csrfFieldName:    "csrf",
	}
}

// WithTemplateRoot sets the root directory inside the [embed.FS] that contains
// the layout files, views/, and partials/ directories. The default is "templates".
func WithTemplateRoot(path string) OptionFunc {
	return func(c *option) {
		c.templateRootName = path
	}
}

// WithCSRF configures the CSRF field name and token generator used by [Template.RequestForm].
func WithCSRF(fieldName string, generator csrfTokenGenerator) OptionFunc {
	return func(c *option) {
		c.csrfFieldName = fieldName
		c.csrfTokenGenerator = generator
	}
}
