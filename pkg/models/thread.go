package forum

import (
	"time"
)

type Thread struct {
	ID           int64     `json:"id"`
	Author       string    `json:"author"`
	Post         int64     `json:"category"`
	Content      string    `json:"content"`
	Like         int       `json:"like"`
	Dislike      int       `json:"dislike"`
	CreationDate time.Time `json:"creation_date"`
}
