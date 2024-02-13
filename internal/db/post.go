package forum

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	s "forum/sessions"
)

func (app *App_db) PostedIdHandler(w http.ResponseWriter, r *http.Request, current_user int64, post models.Post, message string) {
	template, err := template.ParseFiles(
		"web/templates/edit-post.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// potential more work TODO
	// app.Data.Categories = middle.FetchCat(app.DB, app.Data.CurrentPost.Categories)
	app.Data.Categories = middle.FetchCat(app.DB, post.Categories)
	app.Data.ErrMessage = message
	app.Data.CurrentPost = post

	// ReturnCurrentPost(app, w, r, current_user)
	if err := template.Execute(w, app.Data); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func (app *App_db) PostIdHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		"web/templates/post-id.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/comment.html",
		"web/templates/comment-create.html",
		"web/templates/footer.html",
	)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	var current_user int64
	c, _ := r.Cookie("session_token")
	if c != nil {
		current_user = s.GlobalSessions[c.Value].UserID
	}
	if !ReturnCurrentPost(app, w, r, current_user) {
		return
	}
	ReturnComment(app, w, r, current_user)
	switch r.Method {
	case "POST":
		if app.Data.Connected {
			// create comment
			if r.FormValue("content") != "" {
				var comment models.Comment
				comment.AuthorID = current_user
				comment.Content = r.FormValue("content")
				comment.PostID = app.Data.CurrentPost.ID
				middle.CreateComment(app.DB, &comment)
			}
			// like comment
			if r.FormValue("like") != "" {
				like := strings.Split(r.FormValue("like"), " ")[0] == "true"
				id_comment, err := strconv.Atoi(strings.Split(r.FormValue("like"), " ")[1])
				if err != nil {
					ErrorHandler(w, r, http.StatusInternalServerError)
					return
				}
				middle.UpdateLike(app.DB, int64(id_comment), current_user, like)
			}
			// like post
			if r.FormValue("like-post") != "" {
				like := strings.Split(r.FormValue("like-post"), " ")[0] == "true"
				id_post, err := strconv.Atoi(strings.Split(r.FormValue("like-post"), " ")[1])
				if err != nil {
					ErrorHandler(w, r, http.StatusInternalServerError)
					return
				}
				middle.UpdateLikePost(app.DB, int64(id_post), current_user, like)
				ReturnCurrentPost(app, w, r, current_user)
			}
			// delete comment
			if r.FormValue("delete") != "" {
				id_comment, err := strconv.Atoi(r.FormValue("delete"))
				if err != nil {
					ErrorHandler(w, r, http.StatusInternalServerError)
					return
				}
				if s.GlobalSessions[c.Value].Moderator || s.GlobalSessions[c.Value].Admin || s.GlobalSessions[c.Value].ModLight {
					middle.DelCom(app.DB, r)
				} else {
					middle.RemoveComment(app.DB, int64(id_comment), current_user, false)
				}
			}
			// edit comment
			if r.FormValue("edit-comment") != "" {
				id_comment, err := strconv.Atoi(r.FormValue("edit-comment"))
				if err != nil {
					ErrorHandler(w, r, http.StatusInternalServerError)
					return
				}
				if s.GlobalSessions[c.Value].Moderator || s.GlobalSessions[c.Value].Admin || s.GlobalSessions[c.Value].ModLight {
					app.CommentHandler(w, r, int64(id_comment), -1)
				} else {
					app.CommentHandler(w, r, int64(id_comment), current_user)
				}
				return
			}
			if r.FormValue("comment-editor") != "" {
				var comment models.Comment
				comment.Content = r.FormValue("content-editor")
				id, err := strconv.Atoi(r.FormValue("comment-editor"))
				if err != nil {
					ErrorHandler(w, r, http.StatusInternalServerError)
					return
				}
				comment.ID = int64(id)
				middle.UpdateComment(app.DB, &comment)
			}
			// flag comment
			if r.FormValue("report") != "" {
				middle.FlagComment(app.DB, r)
			}
			// flag post
			if r.FormValue("report-post") != "" {
				middle.FlagPost(app.DB, r)
				http.Redirect(w, r, "id?id="+r.FormValue("report-post"), http.StatusFound)
			}
			if app.Data.CurrentPost.AuthorID == current_user || s.GlobalSessions[c.Value].Moderator || s.GlobalSessions[c.Value].Admin {
				// delete post
				if r.FormValue("delete-post") != "" {
					middle.RemovePost(app.DB, app.Data.CurrentPost.ID)
					http.Redirect(w, r, "/post", http.StatusFound)
				}
				// edit post
				if r.FormValue("edit-post") != "" {
					ReturnCurrentPost(app, w, r, current_user)
					app.PostedIdHandler(w, r, current_user, app.Data.CurrentPost, "")
					return
				}
				if r.FormValue("post-editor") != "" {
					img, imgerr := InitImg(r)
					// fmt.Println(imgerr)
					var post models.Post
					post.Title = r.FormValue("title-editor")
					post.Content = r.FormValue("content-editor")
					post.ID = app.Data.CurrentPost.ID
					// to change category
					var cat []int
					for _, v := range r.Form["categories-editor"] {
						temp, err := strconv.Atoi(v)
						if err != nil {
							ErrorHandler(w, r, http.StatusInternalServerError)
							return
						}
						cat = append(cat, temp)
					}
					post.Categories = cat
					if imgerr != nil && imgerr.Error() != "http: no such file" {
						app.PostedIdHandler(w, r, current_user, post, imgerr.Error())
						// app.PageCreatePost(w, r, *post, "Bad format")
						return
					}
					// update img
					if r.FormValue("deleteimg") == "true" {
						middle.UpdateImgPoste(app.DB, post.ID, "")
					} else if imgerr == nil {
						mkdirPostAsset(app, post.ID, &post, img)
					}
					middle.UpdateCategory(app.DB, &post)
					// update category
					middle.UpdatePost(app.DB, &post)
					ReturnCurrentPost(app, w, r, current_user)
				}
			}
		}
	}
	ReturnComment(app, w, r, current_user)
	if err := template.Execute(w, app.Data); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func ReturnCurrentPost(app *App_db, w http.ResponseWriter, r *http.Request, current_user int64) bool {
	var post models.Post
	if r.URL.Query().Has("id") {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			ErrorHandler(w, r, http.StatusNotFound)
			return false
		}

		err = app.DB.QueryRow("SELECT * FROM post where id = ?", id).Scan(
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
			ErrorHandler(w, r, http.StatusNotFound)
			return false
		}
		// get middle table
		rows, err := app.DB.Query("SELECT category_id FROM link_cat_post WHERE post_id = ?", post.ID)
		for rows.Next() {
			var cat_id int
			err = rows.Scan(&cat_id)
			if err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
			}
			post.Categories = append(post.Categories, cat_id)
			var cat_name string
			err = app.DB.QueryRow("SELECT title FROM categories WHERE id=?", cat_id).Scan(&cat_name)
			if err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
			}
			post.CategoriesName = append(post.CategoriesName, cat_name)
		}

		post.User_like, post.User_dislike = linkPost(app, post.ID)
		post.Like, post.Dislike = len(post.User_like), len(post.User_dislike)
		post.IfCurrentUser = post.AuthorID == current_user
		post.Ifimg = post.Img != ""
		app.Data.CurrentPost = post
		if err != nil {
			if err == sql.ErrNoRows {
				ErrorHandler(w, r, http.StatusNotFound)
			} else {
				ErrorHandler(w, r, http.StatusInternalServerError)
			}
			return false
		}
		post.CategoriesName = []string{}
		post.Categories = []int{}
	}
	return true
}

