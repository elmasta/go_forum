package forum

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	models "forum/pkg/models"
)

func CreateComment(db *sql.DB, comment *models.Comment) (int, error) {
	_, err := db.Exec("INSERT INTO comment(author_id, post_id, img, content, creation) VALUES(?,?,?,?, datetime())",
		comment.AuthorID, comment.PostID, "", comment.Content)
	if err != nil {
		return 0, err
	}

	sql, _ := db.Exec("SELECT last_insert_rowid()")
	id, _ := sql.LastInsertId()
	return int(id), nil
}

func UpdateComment(db *sql.DB, comment *models.Comment) error {
	_, err := db.Exec("UPDATE comment SET content = ? WHERE id = ?", comment.Content, comment.ID)
	if err != nil {
		return err
	}
	return nil
}

func ReportComment(db *sql.DB, comment int) error {
	_, err := db.Exec("UPDATE comment SET flagged = ? WHERE id = ?", 1, comment)
	if err != nil {
		return err
	}
	return nil
}

func RemoveComment(db *sql.DB, id_comment, current_user int64, isdeletpost bool) error {
	var err error
	if isdeletpost {
		_, err = db.Exec("DELETE FROM comment WHERE id = ? ", id_comment)
		if err != nil {
			fmt.Println("Remove comment : ", err)
			return err
		}
	} else {
		_, err = db.Exec("DELETE FROM comment WHERE id = ? AND author_id = ? ", id_comment, current_user)
		if err != nil {
			fmt.Println("Remove comment : ", err)
			return err
		}
	}
	_, err = db.Exec("DELETE FROM link_comment WHERE comment_id = ?", id_comment)
	if err != nil {
		fmt.Println("Remove post : ", err)
		return err
	}
	return nil
}

func UpdateLike(db *sql.DB, id_comment, id_user int64, like bool) error {
	var exist bool
	err := db.QueryRow("SELECT EXISTS( SELECT * FROM link_comment WHERE user_id = ? AND comment_id = ?) AS exist", id_user, id_comment).Scan(&exist)
	if err != nil {
		return err
	}
	if exist {
		var current bool
		err := db.QueryRow("SELECT likes FROM link_comment WHERE user_id = ? AND comment_id = ?", id_user, id_comment).Scan(&current)
		if err != nil {
			return err
		}
		if current == like {
			_, err := db.Exec("DELETE FROM link_comment WHERE user_id = ? AND comment_id = ?", id_user, id_comment)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec("UPDATE link_comment SET likes = ? WHERE user_id = ? AND comment_id = ?", like, id_user, id_comment)
			if err != nil {
				return err
			}
		}
	} else {
		_, err := db.Exec("INSERT INTO link_comment(user_id,comment_id,likes) VALUES(?,?,?)", id_user, id_comment, like)
		if err != nil {
			return err
		}
	}
	return nil
}

func FlagComment(db *sql.DB, r *http.Request) error {
	id, _ := strconv.Atoi(r.FormValue("report"))
	_, err := db.Exec("UPDATE comment SET flagged = ? WHERE id = ?", 1, id)
	if err != nil {
		return err
	}
	return nil
}
