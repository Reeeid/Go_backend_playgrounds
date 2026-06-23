package post

type PostQueryRepository interface {
	GetByID(id PostID) (*Post, error)
	GetByAuthorID(authorID AuthorID) ([]*Post, error)
	GetByStatus(status Status) ([]*Post, error)
	GetAll() ([]*Post, error)
}
