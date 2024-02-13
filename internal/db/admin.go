package forum

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	//"database/sql"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	s "forum/sessions"
	//"fmt"
	//"time"
)

func (app *App_db) AdminHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		"web/templates/admin.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
		"web/templates/comment-flagged.html",
		"web/templates/post-flagged.html",
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
	// check if user is admin
	if !s.GlobalSessions[c.Value].Admin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if r.Method == "POST" {
		actions := map[string]func(db *sql.DB, r *http.Request) error{
			"deletion":      middle.RmUser,
			"del_mod":       middle.DelMod,
			"add_mod":       app.addModPrivileges(2),
			"add_mod_light": app.addModPrivileges(4),
			"del_cat":       middle.DelCategory,
			"del_post":      middle.DelPost,
			"del_com":       middle.DelCom,
			"del_com_flag":  middle.DelComFlag,
			"del_post_flag": middle.DelPostFlag,
		}

		for formValue, action := range actions {
			if len(r.FormValue(formValue)) > 0 {
				if err := action(app.DB, r); err != nil {
					ErrorHandler(w, r, http.StatusInternalServerError)
					return
				}
			}
		}
	}

	type Context struct {
		UsersList  []models.User
		Categories []models.Categories
		Comments   []models.Comment
		Posts      []models.Post
		Connected  bool
		Moderator  bool
		ModLight   bool
		Admin      bool
	}
	var context Context
	context.UsersList = middle.FetchUsers(app.DB)
	context.Categories = middle.FetchCat(app.DB, []int{0})
	context.Comments = middle.FetchFlaggedCom(app.DB)
	context.Posts = middle.FetchFlaggedPost(app.DB)
	context.Connected = app.Data.Connected
	context.Moderator = s.GlobalSessions[c.Value].Moderator
	context.ModLight = s.GlobalSessions[c.Value].ModLight
	context.Admin = s.GlobalSessions[c.Value].Admin

	if err := template.Execute(w, context); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func (app *App_db) addModPrivileges(modType int) func(db *sql.DB, r *http.Request) error {
	return func(db *sql.DB, r *http.Request) error {
		var id int
		var err error
		if modType == 2 {
			id, err = strconv.Atoi(r.FormValue("add_mod"))
			if err != nil {
				return err
			}
		} else if modType == 4 {
			id, err = strconv.Atoi(r.FormValue("add_mod_light"))
			if err != nil {
				return err
			}

		}
		return middle.AddMod(db, r, modType, id)
	}
}
