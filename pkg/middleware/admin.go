package forum

import (
	"database/sql"
	"errors"
	models "forum/pkg/models"
	"log"
	"net/http"
	"strconv"
)

// Prevent duplicate credentials in database during register procedure
func CheckAdminRegister(db *sql.DB, confirmation string, user *models.User) error {
	if confirmation != user.Password {
		return errors.New("passwords do not match")
	}

	err := db.QueryRow(
		"SELECT username,email FROM users WHERE username=? OR email=?",
		user.Username,
		user.Email).Scan(&user.Username, &user.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		return errors.New("username or email already exist")
	}
	return nil
}

func RmUser(db *sql.DB, r *http.Request) error {
	id, err := strconv.Atoi(r.FormValue("deletion"))
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}

func AddMod(db *sql.DB, r *http.Request, to_add int, id int) error {
	_, err := db.Exec("UPDATE users SET user_type_id=?, asked_mod=? WHERE id=?", to_add, 0, id)
	if err != nil {
		return err
	}
	return nil
}

func DelMod(db *sql.DB, r *http.Request) error {
	id, err := strconv.Atoi(r.FormValue("del_mod"))
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE users SET user_type_id=? WHERE id=?", 1, id)
	if err != nil {
		return err
	}
	return nil
}

func FetchUsers(db *sql.DB) []models.User {
	rows, err := db.Query("SELECT id, user_type_id, username, email, valid, asked_mod, creation FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var users_list []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID,
			&user.UserType,
			&user.Username,
			&user.Email,
			&user.Validation,
			&user.AskedMod,
			&user.CreationDate)
		if err != nil {
			log.Fatal(err)
		}
		switch user.UserType {
		case 1:
			user.UserTypeTxt = "User"
		case 2:
			user.UserTypeTxt = "Moderator"
		case 3:
			user.UserTypeTxt = "Admin"
		case 4:
			user.UserTypeTxt = "Comment Moderator"
		}
		user.FormattedTime = user.CreationDate.Format("2006-01-02 15:04:05")
		users_list = append(users_list, user)
	}
	return users_list
}
