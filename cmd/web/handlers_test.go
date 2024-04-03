package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"letsgo.skolasinski.me/internal/assert"
)

func TestPing(t *testing.T) {
	// Initialize a new httptest.ResponseRecorder
	rr := httptest.NewRecorder()

	// Initialize a new dummy http.Request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call the ping handler function, passing in
	// the defined response recorder and request
	ping(rr, r)

	// Call the Result() method on the response recorder
	rs := rr.Result()

	// Check the status code
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// And check if the response body equals to "OK"
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}

func TestPingE2E(t *testing.T) {
	// Create a new instance of our application struct. For now, this just
	// contains a couple of mock loggers (which discard anything written to
	// them).
	app := &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}

	// We then use the httptest.NewTLSServer() function to create a new test
	// server, passing in the value returned by our app.routes() method as the
	// handler for the server. This starts up a HTTPS server which listens on a
	// randomly-chosen port of your local machine for the duration of the test.
	// Notice that we defer a call to ts.Close() so that the server is shutdown
	// when the test finishes.
	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	// The network address that the test server is listening on is contained in
	// the ts.URL field. We can use this along with the ts.Client().Get() method
	// to make a GET /ping request against the test server. This returns a
	// http.Response struct containing the response.
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	// Check the status code
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// And check if the response body equals to "OK"
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
