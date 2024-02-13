package forum

import (
	"net/http"
	"time"

	s "forum/sessions"
)

// Disconnect user, remove session and clear cookies before redirecting to home page
func (app *App_db) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	sessionToken := cookie.Value
	delete(s.GlobalSessions, sessionToken)

	stmt, err := app.DB.Prepare("UPDATE users SET session_token = NULL WHERE session_token = ?")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
	}
	_, err = stmt.Exec(sessionToken)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
