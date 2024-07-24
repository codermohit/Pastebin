package main

import (
  "net/http"
  "github.com/justinas/alice"
)

func (app *application) routes() http.Handler{
	mux := http.NewServeMux()

	//file server to serve files from the "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

  dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	//Create the new route, which is restricted to POST requests only
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

  standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders) 

  return standard.Then(mux) 
}
