package view

import (
	"embed"
	"github.com/Masterminds/sprig/v3"
	"html/template"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

//go:embed templates/*
var templatesFS embed.FS

// ParseTemplates parses all HTML templates using the embedded file system.
func ParseTemplates() (*template.Template, error) {
	return template.New("").Funcs(sprig.FuncMap()).ParseFS(templatesFS, "templates/*.html")
}
