package post

import (
	"unicode/utf8"

	"github.com/google/uuid"
)

type PostID struct {
	value uuid.UUID
}

type Title struct {
	value string
}

type Content struct {
	value string
}

type AuthorID struct {
	value uuid.UUID
}
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

func NewTitle(value string) (Title, error) {
	length := utf8.RuneCountInString(value)
	if length == 0 {
		return Title{}, &NilError{FieldName: "title"}
	} else if length > 32 {
		return Title{}, &LengthError{FieldName: "title", Min: 1, Max: 32}
	}
	return Title{value: value}, nil
}

func NewContent(value string) (Content, error) {
	length := utf8.RuneCountInString(value)
	if length == 0 {
		return Content{}, &NilError{FieldName: "content"}
	} else if length > 1000 {
		return Content{}, &LengthError{FieldName: "content", Min: 1, Max: 1000}
	}
	return Content{value: value}, nil
}

func newAuthorID(value uuid.UUID) (AuthorID, error) {
	if value == uuid.Nil {
		return AuthorID{}, &NilError{FieldName: "authorID"}
	}
	return AuthorID{value: value}, nil
}

func NewStatus(value string) (status Status, err error) {
	switch value {
	case string(StatusDraft):
		status = StatusDraft
	case string(StatusPublished):
		status = StatusPublished
	case string(StatusArchived):
		status = StatusArchived
	default:
		err = &StatusError{FieldName: "status"}
	}
	return
}

func NewPostIDFromString(s string) (PostID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return PostID{}, err
	}
	return PostID{value: id}, nil
}
