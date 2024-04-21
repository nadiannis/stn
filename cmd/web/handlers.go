package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nadiannis/stn/internal/models"
	"github.com/nadiannis/stn/internal/validator"
)

type authForm struct {
	Email    string
	Password string
	validator.Validator
}

type linkForm struct {
	URL      string
	BackHalf string
	validator.Validator
}

func (app *application) homeView(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	data := app.newTemplateData(r)
	data.Form = linkForm{}
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) linkRedirect(w http.ResponseWriter, r *http.Request) {
	backHalf := r.PathValue("backhalf")

	link, err := app.links.GetByBackHalf(backHalf)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.links.UpdateEngagements(link.ID.String(), int(link.Engagements)+1)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, link.URL, http.StatusFound)
}

func (app *application) linkListView(w http.ResponseWriter, r *http.Request) {
	userID := app.getAuthenticatedUser(r).ID

	links, err := app.links.GetByUserID(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	summary, err := app.links.GetSummaryByUserID(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Links = links
	data.Summary = summary
	app.render(w, http.StatusOK, "link-list.tmpl.html", data)
}

func (app *application) linkCreateView(w http.ResponseWriter, r *http.Request) {
	userID := app.getAuthenticatedUser(r).ID

	summary, err := app.links.GetSummaryByUserID(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = linkForm{}
	data.Summary = summary
	app.render(w, http.StatusOK, "link-create.tmpl.html", data)
}

func (app *application) linkCreate(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	if from == "" {
		from = "home"
	}

	formValues, err := decodeForm(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := linkForm{
		URL:      formValues.Get("url"),
		BackHalf: formValues.Get("back-half"),
	}

	form.CheckField(validator.NotEmpty(form.URL), "url", "URL is required")
	form.CheckField(validator.Matches(form.URL, validator.URLRegex), "url", "URL is not valid")

	if validator.NotEmpty(form.BackHalf) {
		form.CheckField(validator.Matches(form.BackHalf, validator.BackHalfRegex), "backHalf", "Back-half can only contain letters, numbers, & the characters _-")
	}

	if !form.Valid() {
		if from == "home" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]any{
				"fieldErrors":    form.FieldErrors,
				"nonFieldErrors": form.NonFieldErrors,
			})
			return
		}

		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, from+".tmpl.html", data)
		return
	}

	if !validator.NotEmpty(form.BackHalf) {
		for {
			backHalf := randString(5)
			exists, err := app.links.BackHalfExists(backHalf)
			if err != nil {
				app.serverError(w, err)
				return
			}

			if !exists {
				form.BackHalf = backHalf
				break
			}
		}
	}

	var userID string
	if user := app.getAuthenticatedUser(r); user != nil {
		userID = user.ID
	}

	link, err := app.links.Insert(form.URL, form.BackHalf, userID)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateBackHalf) {
			form.AddFieldError("backHalf", "Back-half is already in use")

			if from == "home" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]any{
					"fieldErrors":    form.FieldErrors,
					"nonFieldErrors": form.NonFieldErrors,
				})
				return
			}

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, from+".tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	if from == "home" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"url":      form.URL,
			"backHalf": form.BackHalf,
		})
		return
	}

	http.Redirect(w, r, "/links/"+link.ID.String(), http.StatusSeeOther)
}

func (app *application) linkDetailView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	link, err := app.links.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	var userID string
	if user := app.getAuthenticatedUser(r); user != nil {
		userID = user.ID
	}

	if !link.UserID.Valid || (link.UserID.Valid && link.UserID.V.String() != userID) {
		http.Redirect(w, r, "/links/list", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.Link = link
	app.render(w, http.StatusOK, "link-detail.tmpl.html", data)
}

func (app *application) linkEditView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	link, err := app.links.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	var userID string
	if user := app.getAuthenticatedUser(r); user != nil {
		userID = user.ID
	}

	if !link.UserID.Valid || (link.UserID.Valid && link.UserID.V.String() != userID) {
		http.Redirect(w, r, "/links/list", http.StatusSeeOther)
		return
	}

	summary, err := app.links.GetSummaryByUserID(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Link = link
	data.Summary = summary
	app.render(w, http.StatusOK, "link-edit.tmpl.html", data)
}

func (app *application) linkEdit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	formValues, err := decodeForm(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := linkForm{
		URL:      formValues.Get("url"),
		BackHalf: formValues.Get("back-half"),
	}

	form.CheckField(validator.NotEmpty(form.URL), "url", "URL is required")
	form.CheckField(validator.Matches(form.URL, validator.URLRegex), "url", "URL is not valid")

	if validator.NotEmpty(form.BackHalf) {
		form.CheckField(validator.Matches(form.BackHalf, validator.BackHalfRegex), "back-half", "Back-half can only contain letters, numbers, & the characters _-")
	}

	if !form.Valid() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"fieldErrors":    form.FieldErrors,
			"nonFieldErrors": form.NonFieldErrors,
		})
		return
	}

	if !validator.NotEmpty(form.BackHalf) {
		for {
			backHalf := randString(5)
			exists, err := app.links.BackHalfExists(backHalf)
			if err != nil {
				app.serverError(w, err)
				return
			}

			if !exists {
				form.BackHalf = backHalf
				break
			}
		}
	}

	err = app.links.Update(id, form.URL, form.BackHalf)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateBackHalf) {
			form.AddFieldError("back-half", "Back-half is already in use")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]any{
				"fieldErrors":    form.FieldErrors,
				"nonFieldErrors": form.NonFieldErrors,
			})
			return
		} else {
			app.serverError(w, err)
		}
		return
	}

	http.Redirect(w, r, "/links/"+id, http.StatusSeeOther)
}

func (app *application) linkDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := app.links.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/links/list", http.StatusSeeOther)
}

func (app *application) signupView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = authForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	formValues, err := decodeForm(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := authForm{
		Email:    formValues.Get("email"),
		Password: formValues.Get("password"),
	}

	form.CheckField(validator.NotEmpty(form.Email), "email", "Email is required")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "Email is not valid")

	form.CheckField(validator.NotEmpty(form.Password), "password", "Password is required")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password should be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email is already in use")

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
	data.Form = authForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	formValues, err := decodeForm(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := authForm{
		Email:    formValues.Get("email"),
		Password: formValues.Get("password"),
	}

	form.CheckField(validator.NotEmpty(form.Email), "email", "Email is required")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "Email is not valid")

	form.CheckField(validator.NotEmpty(form.Password), "password", "Password is required")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	user, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

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
