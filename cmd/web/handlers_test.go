package main

import (
	"bytes"
	"io"
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
	app := newTestApplication(t)

	// We then use the httptest.NewTLSServer() function to create a new test
	// server, passing in the value returned by our app.routes() method as the
	// handler for the server. This starts up a HTTPS server which listens on a
	// randomly-chosen port of your local machine for the duration of the test.
	// Notice that we defer a call to ts.Close() so that the server is shutdown
	// when the test finishes.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	// Check the status code and body
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the response set by our
	// application for different URLs.
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
