package forum

import (
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	s "forum/sessions"
	"html/template"
	"net/http"
)

func (app *App_db) ModHandler(w http.ResponseWriter, r *http.Request) {

	template, err := template.ParseFiles(
		"web/templates/moderation.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
		"web/templates/post-flagged.html",
		"web/templates/comment-flagged.html",
	)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	c, err := r.Cookie("session_token")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	//check if user has mod right
	if !s.GlobalSessions[c.Value].Admin && !s.GlobalSessions[c.Value].Moderator && !s.GlobalSessions[c.Value].ModLight {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	type Context struct {
		Comments  []models.Comment
		Posts     []models.Post
		Connected bool
		Moderator bool
		ModLight  bool
		Admin     bool
	}
	var context Context
	context.Comments = middle.FetchFlaggedCom(app.DB)
	context.Posts = middle.FetchFlaggedPost(app.DB)
	context.Connected = app.Data.Connected
	context.Moderator = s.GlobalSessions[c.Value].Moderator
	context.ModLight = s.GlobalSessions[c.Value].ModLight
	context.Admin = s.GlobalSessions[c.Value].Admin

	if err := template.Execute(w, context); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
	}
}

func (app *App_db) ComModHandler(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest)
	}
	template, err := template.ParseFiles(
		"web/templates/moderation.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
		"web/templates/post-flagged.html",
		"web/templates/comment-flagged.html",
	)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	//check if user has mod right
	if !s.GlobalSessions[c.Value].Admin && !s.GlobalSessions[c.Value].Moderator && !s.GlobalSessions[c.Value].ModLight {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	type Context struct {
		Comments  []models.Comment
		Connected bool
		Moderator bool
		ModLight  bool
		Admin     bool
	}
	var context Context
	context.Comments = middle.FetchFlaggedCom(app.DB)
	context.Connected = app.Data.Connected
	context.Moderator = s.GlobalSessions[c.Value].Moderator
	context.ModLight = s.GlobalSessions[c.Value].ModLight
	context.Admin = s.GlobalSessions[c.Value].Admin

	if err := template.Execute(w, context); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
	}
}
