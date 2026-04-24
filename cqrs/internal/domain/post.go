package domain

import "time"

type Post struct {
	ID        PostID
	Title     Title
	Content   Content
	AuthorID  AuthorID
	CreatedAt time.Time
	UpdatedAt time.Time
}
