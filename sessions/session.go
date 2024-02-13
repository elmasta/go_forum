package forum

import (
	"database/sql"
	"net/http"
	"time"
)

var GlobalSessions = map[string]Session{}

type Session struct {
	Username  string
	UserID    int64
	Admin     bool
	Moderator bool
	ModLight  bool
	EndLife   time.Time
}

// Check if session token is the same as user token to avoid multiple instances of session conflict
func CheckSession(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	var userToken string
	c, err := r.Cookie("session_token")
	if err != nil {
		return err
	}
	session_token := c.Value
	session := GlobalSessions[session_token]
	err = db.QueryRow("SELECT session_token FROM users where id = ?", session.UserID).Scan(&userToken)
	if err != nil {
		return err
	}
	if userToken != session_token {
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}
	return nil
}

// Clear the session token to prevent memory leaks
func CheckActive() {
	for k := range GlobalSessions {
		if GlobalSessions[k].EndLife.Before(time.Now()) {
			delete(GlobalSessions, k)
		}
	}
}
