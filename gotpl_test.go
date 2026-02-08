package gotpl

import (
	"bytes"
	"embed"
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
