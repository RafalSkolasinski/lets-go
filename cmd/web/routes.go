package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Define routes for our application using the httprouter package
func (app *application) routes() http.Handler {
	// Initialize the router
	router := httprouter.New()

	// Set a custom handler function which wraps our notFound() helper.
	// This is so all 404 responses are the same and not partially
	// handled automatically by the httrouter and our helper.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Update the pattern for the fileserver
	fileServer := app.fileServer()
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// New dynamic middleware chain
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Update these routes to use the dynamic middleware chain
	// Note tat ThenFunc() returns an http.Handler rather than http.HandlerFunc.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// Create a middleware chain containing our 'standard' middleware.
	chain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux
	return chain.Then(router)
}
