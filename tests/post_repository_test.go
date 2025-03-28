package service_test

import (
	"1337bo4rd/internal/domain"
	"testing"
	"time"
)

func TestFindUniquePostID(t *testing.T) {
	mockRepo := NewMockPostRepo()
	id, err := mockRepo.FindUniquePostID()
	if err != nil {
		t.Errorf("❌ FindUniquePostID() returned error: %v", err)
	}
	if id != 1 {
		t.Errorf("❌ Expected ID 1, got %d", id)
	}
}

func TestSavePost(t *testing.T) {
	mockRepo := NewMockPostRepo()
	post := domain.Post{ID: 1, Title: "Test Post", Content: "Content", ExpiresAt: time.Now().Add(time.Hour)}

	err := mockRepo.SavePost(post)
	if err != nil {
		t.Errorf("❌ SavePost() returned error: %v", err)
	}

	if _, exists := mockRepo.posts[1]; !exists {
		t.Errorf("❌ Post not saved")
	}
}

func TestGetActivePosts(t *testing.T) {
	mockRepo := NewMockPostRepo()
	mockRepo.SavePost(domain.Post{ID: 1, ExpiresAt: time.Now().Add(time.Hour)})  // Active
	mockRepo.SavePost(domain.Post{ID: 2, ExpiresAt: time.Now().Add(-time.Hour)}) // Expired

	posts, err := mockRepo.GetActivePosts()
	if err != nil {
		t.Errorf("❌ GetActivePosts() returned error: %v", err)
	}
	if len(posts) != 1 {
		t.Errorf("❌ Expected 1 active post, got %d", len(posts))
	}
}

func TestGetArchivePosts(t *testing.T) {
	mockRepo := NewMockPostRepo()
	mockRepo.SavePost(domain.Post{ID: 1, ExpiresAt: time.Now().Add(time.Hour)})  // Active
	mockRepo.SavePost(domain.Post{ID: 2, ExpiresAt: time.Now().Add(-time.Hour)}) // Expired

	posts, err := mockRepo.GetArchivePosts()
	if err != nil {
		t.Errorf("❌ GetArchivePosts() returned error: %v", err)
	}
	if len(posts) != 1 {
		t.Errorf("❌ Expected 1 archived post, got %d", len(posts))
	}
}

func TestGetPost(t *testing.T) {
	mockRepo := NewMockPostRepo()
	mockRepo.SavePost(domain.Post{ID: 1, Title: "Test Post"})

	post, err := mockRepo.GetPost(1)
	if err != nil {
		t.Errorf("❌ GetPost() returned error: %v", err)
	}
	if post.Title != "Test Post" {
		t.Errorf("❌ Expected title 'Test Post', got '%s'", post.Title)
	}

	_, err = mockRepo.GetPost(99) // Non-existent ID
	if err == nil {
		t.Errorf("❌ Expected error for non-existent post, got nil")
	}
}
