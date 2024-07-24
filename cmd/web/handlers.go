package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"capybara.pastebin.xyz/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	pastes, err := app.pastes.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Pastes = pastes

  
	//use the new render helper
	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

// display a specific snippet
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	paste, err := app.pastes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Paste = paste

	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

// display a form for creating a new snippet
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r)

  app.render(w, r , http.StatusOK, "create.tmpl.html", data)
}

// create a new snippet post
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
  r.Body = http.MaxBytesReader(w, r.Body, 4096)

  err := r.ParseForm()
  if err != nil{
    app.clientError(w, http.StatusBadRequest)
    return
  }

  title := r.PostForm.Get("title")
  content := r.PostForm.Get("content")
  expires, err := strconv.Atoi(r.PostForm.Get("expires"))
  
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

	id, err := app.pastes.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
    return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
