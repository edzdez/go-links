package handlers

import (
	"database/sql"
	"errors"
	"go-links/models"
	"html/template"
	"log"
	"net/http"
)

func displayIndex(db *sql.DB, w http.ResponseWriter) {
	log.Println("Serving list page.")

	t, err := template.New("index").ParseFiles("templates/index.go.tmpl")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT * FROM shortcuts ORDER BY name")
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

func displayEdit(db *sql.DB, toEdit string, w http.ResponseWriter) {
	log.Println("Serving edit page")

	t, err := template.ParseFiles("templates/register.go.tmpl")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	row := db.QueryRow("SELECT url FROM shortcuts WHERE name = $1", toEdit)

	shortcut := models.Shortcut{Name: toEdit}
	err = row.Scan(&shortcut.Url)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("Shortcut does not exist")
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("Serving form")
	if err := t.ExecuteTemplate(w, "shortcut", shortcut); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sql.DB)

	if toEdit := r.URL.Query().Get("edit"); toEdit != "" {
		displayEdit(db, toEdit, w)
	} else {
		displayIndex(db, w)
	}
}
