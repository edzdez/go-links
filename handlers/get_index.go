package handlers

import (
	"database/sql"
	"go-links/models"
	"html/template"
	"log"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index").ParseFiles("templates/index.go.tmpl")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	db := r.Context().Value("db").(*sql.DB)

	rows, err := db.Query("SELECT * FROM shortcuts")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if rows.Close() != nil {
			log.Panic(err)
		}
	}()

	var shortcuts []models.Shortcut
	for rows.Next() {
		var s models.Shortcut
		if err := rows.Scan(&s.Name, &s.Url); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		shortcuts = append(shortcuts, s)
	}

	if err := t.ExecuteTemplate(w, "shortcuts", shortcuts); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
