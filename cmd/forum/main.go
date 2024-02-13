package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	. "forum/internal/db"
)

func main() {
	db, err := sql.Open("sqlite3", "config/db/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	app := InitDB(db)
	app.DB.Exec("PRAGMA foreign_keys = ON")
	if err := app.Migrate(); err != nil {
		log.Fatal(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	img := http.FileServer(http.Dir("web/static/upload/img"))
	http.Handle("/img/", http.StripPrefix("/img", img))

	http.HandleFunc("/", app.ForumHandler)
	http.HandleFunc("/admin", app.AdminHandler)
	http.HandleFunc("/moderation", app.ModHandler)
	http.HandleFunc("/com_moderation", app.ComModHandler)
	http.HandleFunc("/profile", app.ProfileHandler)
	http.HandleFunc("/login", app.LoginHandler)
	http.HandleFunc("/register", app.RegisterHandler)
	http.HandleFunc("/logout", app.LogoutHandler)

	// Post related handlers
	http.HandleFunc("/category", app.CategoryHandler)
	http.HandleFunc("/post/create", app.PostCreateHandler)
	http.HandleFunc("/post", app.PostHandler)
	http.HandleFunc("/post/id", app.PostIdHandler)

	fmt.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
