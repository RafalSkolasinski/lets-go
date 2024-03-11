package main

import (
	"log"
	"net/http"
	"os"

	"github.com/spf13/pflag"
)

type application struct {
	errorLog          *log.Logger
	infoLog           *log.Logger
	allowFileBrowsing *bool
}

func main() {
	addr := pflag.StringP("addr", "", ":4000", "HTTP network address")
	allowFileBrowsing := pflag.BoolP("allow-file-browsing", "", false, "Disable file browsing in File Server")

	pflag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,

		allowFileBrowsing: allowFileBrowsing,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,

		// Call the new app.routes() method to get the servemux containing our routes.
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
