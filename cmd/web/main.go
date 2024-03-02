package main

import (
	"log"
	"net/http"
	"os"

	"github.com/spf13/pflag"
)

type config struct {
	addr string
}

var cfg config

func main() {
	pflag.StringVarP(&cfg.addr, "addr", "", ":4000", "HTTP network address")

	pflag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Starting server on %s\n", cfg.addr)
	err := http.ListenAndServe(cfg.addr, mux)
	log.Fatal(err)
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
