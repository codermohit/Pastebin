package main

import (
	"fmt"
	"net/http"
)


//logs method and attribute of the request sent by the user, and sends generic 500 Internal server error response
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error){
  var (
    method = r.Method
    uri = r.URL.RequestURI()
  )

  app.logger.Error(err.Error(), "method", method, "uri", uri)
  http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}


//used when the problem is with the request sent by the user 
func (app *application) clientError(w http.ResponseWriter, status int){
  http.Error(w, http.StatusText(status), status)  
}


//helper method to render templates from the cache
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData){
  ts, ok := app.templateCache[page]

  if !ok {
    err := fmt.Errorf("the template %s doesn't exist", page)
    app.serverError(w, r, err)
    return
  }

  w.WriteHeader(status)
  
  err := ts.ExecuteTemplate(w, "base", data)
  if err != nil {
    app.serverError(w, r , err)
  }

}
