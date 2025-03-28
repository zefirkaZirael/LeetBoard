package service_test

import (
	"1337bo4rd/internal/domain"
	"1337bo4rd/internal/infrastructure/external"
	"1337bo4rd/internal/service"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestCreateSession(t *testing.T) {
	mockRepo := &MockSessionRepo{users: make(map[string]domain.User)}
	sessionService := service.NewUserSessionService(mockRepo, external.DefaultExternalAPI())

	sessionID := "test_session"
	userID, err := sessionService.CreateSession(sessionID)
	if err != nil {
		t.Errorf("❌ CreateSession() returned an error: %v", err)
	}
	if userID != 1 {
		t.Errorf("❌ Expected user ID 1, got %d", userID)
	}
}

func TestIsValidSession(t *testing.T) {
	mockRepo := &MockSessionRepo{users: make(map[string]domain.User)}
	sessionService := service.NewUserSessionService(mockRepo, external.DefaultExternalAPI())

	mockRepo.users["valid_session"] = domain.User{}

	code, err := sessionService.IsValidSession("valid_session")
	if err != nil || code != http.StatusOK {
		t.Errorf("❌ Expected valid session, got error: %v", err)
	}

	code, err = sessionService.IsValidSession("invalid_session")
	if err == nil || code != http.StatusUnauthorized {
		t.Errorf("❌ Expected unauthorized error, got %v", err)
	}
}

func TestDeleteExpiredSessions(t *testing.T) {
	mockRepo := &MockSessionRepo{}
	sessionService := service.NewUserSessionService(mockRepo, external.DefaultExternalAPI())

	err := sessionService.DeleteExpiredSessions()
	if err != nil {
		t.Errorf("❌ DeleteExpiredSessions() returned an error: %v", err)
	}
}

func TestGetSession(t *testing.T) {
	mockRepo := &MockSessionRepo{users: make(map[string]domain.User)}
	sessionService := service.NewUserSessionService(mockRepo, external.DefaultExternalAPI())

	mockRepo.users["test_session"] = domain.User{
		ID:         1,
		Token_ID:   "test_session",
		Expires_at: time.Now().Add(1 * time.Hour),
	}

	sessionService.Sessions["test_session"] = domain.UserSession{
		ID:      strconv.Itoa(1),
		Expires: time.Now().Add(1 * time.Hour),
	}
	session, err := sessionService.GetSession("test_session")
	if err != nil {
		t.Errorf("❌ GetSession() returned an error: %v", err)
	}
	if session.ID != strconv.Itoa(1) {
		t.Errorf("❌ Expected user ID 1, got %s", session.ID)
	}

	_, err = sessionService.GetSession("invalid_session")
	if err == nil {
		t.Errorf("❌ Expected error for non-existent session, got nil")
	}
}
