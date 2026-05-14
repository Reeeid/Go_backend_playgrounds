package post

import "time"

type Post struct {
	id        PostID
	title     Title
	content   Content
	authorID  AuthorID
	status    Status
	createdAt time.Time
	updatedAt time.Time
}

func NewPost(postID PostID, title Title, content Content, authorID AuthorID, status Status) *Post {
	return &Post{
		id:        postID,
		title:     title,
		content:   content,
		authorID:  authorID,
		status:    status,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}

func (p *Post) Update(title Title, content Content, status Status) {
	p.title = title
	p.content = content
	p.status = status
	p.updatedAt = time.Now()
}

func (p *Post) Delete() {
	//command HandlerでPostを削除するためのメソッド
}
