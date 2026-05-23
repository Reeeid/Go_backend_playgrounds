package command

type CreatePostCommand struct {
	Title    string
	Content  string
	AuthorID string
}

type UpdatePostCommand struct {
	ID      string
	Title   string
	Content string
	Status  string
}

type DeletePostCommand struct {
	ID string
}
