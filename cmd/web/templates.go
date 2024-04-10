package main

import (
	"io/fs"
	"path/filepath"
	"text/template"
	"time"

	"letsgo.skolasinski.me/internal/models"
	"letsgo.skolasinski.me/ui"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string

	// Add an IsAuthenticated field to the templateData struct.
	IsAuthenticated bool

	// Add a CSRFToken field.
	CSRFToken string
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable.
// This is essentially a string-keyed map which acts as a lookup between
// the names of our custom template functions and the function themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl.html". This will essentially give
	// us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl.html ui/html/pages/view.tmpl.html].
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one
	for _, page := range pages {
		// Extract the file name (like home.tmpl.html) from full filepath
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates we
		// want to parse.
		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		ts := template.New(name).Funcs(functions)

		// Use ParseFS() instead of ParseFiles() to parse the template files
		// from the ui.Files embedded filesystem.
		ts, err := ts.ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map as normal.
		cache[name] = ts
	}
	return cache, nil
}
