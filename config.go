package gotpl

// OptionFunc is a functional option for configuring a [Template].
type OptionFunc = func(c *option)

type option struct {
	templateRootName string
}

func defaultOption() *option {
	return &option{
		templateRootName: "templates",
	}
}

// WithTemplateRoot sets the root directory inside the [embed.FS] that contains
// the layout files, views/, and partials/ directories. The default is "templates".
func WithTemplateRoot(path string) OptionFunc {
	return func(c *option) {
		c.templateRootName = path
	}
}
