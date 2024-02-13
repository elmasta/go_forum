package forum

type UserType struct {
	ID      int64   `json:"id"`
	Rank	int64	`json:"rank"`
	Label	string	`json:"label"`
}