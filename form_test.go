package gotpl

import (
	"testing"
)

func TestNewForm(t *testing.T) {
	form := NewForm("", "")

	if form.Values == nil {
		t.Fatal("Values map is nil")
	}
	if form.FieldErrors == nil {
		t.Fatal("FieldErrors map is nil")
	}
	if !form.Valid() {
		t.Error("new form should be valid")
	}
}

func TestFormSetGet(t *testing.T) {
	form := NewForm("", "")
	form.Set("email", "test@example.com")

	if got := form.Get("email"); got != "test@example.com" {
		t.Errorf("Get(email) = %q, want %q", got, "test@example.com")
	}
}

func TestFormGetMissing(t *testing.T) {
	form := NewForm("", "")

	if got := form.Get("missing"); got != "" {
		t.Errorf("Get(missing) = %q, want empty string", got)
	}
}

func TestFormAddFieldError(t *testing.T) {
	form := NewForm("", "")
	form.AddFieldError("email", "Email is required")

	if got := form.FieldErrors["email"]; got != "Email is required" {
		t.Errorf("FieldErrors[email] = %q, want %q", got, "Email is required")
	}
	if form.Valid() {
		t.Error("form with field error should not be valid")
	}
}

func TestFormAddFieldErrorKeepsFirst(t *testing.T) {
	form := NewForm("", "")
	form.AddFieldError("email", "first error")
	form.AddFieldError("email", "second error")

	if got := form.FieldErrors["email"]; got != "first error" {
		t.Errorf("FieldErrors[email] = %q, want %q (first error)", got, "first error")
	}
}

func TestFormAddError(t *testing.T) {
	form := NewForm("", "")
	form.AddError("invalid credentials")

	if len(form.Errors) != 1 || form.Errors[0] != "invalid credentials" {
		t.Errorf("Errors = %v, want [invalid credentials]", form.Errors)
	}
	if form.Valid() {
		t.Error("form with error should not be valid")
	}
}

func TestFormValid(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		form := NewForm("", "")
		form.Set("name", "Alice")
		if !form.Valid() {
			t.Error("form without errors should be valid")
		}
	})

	t.Run("field error only", func(t *testing.T) {
		form := NewForm("", "")
		form.AddFieldError("name", "required")
		if form.Valid() {
			t.Error("form with field error should not be valid")
		}
	})

	t.Run("form error only", func(t *testing.T) {
		form := NewForm("", "")
		form.AddError("something went wrong")
		if form.Valid() {
			t.Error("form with form error should not be valid")
		}
	})

	t.Run("both errors", func(t *testing.T) {
		form := NewForm("", "")
		form.AddFieldError("name", "required")
		form.AddError("something went wrong")
		if form.Valid() {
			t.Error("form with both errors should not be valid")
		}
	})
}

func TestCsrfField(t *testing.T) {
	form := NewForm("CSRF", "csrf-token-1234")
	field := form.CSRF()

	if field != `<input type="hidden" name="CSRF" value="csrf-token-1234">` {
		t.Errorf("CSRF field is not valid")
	}
}
