package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/spf13/pflag"

	_ "github.com/go-sql-driver/mysql"
)

const DEFAULT_DSN = "snippetbox:snippetbox@tcp(localhost:3306)/snippetbox?parseTime=true"

type application struct {
	errorLog          *log.Logger
	infoLog           *log.Logger
	allowFileBrowsing *bool
}

func main() {
	addr := pflag.StringP("addr", "", ":4000", "HTTP network address")
	dsn := pflag.StringP("dsn", "", DEFAULT_DSN, "MySQL data source name")

	allowFileBrowsing := pflag.BoolP("allow-file-browsing", "", false, "Disable file browsing in File Server")

	pflag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

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
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
