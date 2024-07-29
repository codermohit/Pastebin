package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"golang.org/x/net/context"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//set Connection close header on response
				w.Header().Set("Connection", "close")

				//send Internal server error
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if !app.isAuthenticated(r){
      http.Redirect(w, r, "/user/login", http.StatusSeeOther)
      return
    }

    w.Header().Add("Cache-Control", "no-store")
    next.ServeHTTP(w, r)
  })
}

func noSurf(next http.Handler) http.Handler{
  csrfHandler := nosurf.New(next)
  csrfHandler.SetBaseCookie(http.Cookie{
    HttpOnly: true,
    Path: "/",
    Secure: true,
  })

  return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    /*
    Check if the client request has the authenticatedUserID in the session cookie 
    1) If not, let it pass, will be taken care of by the requireAuthentication() middleware
    in case it tries to access authorized resources
    2) If authenticatedUserID exists in the session cookie, check for the authenticatedUserID in the 
    database , if it doesn't exist in the DB throw an error else call the next handler 
    */
    id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
    if id == 0 {
      next.ServeHTTP(w, r)
      return
    }

    //check if a user with 'id' exists
    exists, err := app.users.Exists(id)
    if err != nil {
      app.serverError(w, r , err)
      return
    }

    if exists {
      ctx := context.WithValue(r.Context(),isAuthenticatedContextKey,true)
      r = r.WithContext(ctx)
    }

    next.ServeHTTP(w, r)
  })
}
