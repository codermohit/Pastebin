package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"capybara.pastebin.xyz/internal/models"
	"capybara.pastebin.xyz/ui"
)

type templateData struct {
	Paste           models.Paste
	Pastes          []models.Paste
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
  if t.IsZero() {
    return ""
  }

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page) //will return the name of the html files like "home.html.tmpl"

    patterns := []string{
      "html/base.tmpl.html",
      "html/partials/*.html",
      page,
    }

    //to parse the template files from the ui.Files embedded filesystem
    ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
    if err != nil {
      return nil , err
    }

		cache[name] = ts
	}

	return cache, nil
}
