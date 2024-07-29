package main

import (
	"html/template"
	"path/filepath"
	"time"

	"capybara.pastebin.xyz/internal/models"
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
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page) //will return the name of the html files like "home.html.tmpl"

		//parse the base template file into a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		//add partials to the template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		//add the page template to the template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
