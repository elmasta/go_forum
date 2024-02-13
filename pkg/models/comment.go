package forum

type Comment struct {
	ID            int64  `json:"id"`
	AuthorID      int64  `json:"author_id"`
	Author        string `json:"author"`
	IfCurrentUser bool
	Img           string
	PostID        int64 `json:"post_id"`
	Post          string
	Content       string `json:"content"`
	Like          int
	User_like     map[int64]bool
	Dislike       int
	User_dislike  map[int64]bool
	Flagged       int
	CreationDate  string `json:"time"`
}
