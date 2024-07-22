package main

import (
	"html/template"
	"path/filepath"

	"capybara.pastebin.xyz/internal/models"
)

type templateData struct {
	Paste  models.Paste
	Pastes []models.Paste
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page) //will return the name of the html files like home.html.tmpl
		files := []string{
			"./ui/html/base.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
      page,
		}

		//parse the template files into a template set
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
