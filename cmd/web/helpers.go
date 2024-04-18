package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
)

type AuthenticatedUser struct {
	ID    string
	Email string
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func decodePostForm(r *http.Request) (url.Values, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	return r.PostForm, nil
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		IsAuthenticated:   app.isAuthenticated(r),
		AuthenticatedUser: app.getAuthenticatedUser(r),
		Flash:             app.sessionManager.PopString(r.Context(), "flash"),
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedCtxKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func (app *application) getAuthenticatedUser(r *http.Request) *AuthenticatedUser {
	authenticatedUser, ok := app.sessionManager.Get(r.Context(), "authenticatedUser").(AuthenticatedUser)
	if !ok {
		return nil
	}

	return &authenticatedUser
}
