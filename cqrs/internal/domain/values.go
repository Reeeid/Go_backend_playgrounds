package domain

type PostID string
type Title string
type Content string
type AuthorID string
type Status string //draft|published|archived

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

func NewTitle(title string) (Title, error) {
	if len(title) == 0 {
		return "", ErrInvalidTitle
	}
}
