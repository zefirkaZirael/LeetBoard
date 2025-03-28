package service

import (
	"1337bo4rd/internal/domain"
	"errors"
	"net/http"
	"time"
)

// UserSessionService struct
type UserSessionService struct {
	Sessions map[string]domain.UserSession
	Repo     domain.SessionRepoInt
	External domain.ExternalAPI
}

// Constructor
func NewUserSessionService(repo domain.SessionRepoInt, external domain.ExternalAPI) *UserSessionService {
	return &UserSessionService{Sessions: make(map[string]domain.UserSession), Repo: repo, External: external}
}

var _ domain.UserSessionService = (*UserSessionService)(nil)

func (s *UserSessionService) GetUser(sessionID string) (domain.User, error) {
	return s.Repo.GetUser(sessionID)
}

// Create a new session
func (s *UserSessionService) CreateSession(session_id string) (int, error) {
	var newUser domain.User
	var err error

	newUser.ID, err = s.Repo.FindUniqueUserID()
	if err != nil {
		return 0, err
	}
	newUser.Token_ID = session_id

	avatarCount, err := s.External.GetAvatarCount()
	if err != nil {
		return 0, err
	}
	oldId := newUser.ID
	if newUser.ID > avatarCount {
		newUser.ID %= avatarCount
	}
	err = s.External.GetCharacter(&newUser)
	if err != nil {
		return 0, err
	}
	newUser.ID = oldId
	newUser.TokenDate = time.Now()
	newUser.Expires_at = time.Now().Add(7 * 24 * time.Hour)
	err = s.Repo.SaveUser(newUser)
	if err != nil {
		return 0, err
	}
	return newUser.ID, nil
}

func (s *UserSessionService) DeleteExpiredSessions() error {
	err := s.Repo.DeleteExpiredSessions()
	if err != nil {
		return err
	}
	return nil
}

// Retrieve a session
func (s *UserSessionService) GetSession(id string) (domain.UserSession, error) {
	session, exists := s.Sessions[id]
	if !exists || time.Now().After(session.Expires) {
		return domain.UserSession{}, errors.New("session not found or expired")
	}
	return session, nil
}

func (s *UserSessionService) IsValidSession(session_id string) (int, error) {
	exist, err := s.Repo.IsSessionExist(session_id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusUnauthorized, errors.New("user is not authorized")
	}
	return http.StatusOK, nil
}
