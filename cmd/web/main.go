package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"capybara.pastebin.xyz/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger         *slog.Logger
	pastes         *models.PasteModel
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	//data source name : connection string that describes how to connect to your database
	dsn := flag.String("dsn", "web:pass@tcp(127.0.0.1:3306)/pastebin?parseTime=true", "MySQL data source name")
	//to be called before using the addr variable
	flag.Parse()

  //standard logger for logging information
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

  //connecting to database 
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	//initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

  //session manager for session management 
  sessionManager := scs.New()
  sessionManager.Store = mysqlstore.New(db)
  sessionManager.Lifetime = 12*time.Hour
  sessionManager.Cookie.Secure = true

  //application struct for dependency injection 
	app := &application{
		logger:        logger,
		pastes:        &models.PasteModel{DB: db},
		templateCache: templateCache,
    sessionManager: sessionManager,
	}


  tlsConfig := &tls.Config{
    CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
  }

  //new http.Server struct
  srv := &http.Server{
    Addr: *addr,
    Handler: app.routes(),
    ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
    TLSConfig: tlsConfig,

    IdleTimeout: time.Minute,
    ReadTimeout: 5*time.Second,
    WriteTimeout: 10*time.Second,
  }

	logger.Info("starting server", "addr", *addr)

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem") 

	logger.Error(err.Error())
	os.Exit(1)
}


func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
    return nil, err
	}
	fmt.Println("Connected to database")
	return db, nil
}
