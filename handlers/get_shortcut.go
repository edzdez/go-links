package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
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

		t, err := template.New("register").ParseFiles("templates/register.go.tmpl")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := t.ExecuteTemplate(w, "name", name); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	} else {
		log.Printf("shortcut %s links to %s", name, url)
		http.Redirect(w, r, fmt.Sprintf("%s%s", url, remaining), http.StatusFound)
	}
}
