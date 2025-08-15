package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

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
