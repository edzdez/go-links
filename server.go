package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func ShortcutHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	remaining, _ := strings.CutPrefix(r.URL.Path, fmt.Sprintf("/%s", name))

	_, err := io.WriteString(w, fmt.Sprintf("shortcut: %s, remaining: %s", name, remaining))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	port := flag.Int("port", 5200, "port to listen on.")
	// db := flag.String("db", "", "database to connect to.")
	flag.Parse()

	router := http.NewServeMux()

	router.HandleFunc("GET /{name}/", ShortcutHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: withLogging(router),
	}

	log.Printf("Listening on port %d", *port)
	log.Fatal(server.ListenAndServe())
}
