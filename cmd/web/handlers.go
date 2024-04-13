package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Home"))
}

func (app *application) linkList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("View all links"))
}

func (app *application) linkCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("View create short link page"))
	case http.MethodPost:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Create short link"))
	default:
		allowedMethods := fmt.Sprintf("%s, %s", http.MethodGet, http.MethodPost)
		w.Header().Set("Allow", allowedMethods)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
