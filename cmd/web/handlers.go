package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/nadiannis/stn/internal/models"
	"github.com/nadiannis/stn/internal/validator"
)

type signupForm struct {
	Email       string
	Password    string
	FieldErrors map[string]string
}

func (app *application) homeView(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.render(w, http.StatusOK, "home.tmpl.html", nil)
}

func (app *application) linkListView(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "link-list.tmpl.html", nil)
}

func (app *application) linkCreateView(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "link-create.tmpl.html", nil)
}

func (app *application) linkDetailView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Println("Link ID:", id)
	app.render(w, http.StatusOK, "link-detail.tmpl.html", nil)
}

func (app *application) linkEditView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Println("Link ID:", id)
	app.render(w, http.StatusOK, "link-edit.tmpl.html", nil)
}

func (app *application) signupView(w http.ResponseWriter, r *http.Request) {
	data := newTemplateData()
	data.Form = signupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	formValues, err := decodePostForm(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := signupForm{
		Email:       formValues.Get("email"),
		Password:    formValues.Get("password"),
		FieldErrors: make(map[string]string),
	}

	if strings.TrimSpace(form.Email) == "" {
		form.FieldErrors["email"] = "Email is required"
	} else if !validator.EmailRX.MatchString(form.Email) {
		form.FieldErrors["email"] = "Email is not valid"
	}

	if strings.TrimSpace(form.Password) == "" {
		form.FieldErrors["password"] = "Password is required"
	} else if utf8.RuneCountInString(form.Password) < 8 {
		form.FieldErrors["password"] = "Password should be at least 8 characters long"
	}

	if len(form.FieldErrors) != 0 {
		data := newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.FieldErrors["email"] = "Email is already in use"

			data := newTemplateData()
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (app *application) loginView(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "login.tmpl.html", nil)
}
