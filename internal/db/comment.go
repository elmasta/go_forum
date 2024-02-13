package forum

import (
	"database/sql"
	"html/template"
	"net/http"

	models "forum/pkg/models"
)

func (app *App_db) CommentHandler(w http.ResponseWriter, r *http.Request, id_comment, current_user int64) {
	template, err := template.ParseFiles(
		"web/templates/edit-comment.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	var comment models.Comment
	err = app.DB.QueryRow("Select id,author_id,content,post_id From comment where id = ?", id_comment).Scan(
		&comment.ID,
		&comment.AuthorID,
		&comment.Content,
		&comment.PostID,
	)
	if current_user != comment.AuthorID && current_user != -1 {
		return
	}
	if err != nil {
		if err == sql.ErrNoRows {
			ErrorHandler(w, r, http.StatusNotFound)
		} else {
			ErrorHandler(w, r, http.StatusInternalServerError)
		}
		return
	}
	app.Data.CurrentComment = comment
	if err := template.Execute(w, app.Data); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func ReturnComment(app *App_db, w http.ResponseWriter, r *http.Request, current_user int64) {
	var tab_comment []models.Comment
	var comment models.Comment
	rows, err := app.DB.Query("SELECT id, author_id, post_id, img, content, creation, flagged FROM comment WHERE post_id = ?", app.Data.CurrentPost.ID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		rows.Scan(&comment.ID, &comment.AuthorID, &comment.PostID, &comment.Img, &comment.Content, &comment.CreationDate, &comment.Flagged)
		comment.PostID = app.Data.CurrentPost.ID
		err = app.DB.QueryRow("SELECT username FROM users where id = ?", comment.AuthorID).Scan(&comment.Author)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
		comment.IfCurrentUser = comment.AuthorID == current_user
		comment.User_like, comment.User_dislike = linkComment(app, w, r, comment.ID)
		comment.Like = len(comment.User_like)
		comment.Dislike = len(comment.User_dislike)
		tab_comment = append(tab_comment, comment)
		comment = models.Comment{}
	}
	app.Data.CurrentPost.Tab_comment = tab_comment
}

func linkComment(app *App_db, w http.ResponseWriter, r *http.Request, id int64) (tab_like map[int64]bool, tab_dislike map[int64]bool) {
	tab_like, tab_dislike = make(map[int64]bool), make(map[int64]bool)
	rows, err := app.DB.Query("SELECT user_id,likes FROM link_comment WHERE comment_id = ?", id)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return nil, nil
	}
	var user_id int64
	var like bool
	for rows.Next() {
		rows.Scan(&user_id, &like)
		if like {
			tab_like[user_id] = true
		} else {
			tab_dislike[user_id] = true
		}
	}
	return
}
