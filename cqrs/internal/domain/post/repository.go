package post

type PostCommandRepository interface {
	Save(post *Post) error
	Update(post *Post) error
	Delete(id PostID) error
}

type PostQueryRepository interface {
	GetByID(id PostID) (*Post, error)
	GetByAuthorID(authorID AuthorID) ([]*Post, error)
	GetByStatus(status Status) ([]*Post, error)
	GetAll() ([]*Post, error)
}
