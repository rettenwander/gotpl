package gotpl

import "net/http"

// Form holds submitted field values and validation errors for template rendering.
//
// Handlers populate field errors and form-level errors after parsing a request;
// templates read them back to display inline messages and re-populate inputs.
//
//	form := gotpl.FormFromRequest(r)
//
//	if form.Get("email") == "" {
//	    form.AddFieldError("email", "Email is required")
//	}
//
//	if !form.Valid() {
//	    tmpl.Render(w, "layout/contact.html", gotpl.PageData{Form: form})
//	    return
//	}
type Form struct {
	// Values holds the raw string values keyed by field name.
	Values map[string]string

	// FieldErrors holds a single validation error per field, keyed by field name.
	// In templates: {{with .Form.FieldErrors.email}}<p class="error">{{.}}</p>{{end}}
	FieldErrors map[string]string

	// Errors holds form-level errors not tied to a specific field
	// (e.g. "invalid credentials").
	Errors []string
}

// NewForm returns an empty [Form] ready for use.
func NewForm() *Form {
	return &Form{
		Values:      make(map[string]string),
		FieldErrors: make(map[string]string),
	}
}

// FormFromRequest creates a [Form] pre-populated with all POST body values
// from the given request. It takes the first value for each field.
func FormFromRequest(r *http.Request) *Form {
	r.ParseForm()

	form := NewForm()
	for field, values := range r.PostForm {
		if len(values) > 0 {
			form.Values[field] = values[0]
		}
	}
	return form
}

// Set stores a field value.
func (f *Form) Set(field, value string) {
	f.Values[field] = value
}

// Get returns the value of a field, or "" if unset.
func (f *Form) Get(field string) string {
	return f.Values[field]
}

// AddFieldError adds a validation error for a specific field.
// Only the first error per field is kept.
func (f *Form) AddFieldError(field, message string) {
	if _, exists := f.FieldErrors[field]; !exists {
		f.FieldErrors[field] = message
	}
}

// AddError adds a form-level error not tied to a specific field.
func (f *Form) AddError(message string) {
	f.Errors = append(f.Errors, message)
}

// Valid reports whether the form has no errors of any kind.
func (f *Form) Valid() bool {
	return len(f.FieldErrors) == 0 && len(f.Errors) == 0
}
