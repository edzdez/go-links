package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

func ShortcutHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET %s", r.URL.Path)

	_, err := io.WriteString(w, r.URL.Path)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	port := flag.Int("port", 5200, "port to listen on.")
	// db := flag.String("db", "", "database to connect to.")
	flag.Parse()

	http.HandleFunc("GET /{name}/", ShortcutHandler)

	log.Printf("Listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
