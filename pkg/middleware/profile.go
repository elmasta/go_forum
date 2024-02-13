package forum

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	models "forum/pkg/models"
)

func FetchUser(db *sql.DB, cookie string) models.User {
	var currentUser models.User
	err := db.QueryRow("SELECT id, user_type_id, username, email, valid, asked_mod, creation FROM users WHERE session_token=?", cookie).Scan(
		&currentUser.ID,
		&currentUser.UserType,
		&currentUser.Username,
		&currentUser.Email,
		&currentUser.Validation,
		&currentUser.AskedMod,
		&currentUser.CreationDate)
	if err != nil {
		log.Fatal(err)
	}
	return currentUser
}

func AskModerator(db *sql.DB, r *http.Request, asked int, id int) error {
	_, err := db.Exec("UPDATE users SET asked_mod=? WHERE id=?", asked, id)
	if err != nil {
		return err
	}
	return nil
}

func Likedpost(db *sql.DB, user_id int64) []models.Post {
	rows, err := db.Query(" SELECT post.id,post.title,post.creation FROM post INNER JOIN link_post ON post.id = link_post.post_id WHERE user_id = ? AND likes = 1;", user_id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var tabpost []models.Post
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		rows.Scan(
			&post.ID,
			&post.Title,
			&post.CreationDate)
		tabpost = append(tabpost, post)
	}
	// fmt.Println(tabpost)
	return tabpost
}

func Dislikedpost(db *sql.DB, user_id int64) []models.Post {
	rows, err := db.Query(" SELECT post.id,post.title,post.creation FROM post INNER JOIN link_post ON post.id = link_post.post_id WHERE user_id = ? AND likes = 0;", user_id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var tabpost []models.Post
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		rows.Scan(
			&post.ID,
			&post.Title,
			&post.CreationDate)
		tabpost = append(tabpost, post)
	}
	// fmt.Println(tabpost)
	return tabpost
}

func Likedcomment(db *sql.DB, user_id int64) []models.Commentpost {
	rows, err := db.Query(` SELECT post.id,post.title,comment.id,comment.content,comment.creation FROM comment 
							INNER JOIN link_comment ON comment.id = link_comment.comment_id
							INNER JOIN post ON comment.post_id = post.id
							WHERE link_comment.user_id = ? AND link_comment.likes = 1;`, user_id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	var tab []models.Commentpost
	for rows.Next() {
		var commentpost models.Commentpost
		rows.Scan(
			&commentpost.Post.ID,
			&commentpost.Post.Title,
			&commentpost.Comment.ID,
			&commentpost.Comment.Content,
			&commentpost.Comment.CreationDate)
		tab = append(tab, commentpost)
	}
	// fmt.Println(tab)
	return tab
}

func Dislikedcomment(db *sql.DB, user_id int64) []models.Commentpost {
	rows, err := db.Query(` SELECT post.id,post.title,comment.id,comment.content,comment.creation FROM comment 
							INNER JOIN link_comment ON comment.id = link_comment.comment_id
							INNER JOIN post ON comment.post_id = post.id
							WHERE link_comment.user_id = ? AND link_comment.likes = 0;`, user_id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	var tab []models.Commentpost
	for rows.Next() {
		var commentpost models.Commentpost
		rows.Scan(
			&commentpost.Post.ID,
			&commentpost.Post.Title,
			&commentpost.Comment.ID,
			&commentpost.Comment.Content,
			&commentpost.Comment.CreationDate)
		tab = append(tab, commentpost)
	}
	// fmt.Println(tab)
	return tab
}

func ProfilNotified(db *sql.DB, user_id int64) []models.Notified {
	rows, err := db.Query("Select post.id,post.title From post Where author_id = ? ORDER BY id DESC LIMIT 5", user_id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var tab []models.Notified
	for rows.Next() {
		var notified models.Notified
		rows.Scan(
			&notified.Post.ID,
			&notified.Post.Title)
		tab = append(tab, notified)
	}
	rows.Close()
	for index := range tab {
		var tabcomment []models.Comment
		var tabliked []models.Like
		rows, err := db.Query(`	Select comment.content,comment.author_id,users.username From comment
								INNER JOIN	users ON users.id = comment.author_id
		 						Where post_id = ? ORDER BY comment.id DESC LIMIT 5`, tab[index].Post.ID)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		for rows.Next() {
			var comment models.Comment
			rows.Scan(
				&comment.Content,
				&comment.AuthorID,
				&comment.Author)
			tabcomment = append(tabcomment, comment)
		}
		rows.Close()
		rows, err = db.Query(`	Select link_post.likes,users.username From link_post
								INNER JOIN	users ON users.id = link_post.user_id
		 						Where post_id = ? ORDER BY link_post.id DESC LIMIT 5`, tab[index].Post.ID)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		for rows.Next() {
			var like models.Like
			rows.Scan(
				&like.Is_Liked,
				&like.Username)
			tabliked = append(tabliked, like)
		}
		rows.Close()
		tab[index].Tab_comment = tabcomment
		tab[index].Tabliked = tabliked
	}
	// fmt.Println(tab)
	// for _, data := range tab {
	// 	fmt.Println("Post:")
	// 	fmt.Println(data.Post)
	// 	fmt.Println("Tab comment:")
	// 	for _, comment := range data.Tab_comment {
	// 		fmt.Println("	", comment)
	// 	}
	// 	fmt.Println("Tab like:")
	// 	for _, like := range data.Tabliked {
	// 		fmt.Println("	", like)
	// 	}
	// }
	return tab
}
