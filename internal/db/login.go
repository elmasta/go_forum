package forum

import (
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"net/http"
	"net/url"
	"text/template"
)

// Connect existing user to forum by comparing their credentials with database
func (app *App_db) LoginHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		"web/templates/login.html",
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
		if err := middle.Auth(app.DB, w, r, user); err != nil {
			if err.Error() == "invalid email" {
				errMsg = "Invalid email address"
			}
			if err.Error() == "invalid password" {
				errMsg = "Invalid password"
			}
			http.Redirect(w, r, "/login?error="+url.QueryEscape(errMsg), http.StatusFound)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if err := template.Execute(w, map[string]string{"ErrorMessage": errMsg}); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
	}
}
