package gotpl

import (
	"bytes"
	"embed"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

//go:embed testdata/templates/*
var testFS embed.FS

func TestValidateAndRender(t *testing.T) {
	tmpl := NewTemplate(testFS, WithTemplateRoot("testdata/templates"))

	if err := tmpl.Validate(); err != nil {
		t.Fatalf("Validate() error: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Render(&buf, "layout/home.html", "World"); err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	got := buf.String()
	want := "<html><header>Header</header>Hello, World!</html>"
	if got != want {
		t.Errorf("Render() =\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderMultipleLayouts(t *testing.T) {
	tmpl := NewTemplate(testFS, WithTemplateRoot("testdata/templates"))

	if err := tmpl.Validate(); err != nil {
		t.Fatalf("Validate() error: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Render(&buf, "app/dashboard.html", nil); err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	got := buf.String()
	want := "<app>Dashboard</app>"
	if got != want {
		t.Errorf("Render() =\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderBeforeValidate(t *testing.T) {
	tmpl := NewTemplate(testFS, WithTemplateRoot("testdata/templates"))

	var buf bytes.Buffer
	err := tmpl.Render(&buf, "layout/home.html", nil)
	if err != ErrTemplateNotFound {
		t.Errorf("Render() error = %v, want %v", err, ErrTemplateNotFound)
	}
}

func TestRenderNotFound(t *testing.T) {
	tmpl := NewTemplate(testFS, WithTemplateRoot("testdata/templates"))

	if err := tmpl.Validate(); err != nil {
		t.Fatalf("Validate() error: %v", err)
	}

	var buf bytes.Buffer
	err := tmpl.Render(&buf, "layout/nonexistent.html", nil)
	if err != ErrTemplateNotFound {
		t.Errorf("Render() error = %v, want %v", err, ErrTemplateNotFound)
	}
}

func TestValidateMissingRoot(t *testing.T) {
	tmpl := NewTemplate(testFS, WithTemplateRoot("nonexistent"))

	err := tmpl.Validate()
	if err == nil {
		t.Fatal("Validate() expected error for missing root, got nil")
	}
}

func TestWithTemplateRoot(t *testing.T) {
	tmpl := NewTemplate(testFS, WithTemplateRoot("testdata/templates"))

	if tmpl.config.templateRootName != "testdata/templates" {
		t.Errorf("templateRootName = %q, want %q", tmpl.config.templateRootName, "testdata/templates")
	}
}

func TestDefaultTemplateRoot(t *testing.T) {
	tmpl := NewTemplate(testFS)

	if tmpl.config.templateRootName != "templates" {
		t.Errorf("templateRootName = %q, want %q", tmpl.config.templateRootName, "templates")
	}
}

func TestFormFromRequest(t *testing.T) {
	body := url.Values{
		"email":   {"test@example.com"},
		"message": {"hello"},
	}

	r, err := http.NewRequest(http.MethodPost, "/contact", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tmpl := NewTemplate(testFS)
	form, err := tmpl.FormFromRequest(r)
	if err != nil {
		t.Fatal(err)
	}

	if got := form.Get("email"); got != "test@example.com" {
		t.Errorf("Get(email) = %q, want %q", got, "test@example.com")
	}
	if got := form.Get("message"); got != "hello" {
		t.Errorf("Get(message) = %q, want %q", got, "hello")
	}
	if !form.Valid() {
		t.Error("form from request should be valid (no errors added)")
	}
}

func TestFormFromRequestTakesFirstValue(t *testing.T) {
	body := url.Values{
		"color": {"red", "blue"},
	}

	r, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tmpl := NewTemplate(testFS)
	form, err := tmpl.FormFromRequest(r)
	if err != nil {
		t.Fatal(err)
	}

	if got := form.Get("color"); got != "red" {
		t.Errorf("Get(color) = %q, want %q (first value)", got, "red")
	}
}

func TestFormFromRequestIgnoresQueryParams(t *testing.T) {
	t.Run("query only", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodPost, "/?secret=fromquery", strings.NewReader(""))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		tmpl := NewTemplate(testFS)
		form, err := tmpl.FormFromRequest(r)
		if err != nil {
			t.Fatal(err)
		}

		if got := form.Get("secret"); got != "" {
			t.Errorf("Get(secret) = %q, want empty (query-only param should be ignored)", got)
		}
	})

	t.Run("body wins over query", func(t *testing.T) {
		body := url.Values{"secret": {"frombody"}}
		r, err := http.NewRequest(http.MethodPost, "/?secret=fromquery", strings.NewReader(body.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		tmpl := NewTemplate(testFS)
		form, err := tmpl.FormFromRequest(r)
		if err != nil {
			t.Fatal(err)
		}

		if got := form.Get("secret"); got != "frombody" {
			t.Errorf("Get(secret) = %q, want %q (body should win over query)", got, "frombody")
		}
	})
}

func TestFormFromRequestUsesCSRFGenerator(t *testing.T) {
	body := url.Values{
		"email":   {"test@example.com"},
		"message": {"hello"},
	}

	r, err := http.NewRequest(http.MethodPost, "/contact", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tmpl := NewTemplate(testFS, WithCSRF("CSRF", func(*http.Request) string { return "csrf-token-1234" }))
	form, err := tmpl.FormFromRequest(r)
	if err != nil {
		t.Fatal(err)
	}

	if form.csrfToken != "csrf-token-1234" {
		t.Errorf("Get(CSRF) = %q, want %q", form.csrfToken, "csrf-token-1234")
	}
}
