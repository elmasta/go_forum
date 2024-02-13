package forum

import (
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Create a new entry in the database with the given information about the new user
func (app *App_db) CreateUser(user *models.User) error {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = app.DB.Exec(
		"INSERT INTO users(username, user_type_id, pwd, email, valid, creation) VALUES (?,?,?,?,?,?)",
		user.Username,
		user.UserType,
		string(hashPass),
		user.Email,
		1,
		time.Now(),
	)
	return err
}

// Register new user in application
func (app *App_db) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		"web/templates/register.html",
		"web/templates/head.html",
		"web/templates/footer.html",
	)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	errMsg := r.URL.Query().Get("error")
	if r.Method == "POST" {
		user := &models.User{}
		if err := middle.CheckRegister(app.DB, r, user); err != nil {
			if err.Error() == "username or email already exist" {
				errMsg = "Username or Email already exist!"
			}
			if err.Error() == "passwords do not match" {
				errMsg = "Passwords do not match!"
			}
			if err.Error() == "passwords parsing error" {
				errMsg = "password length must be at least 8 characters long and can only contain alphanumerical characters!"
			}
			if err.Error() == "email not valid" {
				errMsg = "Email is not valid!"
			}
			http.Redirect(w, r, "/register?error="+url.QueryEscape(errMsg), http.StatusFound)
			return
		}
		if err := app.CreateUser(user); err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
	if err := template.Execute(w, map[string]string{"ErrorMessage": errMsg}); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}
