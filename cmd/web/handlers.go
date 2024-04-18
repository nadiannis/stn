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
	Email          string
	Password       string
	FieldErrors    map[string]string
	NonFieldErrors []string
}

type loginForm struct {
	Email          string
	Password       string
	FieldErrors    map[string]string
	NonFieldErrors []string
}

func (app *application) homeView(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) linkListView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "link-list.tmpl.html", data)
}

func (app *application) linkCreateView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "link-create.tmpl.html", data)
}

func (app *application) linkDetailView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Println("Link ID:", id)
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "link-detail.tmpl.html", data)
}

func (app *application) linkEditView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Println("Link ID:", id)
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "link-edit.tmpl.html", data)
}

func (app *application) signupView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
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
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.FieldErrors["email"] = "Email is already in use"

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "You have successfully registered")

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (app *application) loginView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = loginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	formValues, err := decodePostForm(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := loginForm{
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
	}

	if len(form.FieldErrors) != 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	user, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.NonFieldErrors = append(form.NonFieldErrors, "Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUser", AuthenticatedUser{
		ID:    user.ID.String(),
		Email: user.Email,
	})

	http.Redirect(w, r, "/links/list", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUser")
	app.sessionManager.Put(r.Context(), "flash", "You have successfully logged out")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
