package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
)

func RegisterShortcutHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "Bad Form Input", http.StatusBadRequest)
		return
	}

	name := r.PathValue("name")
	formUrl := r.FormValue("url")

	parsedUrl, err := url.Parse(formUrl)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad URL", http.StatusBadRequest)
		return
	}

	db := r.Context().Value("db").(*sql.DB)
	_, err = db.Exec("INSERT OR REPLACE INTO shortcuts (name, url) VALUES (?, ?)", name, parsedUrl.String())
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Register shortcut %s to %s", name, parsedUrl.String())
	http.Redirect(w, r, parsedUrl.String(), http.StatusFound)
}
