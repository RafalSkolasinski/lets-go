package main

import (
	"net/http"
	"os"

	"github.com/justinas/alice"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	var fileServer http.Handler
	if *app.allowFileBrowsing {
		app.infoLog.Println("FileServer allows for file browsing!")
		fileServer = http.FileServer(http.Dir("./ui/static/"))
	} else {
		fileServer = http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	}
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Create a middleware chain containing our 'standard' middleware.
	chain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux
	return chain.Then(mux)
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
