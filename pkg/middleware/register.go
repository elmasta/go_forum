package forum

import (
	"database/sql"
	"errors"
	models "forum/pkg/models"
	"net/http"
	"net/mail"
	"regexp"
)

func isAlphanumeric(str string) bool {
	var alphanumeric = regexp.MustCompile("^[a-zA-Z0-9]*$")
	return alphanumeric.MatchString(str)
}

// Prevent duplicate credentials in database during register procedure
func CheckRegister(db *sql.DB, r *http.Request, user *models.User) error {
	username, email := r.FormValue("username"), r.FormValue("email")
	password := r.FormValue("password")
	if len(password) < 8 || !isAlphanumeric(password) {
		return errors.New("passwords parsing error")
	} else if password != r.FormValue("confirmation") {
		return errors.New("passwords do not match")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("email not valid")
	}

	err = db.QueryRow(
		"SELECT username,email FROM users WHERE username=? OR email=?",
		username,
		email).Scan(&user.Username, &user.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		return errors.New("username or email already exist")
	}
	user.Username, user.Email, user.Password, user.UserType = username, email, password, 1
	return nil
}
