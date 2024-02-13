package forum

type Dataprofile struct {
	Likedcomment    []Commentpost
	Dislikedcomment []Commentpost

	Likedpost    []Post
	Dislikedpost []Post

	Notified []Notified
}

type Commentpost struct {
	Post    Post
	Comment Comment
}

type Notified struct {
	Post        Post
	Tab_comment []Comment
	Tabliked    []Like
}

type Like struct {
	Username string
	Is_Liked bool
}
