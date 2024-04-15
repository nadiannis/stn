package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
}

func main() {
	port := flag.Int("port", 8080, "Web server port")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", *port),
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on :%d", *port)
	errorLog.Fatal(srv.ListenAndServe())
}
