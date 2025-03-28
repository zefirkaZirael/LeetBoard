package service_test

import (
	"1337bo4rd/internal/domain"
	"1337bo4rd/internal/infrastructure/s3"
	"1337bo4rd/internal/service"
	"testing"
	"time"
)

func TestCreateComment(t *testing.T) {
	mockCommentRepo := NewMockCommentRepo()
	mockPostRepo := NewMockPostRepo()
	mockSessionRepo := NewMockSessionRepo()
	mockStorage := s3.NewS3Repo()

	mockSessionRepo.users["session_id"] = domain.User{
		ID:       1,
		Name:     "Test User",
		ImageURL: "http://avatar.com/avatar.jpg",
	}

	mockPostRepo.posts[1] = domain.Post{
		ID:        1,
		Title:     "Test Post",
		Content:   "Some content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	commentService := service.NewCommentService(mockCommentRepo, mockPostRepo, mockSessionRepo, mockStorage)

	commentID, err := commentService.CreateComment(1, 0, "Test Comment", "session_id", nil)
	if err != nil {
		t.Errorf("❌ CreateComment() returned error: %v", err)
	}
	if commentID != 1 {
		t.Errorf("❌ Expected comment ID 1, got %d", commentID)
	}
}

func TestGetCommentsByPost(t *testing.T) {
	mockCommentRepo := NewMockCommentRepo()
	mockPostRepo := NewMockPostRepo()
	mockSessionRepo := NewMockSessionRepo()
	mockStorage := s3.NewS3Repo()

	mockSessionRepo.users["session_id"] = domain.User{
		ID:       1,
		Name:     "Test User",
		ImageURL: "http://avatar.com/avatar.jpg",
	}

	mockPostRepo.posts[1] = domain.Post{
		ID:        1,
		Title:     "Test Post",
		Content:   "Some content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	commentService := service.NewCommentService(mockCommentRepo, mockPostRepo, mockSessionRepo, mockStorage)
	// Creating 2 comments for post number 1
	commentService.CreateComment(1, 0, "Comment 1", "session_id", nil)
	commentService.CreateComment(1, 0, "Comment 2", "session_id", nil)

	comments, err := commentService.GetCommentsByPost(1)
	if err != nil {
		t.Errorf("❌ GetCommentsByPost() returned error: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("❌ Expected 2 comments, got %d", len(comments))
	}
}
