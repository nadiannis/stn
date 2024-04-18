package main

import (
	"database/sql"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nadiannis/stn/internal/models"
)

type application struct {
	infoLog        *log.Logger
	errorLog       *log.Logger
	templateCache  map[string]*template.Template
	users          *models.UserModel
	sessionManager *scs.SessionManager
}

func main() {
	port := flag.Int("port", 8080, "Web server port")
	dsn := flag.String("db-dsn", "", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = true

	gob.Register(AuthenticatedUser{})

	app := application{
		infoLog:        infoLog,
		errorLog:       errorLog,
		templateCache:  templateCache,
		users:          &models.UserModel{DB: db},
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", *port),
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on :%d", *port)
	errorLog.Fatal(srv.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
