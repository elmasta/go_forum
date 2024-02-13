package forum

type Post struct {
	ID             int64  `json:"id"`
	AuthorID       int64  `json:"author_id"`
	Author         string `json:"author"`
	IfCurrentUser  bool
	Categories     []int
	CategoriesName []string
	Img            string
	Ifimg          bool
	Title          string `json:"title"`
	Content        string `json:"content"`
	Like           int
	User_like      map[int64]bool
	Dislike        int
	User_dislike   map[int64]bool
	CreationDate   string `json:"time"`
	Flagged        int    `json:"flagged"`
	Tab_comment    []Comment
}
