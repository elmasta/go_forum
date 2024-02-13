package forum

import (
	"database/sql"
	models "forum/pkg/models"
	"log"
	"net/http"
	"strconv"
)

func FetchFlaggedCom(db *sql.DB) []models.Comment {

	rows, err := db.Query("SELECT id, author_id, post_id, content, creation, flagged FROM comment WHERE flagged=?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err = rows.Scan(&comment.ID, &comment.AuthorID, &comment.PostID, &comment.Content, &comment.CreationDate, &comment.Flagged)
		if err != nil {
			log.Fatal(err)
		}
		//get author name and category name
		err = db.QueryRow("SELECT username FROM users WHERE id = ?", comment.AuthorID).Scan(&comment.Author)
		if err != nil {
			log.Fatal(err)
		}
		err = db.QueryRow("SELECT title FROM post WHERE id = ?", comment.PostID).Scan(&comment.Post)
		if err != nil {
			log.Fatal(err)
		}
		comments = append(comments, comment)
	}
	return comments
}

func FetchFlaggedPost(db *sql.DB) []models.Post {

	rows, err := db.Query("SELECT id, author_id, content, creation, flagged FROM post WHERE flagged=?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.CreationDate, &post.Flagged)
		if err != nil {
			log.Fatal(err)
		}
		//get author name and category name
		err = db.QueryRow("SELECT username FROM users WHERE id = ?", post.AuthorID).Scan(&post.Author)
		if err != nil {
			log.Fatal(err)
		}

		//get the categories
		cat_rows, err_row := db.Query("SELECT category_id FROM link_cat_post WHERE post_id=?", post.ID)
		if err_row != nil {
			log.Fatal(err)
		}
		defer cat_rows.Close()
		for cat_rows.Next() {
			var cat_id int
			err = cat_rows.Scan(&cat_id)
			if err != nil {
				log.Fatal(err)
			}
			//get the categories names
			temp := ""
			err = db.QueryRow("SELECT title FROM categories WHERE id = ?", cat_id).Scan(&temp)
			if err != nil {
				log.Fatal(err)
			}
			post.CategoriesName = append(post.CategoriesName, temp)
			post.Categories = append(post.Categories, cat_id)
		}
		posts = append(posts, post)
	}
	return posts
}

func DelPost(db *sql.DB, r *http.Request) error {

	id, err := strconv.Atoi(r.FormValue("del_post"))
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM post WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}

func DelPostFlag(db *sql.DB, r *http.Request) error {

	id, err := strconv.Atoi(r.FormValue("del_post_flag"))
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE post SET flagged=? WHERE id=?", 0, id)
	if err != nil {
		return err
	}
	return nil
}

func DelCom(db *sql.DB, r *http.Request) error {

	id := 0
	var err error
	if len(r.FormValue("del_com")) == 0 {
		id, err = strconv.Atoi(r.FormValue("delete"))
	} else {
		id, err = strconv.Atoi(r.FormValue("del_com"))
	}

	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM comment WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}

func DelComFlag(db *sql.DB, r *http.Request) error {

	id, err := strconv.Atoi(r.FormValue("del_com_flag"))
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE comment SET flagged=? WHERE id=?", 0, id)
	if err != nil {
		return err
	}
	return nil
}
