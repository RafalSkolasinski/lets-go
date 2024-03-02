package main

import (
	"log"
	"net/http"
	"os"

	"github.com/spf13/pflag"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := pflag.StringP("addr", "", ":4000", "HTTP network address")

	pflag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
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
