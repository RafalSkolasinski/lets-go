package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"letsgo.skolasinski.me/ui"
)

// Add a newTemplateData() helper witch returns a pointer to templateData.
// Note that we do not use http.Request now but we will later.
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.sessionManager.PopString(r.Context(), "flash"),
		// Add the authentication status to the template data.
		IsAuthenticated: app.isAuthenticated(r),

		// Add the CSRF token.
		CSRFToken: nosurf.Token(r),
	}
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
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initialize a new buffer
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError()
	// helper and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer without any errors,
	// we are safe to ho ahead with writing to the http.ResponseWriter
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) fileServer() http.Handler {
	dir := http.FS(ui.Files)

	if app.allowFileBrowsing == nil || !(*app.allowFileBrowsing) {
		return http.FileServer(neuteredFileSystem{dir})
	} else {
		app.infoLog.Println("FileServer allows for file browsing!")
		return http.FileServer(dir)
	}
}

// Return true if the current request is from an authenticated user.
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

type neuteredFileSystem struct {
	httpDir http.FileSystem
}

func (fs neuteredFileSystem) Open(name string) (http.File, error) {
	f, err := fs.httpDir.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}

// Create a new decodePostForm() helper method. The second parameter,
// dst, is the the destination we want to decode form into.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// call ParseForm on the request
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call Decode() on our decoder instance, passing the target destination.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use invalid destination, the Decode() method will
		// return an error with the type *form.InvalidDecoderError.
		var InvalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &InvalidDecoderError) {
			panic(err)
		}

		// For regular errors
		return err
	}

	return nil
}
