package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
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

	db := r.Context().Value("db").(*sql.DB)
	row := db.QueryRow("SELECT url FROM shortcuts WHERE name = $1", name)

	var url string
	err := row.Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("shortcut %s not found", name)
		w.WriteHeader(http.StatusNotFound)
	} else if err != nil {
		log.Panic(err)
	} else {
		log.Printf("shortcut %s links to %s", name, url)
		http.Redirect(w, r, fmt.Sprintf("%s%s", url, remaining), http.StatusFound)
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS shortcuts (name TEXT, url TEXT)")
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
		log.Panic("Something went wrong")
	}

	db, err := openDb(*dbPath)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Panic(err)
		}
	}()

	router := http.NewServeMux()

	router.HandleFunc("GET /{name}/", ShortcutHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: WithLogging(WithDb(db, router)),
	}

	log.Printf("Listening on port %d", *port)
	log.Panic(server.ListenAndServe())
}
