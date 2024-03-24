package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"letsgo.skolasinski.me/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because httprouter matches the "/" path exactly, we can now remove
	// he manual check of the r.URL.Path != "/" from this handler.

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// When httprouter is passing a request, the values of any named parameters
	// will be stored in the request context. More on request contexts later.
	params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName() method to get the value of the "id"
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)

}

// Add a new snippetCreate handler, which now is just a placeholder
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

// Rename this handler to snippetCreatePost
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way to PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError() helper.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// We use r.PostForm.Get() to retrieve the title and content
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// The r.PostForm.Get() will always return form data as a *string*.
	// We need to convert to number on our own
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Initialize a map to hold any validation errors for the form fields.
	fieldErrors := make(map[string]string)

	// Check that the Title value is not bland and is not more than 100 characters
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// Check if Content field is not blank.
	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	// Check the Expires matches one of the permitted values (1, 7 or 365).
	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	// If there are any errors dump plain HTTP response and return
	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}

	// Update the redirect path to use the new clean URL format
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
