package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"
)

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func WithDb(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db", db)
		next.ServeHTTP(w, r.WithContext(ctx))
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

func openDb(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		defaultDbPath, err := xdg.DataFile("go-links/shortcuts.sqlite")
		if err != nil {
			return nil, err
		}

		if _, err := os.Stat(defaultDbPath); errors.Is(err, os.ErrNotExist) {
			_, err := os.Create(defaultDbPath)
			if err != nil {
				return nil, err
			}
		}

		dbPath = defaultDbPath
	}

	log.Printf("Opening db at %s", dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	port := flag.Int("port", 5200, "port to listen on.")
	dbPath := flag.String("db", "", "database to connect to.")
	flag.Parse()

	if port == nil || dbPath == nil {
		log.Fatal("Something went wrong")
	}

	db, err := openDb(*dbPath)
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /{name}/", ShortcutHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: WithLogging(WithDb(db, router)),
	}

	log.Printf("Listening on port %d", *port)
	log.Fatal(server.ListenAndServe())
}
