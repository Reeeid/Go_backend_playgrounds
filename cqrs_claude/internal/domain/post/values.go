package post

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

type PostID struct{ value uuid.UUID }
type Title struct{ value string }
type Content struct{ value string }
type AuthorID struct{ value uuid.UUID }
type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

func NewPostID() (PostID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return PostID{}, err
	}
	return PostID{value: id}, nil
}

func NewTitle(s string) (Title, error) {
	l := utf8.RuneCountInString(s)
	if l == 0 {
		return Title{}, errors.New("title is required")
	}
	if l > 100 {
		return Title{}, errors.New("title must be 100 chars or less")
	}
	return Title{value: s}, nil
}

func NewContent(s string) (Content, error) {
	if utf8.RuneCountInString(s) == 0 {
		return Content{}, errors.New("content is required")
	}
	return Content{value: s}, nil
}

func NewAuthorID(id uuid.UUID) (AuthorID, error) {
	if id == uuid.Nil {
		return AuthorID{}, errors.New("authorID cannot be nil")
	}
	return AuthorID{value: id}, nil
}

func NewStatus(s string) (Status, error) {
	switch Status(s) {
	case StatusDraft, StatusPublished, StatusArchived:
		return Status(s), nil
	}
	return "", errors.New("invalid status: " + s)
}

func NewPostIDFromUUID(id uuid.UUID) (PostID, error) {
	if id == uuid.Nil {
		return PostID{}, errors.New("postID cannot be nil")
	}
	return PostID{value: id}, nil
}

func (p PostID) String() string   { return p.value.String() }
func (t Title) String() string    { return t.value }
func (c Content) String() string  { return c.value }
func (a AuthorID) String() string { return a.value.String() }
