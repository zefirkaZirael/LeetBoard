package service_test

import (
	"1337bo4rd/internal/domain"
	"errors"
	"time"
)

// MockPostRepo (mock for PostRepository)
type MockPostRepo struct {
	posts   map[int]domain.Post
	nextID  int
	isError bool
}

func NewMockPostRepo() *MockPostRepo {
	return &MockPostRepo{
		posts:  make(map[int]domain.Post),
		nextID: 1,
	}
}

func (m *MockPostRepo) FindUniquePostID() (int, error) {
	if m.isError {
		return 0, errors.New("DB error")
	}
	id := m.nextID
	m.nextID++
	return id, nil
}

func (m *MockPostRepo) SavePost(post domain.Post) error {
	if m.isError {
		return errors.New("DB error")
	}
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepo) GetActivePosts() ([]domain.Post, error) {
	if m.isError {
		return nil, errors.New("DB error")
	}
	var active []domain.Post
	for _, p := range m.posts {
		if time.Now().Before(p.ExpiresAt) {
			active = append(active, p)
		}
	}
	return active, nil
}

func (m *MockPostRepo) GetArchivePosts() ([]domain.Post, error) {
	if m.isError {
		return nil, errors.New("DB error")
	}
	var archive []domain.Post
	for _, p := range m.posts {
		if time.Now().After(p.ExpiresAt) {
			archive = append(archive, p)
		}
	}
	return archive, nil
}

func (m *MockPostRepo) GetPost(id int) (domain.Post, error) {
	if m.isError {
		return domain.Post{}, errors.New("DB error")
	}
	post, exists := m.posts[id]
	if !exists {
		return domain.Post{}, errors.New("post not found")
	}
	return post, nil
}

func (m *MockPostRepo) ArchiveExpiredPosts() error {
	// Simulate archiving posts (do nothing or remove expired posts)
	return nil
}

func (m *MockPostRepo) GetArchivePost(id int) (domain.Post, error) {
	post, exists := m.posts[id]
	if !exists {
		return domain.Post{}, errors.New("archive post not found")
	}
	return post, nil
}

func (m *MockPostRepo) IsPostExist(id int) (bool, error) {
	_, exists := m.posts[id]
	return exists, nil
}

// MockSessionRepo (мок для SessionRepository)
type MockSessionRepo struct {
	users          map[string]domain.User
	expiredSession bool
}

func NewMockSessionRepo() *MockSessionRepo {
	return &MockSessionRepo{
		users: make(map[string]domain.User),
	}
}

func (m *MockSessionRepo) FindUniqueUserID() (int, error) {
	return 1, nil
}

func (m *MockSessionRepo) SaveUser(user domain.User) error {
	m.users[user.Token_ID] = user
	return nil
}

func (m *MockSessionRepo) GetUser(sessionID string) (domain.User, error) {
	user, exists := m.users[sessionID]
	if !exists {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func (m *MockSessionRepo) DeleteExpiredSessions() error {
	if m.expiredSession {
		return errors.New("failed to delete sessions")
	}
	return nil
}

func (m *MockSessionRepo) IsSessionExist(sessionID string) (bool, error) {
	_, exists := m.users[sessionID]
	return exists, nil
}

func (m *MockSessionRepo) GetUserByID(userID int) (domain.User, error) {
	for _, user := range m.users {
		if user.ID == userID {
			return user, nil
		}
	}
	return domain.User{}, errors.New("user not found")
}

// Unused methods
func (m *MockSessionRepo) ChangeUserName(changed_name string, user_id int) error { return nil }

func (m *MockSessionRepo) IsNameEqual(changed_name string, user_id int) (bool, error) {
	return true, nil
}

// MockCommentRepo (mock for CommentRepository)
type MockCommentRepo struct {
	comments map[int][]domain.Comment
	isError  bool
}

func NewMockCommentRepo() *MockCommentRepo {
	return &MockCommentRepo{
		comments: make(map[int][]domain.Comment),
	}
}

var _ domain.CommentRepoInt = (*MockCommentRepo)(nil)

func (m *MockCommentRepo) IsReplyIdExist(post_id int, reply_id int) (bool, error) {
	return false, nil
}

func (m *MockCommentRepo) CreateComment(comment domain.Comment) error {
	if m.isError {
		return errors.New("DB error")
	}
	m.comments[comment.PostID] = append(m.comments[comment.PostID], comment)
	return nil
}

func (m *MockCommentRepo) GetCommentsByPost(postID int) ([]domain.Comment, error) {
	if m.isError {
		return nil, errors.New("DB error")
	}
	return m.comments[postID], nil
}

func (m *MockCommentRepo) FindUniqueCommentID() (int, error) {
	return len(m.comments) + 1, nil
}

func (m *MockPostRepo) UpdatePostExpiration(postID int, newExpiration time.Time) error {
	if m.isError {
		return errors.New("DB error")
	}
	post, exists := m.posts[postID]
	if !exists {
		return errors.New("post not found")
	}
	post.ExpiresAt = newExpiration
	m.posts[postID] = post
	return nil
}
