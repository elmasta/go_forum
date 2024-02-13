package forum

import (
	"time"
)

type User struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	UserType      int64     `json:"user_type"` // foreign key to UserType table
	UserTypeTxt   string    `json:"user_type_txt"`
	Validation    int64     `json:"validation"` //0 false 1 true
	AskedMod      int       `json:"asked_mod"`
	CreationDate  time.Time `json:"creation_time"`
	FormattedTime string    `json:"formatted_time"`
}
