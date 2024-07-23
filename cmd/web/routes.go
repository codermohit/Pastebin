package main

import "net/http"

func (app *application) routes() http.Handler{
	mux := http.NewServeMux()

	//file server to serve files from the "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	//Create the new route, which is restricted to POST requests only
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

  return commonHeaders(mux)
}
