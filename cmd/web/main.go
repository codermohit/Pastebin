package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"capybara.pastebin.xyz/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger        *slog.Logger
	pastes        *models.PasteModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	//data source name : connection string that describes how to connect to your database
	dsn := flag.String("dsn", "web:pass@tcp(127.0.0.1:3306)/pastebin?parseTime=true", "MySQL data source name")
	//to be called before using the addr variable
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
  if err != nil {
    logger.Error(err.Error())
    os.Exit(1)
  }
	defer db.Close()


  //initialize a new template cache 
  templateCache, err := newTemplateCache()
  if err!=nil{
    logger.Error(err.Error())
    os.Exit(1)
  }


	app := &application{
		logger: logger,
		pastes: &models.PasteModel{DB: db},
    templateCache: templateCache,
	}

	//use the Info() method to log the starting server message at Info severity
	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to database")

	err = db.Ping()
	if err != nil {
		db.Close()
	}

	return db, nil
}
