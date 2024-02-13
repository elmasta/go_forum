package forum

import (
	"html/template"
	"net/http"
	"strconv"

	//"database/sql"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	s "forum/sessions"
	//"time"
)

func (app *App_db) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		"web/templates/profile.html",
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		if len(r.FormValue("ask_mod")) > 0 {
			id, err := strconv.Atoi(r.FormValue("ask_mod"))
			if err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
				return
			}
			if err := middle.AskModerator(app.DB, r, 1, id); err != nil {
				ErrorHandler(w, r, http.StatusBadRequest)
				return
			}
		} else if len(r.FormValue("ask_light_mod")) > 0 {
			id, err := strconv.Atoi(r.FormValue("ask_light_mod"))
			if err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
				return
			}
			if err := middle.AskModerator(app.DB, r, 2, id); err != nil {
				ErrorHandler(w, r, http.StatusBadRequest)
				return
			}
		}
	}

	type Context struct {
		User      models.User
		Connected bool
		Moderator bool
		ModLight  bool
		Admin     bool
		Data      models.Dataprofile
	}
	var context Context
	if cookie, err := r.Cookie("session_token"); err == nil {
		context.User = middle.FetchUser(app.DB, cookie.Value)
		context.Connected = app.Data.Connected
		context.Moderator = s.GlobalSessions[c.Value].Moderator
		context.ModLight = s.GlobalSessions[c.Value].ModLight
		context.Admin = s.GlobalSessions[c.Value].Admin
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	currentuserid := s.GlobalSessions[c.Value].UserID
	context.Data.Likedcomment = middle.Likedcomment(app.DB, currentuserid)
	context.Data.Dislikedcomment = middle.Dislikedcomment(app.DB, currentuserid)
	context.Data.Likedpost = middle.Likedpost(app.DB, currentuserid)
	context.Data.Dislikedpost = middle.Dislikedpost(app.DB, currentuserid)
	context.Data.Notified = middle.ProfilNotified(app.DB, currentuserid)

	if err := template.Execute(w, context); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}
