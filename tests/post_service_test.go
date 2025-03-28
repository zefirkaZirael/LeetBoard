package service_test

import (
	"1337bo4rd/internal/domain"
	"1337bo4rd/internal/infrastructure/s3"
	"1337bo4rd/internal/service"
	"database/sql"
	"net/http"
	"testing"
)

func TestCreatePost(t *testing.T) {
	mockPostRepo := NewMockPostRepo()
	mockSessionRepo := NewMockSessionRepo()
	mockS3 := &s3.S3{}
	mockDB := &sql.DB{}

	postService := service.NewPostService(mockPostRepo, mockSessionRepo, mockS3, mockDB)

	post := domain.Post{
		AuthorID:     1,
		Title:        "Test Title",
		Content:      "Test Content",
		Changed_name: "New Name",
		File:         nil,
	}

	// Call method
	status, err := postService.CreatePost(post)
	// We check that there is no error and the status is 200
	if err != nil {
		t.Errorf("❌ CreatePost() вернул ошибку: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("❌ Ожидался статус %d, но получен %d", http.StatusOK, status)
	}
}
