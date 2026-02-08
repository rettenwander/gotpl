package gotpl

import "errors"

// ErrTemplateNotFound is returned by [Template.Render] when the requested view
// has not been parsed by [Template.Validate].
var (
	ErrTemplateNotFound = errors.New("template not found")
)
