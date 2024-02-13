package forum

import (
	middle "forum/pkg/middleware"
	s "forum/sessions"
	"html/template"
	"net/http"
	"strconv"
)

func (app *App_db) CategoryHandler(w http.ResponseWriter, r *http.Request) {

	template, err := template.ParseFiles(
		"web/templates/category-create.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	c, err := r.Cookie("session_token")
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	//check if user is admin
	if !s.GlobalSessions[c.Value].Admin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.Method == "POST" {
		if r.FormValue("create_cat") == "create" {
			if err := middle.AddCategory(app.DB, r); err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
			}
		} else {
			if err := middle.ModCategory(app.DB, r); err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
			}
		}
		http.Redirect(w, r, "/admin", http.StatusFound)
	}

	type Context struct {
		Connected   bool
		Moderator   bool
		Admin       bool
		ModLight    bool
		ID          int
		Title       string
		Description string
	}
	var context Context
	context.Connected = app.Data.Connected
	context.Moderator = s.GlobalSessions[c.Value].Moderator
	context.ModLight = s.GlobalSessions[c.Value].ModLight
	context.Admin = s.GlobalSessions[c.Value].Admin

	if r.URL.Query().Has("id") {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest)
			return
		}
		context.ID = id
		title := ""
		description := ""
		err = app.DB.QueryRow("SELECT title, descriptions FROM categories WHERE id=?", id).Scan(&title, &description)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest)
			return
		}
		context.Title = title
		context.Description = description
	}

	if err := template.Execute(w, context); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}
