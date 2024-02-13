package forum

import (
	"database/sql"
	"net/http"
	"time"

	models "forum/pkg/models"

	"github.com/gofrs/uuid"
)

// Set Token and send it to the server session and user cookie
func SetToken(db *sql.DB, w http.ResponseWriter, r *http.Request, user *models.User) error {
	sessionToken, _ := uuid.NewV4()
	expiresAt := time.Now().Add(3600 * time.Second)

	var mod, light_mod, admin bool

	if user.UserType == 2 {
		mod = true
	} else if user.UserType == 3 {
		admin = true
	} else if user.UserType == 4 {
		light_mod = true
	}

	GlobalSessions[sessionToken.String()] = Session{
		Username:  user.Username,
		UserID:    user.ID,
		Moderator: mod,
		Admin:     admin,
		ModLight:  light_mod,
		EndLife:   expiresAt,
	}

	stmt, err := db.Prepare("UPDATE users SET session_token = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sessionToken.String(), user.ID)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: expiresAt,
	})
	return nil
}
