package query

import (
	"errors"
	"sync"
)

// PostReadModel はクエリ用の非正規化モデル（DBのビューに相当）
type PostReadModel struct {
	ID       string
	Title    string
	Content  string
	AuthorID string
	Status   string
}

// PostReadStore はRead Model のインメモリストア
type PostReadStore struct {
	mu    sync.RWMutex
	posts map[string]PostReadModel
}

func NewPostReadStore() *PostReadStore {
	return &PostReadStore{posts: make(map[string]PostReadModel)}
}

func (s *PostReadStore) Save(p PostReadModel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.posts[p.ID] = p
}

func (s *PostReadStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.posts, id)
}

func (s *PostReadStore) GetByID(id string) (PostReadModel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.posts[id]
	if !ok {
		return PostReadModel{}, errors.New("post not found: " + id)
	}
	return p, nil
}

func (s *PostReadStore) GetAll() []PostReadModel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	posts := make([]PostReadModel, 0, len(s.posts))
	for _, p := range s.posts {
		posts = append(posts, p)
	}
	return posts
}
