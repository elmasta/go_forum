package forum

import (
	"time"
)

type Categories struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"category"`
	CreationDate time.Time `json:"time"`
	IfCurrentCat bool
}