// Handler that shows the post creation page and ensures that users are certified to create posts.
func (app *App_db) PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Checking for rights to access this page
	cookie, errCookie := r.Cookie("session_token")
	if errCookie != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	s.CheckActive()
	_, ok := s.GlobalSessions[cookie.Value]
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case "GET":
		app.PageCreatePost(w, r, models.Post{}, "")
		return
	case "POST":
		var post *models.Post
		img, errimg := InitImg(r)
		errParse := r.ParseForm()
		if errParse != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}

		var cat []int
		for _, v := range r.Form["categories"] {
			temp, _ := strconv.Atoi(v)
			cat = append(cat, temp)
		}
		post = &models.Post{
			Title:      r.FormValue("title"),
			Content:    r.FormValue("content"),
			Categories: cat,
		}
		if errimg != nil && errimg.Error() != "http: no such file" {
			app.PageCreatePost(w, r, *post, errimg.Error())
			return
		}

		err := app.DB.QueryRow("SELECT id, username FROM users where session_token = ?", cookie.Value).Scan(&post.AuthorID, &post.Author)
		if err != nil {
			log.Fatal(err)
		}

		id, err_create_post := middle.CreatePost(app.DB, post)
		if err_create_post != nil {
			ErrorHandler(w, r, http.StatusBadRequest)
			return
		}
		mkdirPostAsset(app, int64(id), post, img)
		http.Redirect(w, r, "/post/id?id="+strconv.Itoa(id), http.StatusFound)
	}
}

func (app *App_db) PageCreatePost(w http.ResponseWriter, r *http.Request, post models.Post, message string) {
	template, err := template.ParseFiles(
		"web/templates/post-create.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	// app.Data.Categories = middle.FetchCat(app.DB, []int{0})
	app.Data.ErrMessage = message
	app.Data.CurrentPost = post
	app.Data.Categories = middle.FetchCat(app.DB, post.Categories)

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

	if err := template.Execute(w, app.Data); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}
