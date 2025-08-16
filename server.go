package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-links/handlers"
	"go-links/middleware"

	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"
)

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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS shortcuts (name TEXT PRIMARY KEY, url TEXT)")
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

	router.HandleFunc("GET /", handlers.IndexHandler)
	router.HandleFunc("GET /{name}/", handlers.ShortcutHandler)
	router.HandleFunc("POST /{name}", handlers.RegisterShortcutHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: middleware.WithLogging(middleware.WithStatic(middleware.WithDb(db, router))),
	}

	log.Printf("Listening on port %d", *port)
	log.Panic(server.ListenAndServe())
}
