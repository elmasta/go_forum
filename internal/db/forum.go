package forum

import (
	"html/template"
	"net/http"

	models "forum/pkg/models"
	s "forum/sessions"
)

// Display the home page handler
func (app *App_db) ForumHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	app.Data.Posts = nil

	template, err := template.ParseFiles(
		"web/templates/index.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	app.Data.Connected = func() bool {
		if c, err := r.Cookie("session_token"); err == nil {
			s.CheckSession(app.DB, w, r)
			app.Data.Moderator = s.GlobalSessions[c.Value].Moderator
			app.Data.Admin = s.GlobalSessions[c.Value].Admin
			app.Data.ModLight = s.GlobalSessions[c.Value].ModLight
			return true
		}
		s.CheckActive()
		return false
	}()

	err = GetRecentPosts(app, w, r)
	if err != nil {
		return
	}

	if err := template.Execute(w, app.Data); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func GetRecentPosts(app *App_db, w http.ResponseWriter, r *http.Request) error {
	var post models.Post
	rows, err := app.DB.Query("SELECT * FROM post ORDER BY rowid DESC LIMIT 5")
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Author,
			&post.Img,
			&post.Title,
			&post.Content,
			&post.CreationDate,
			&post.Flagged,
		)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return err
		}
		// get cat ids from mid table
		cat_rows, err_row := app.DB.Query("SELECT category_id FROM link_cat_post WHERE post_id=?", post.ID)
		if err_row != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return err
		}
		for cat_rows.Next() {
			var cat_id int
			err = cat_rows.Scan(&cat_id)
			if err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
				return err
			}
			post.Categories = append(post.Categories, cat_id)
			var cat_title string
			err = app.DB.QueryRow("SELECT title FROM categories WHERE id=?", cat_id).Scan(&cat_title)
			if err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
				return err
			}
			post.CategoriesName = append(post.CategoriesName, cat_title)
		}
		app.Data.Posts = append(app.Data.Posts, post)
		post.Categories = []int{}
		post.CategoriesName = []string{}
	}

	if err := rows.Err(); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return err
	}
	return nil
}
