package middleware

import (
	"1337bo4rd/internal/domain"
	"1337bo4rd/internal/infrastructure/repository"
	"1337bo4rd/utils"
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

type SessionHandler struct {
	serv domain.UserSessionService
	DB   *sql.DB
}

func NewSessionHandler(serv domain.UserSessionService, db *sql.DB) *SessionHandler {
	return &SessionHandler{serv: serv, DB: db}
}

func (h *SessionHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	_, err := h.Authorization(w, r)
	if err != nil {
		slog.Error("❌ Authorization error: " + err.Error())
		utils.ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	temp, err := template.ParseFiles("web/templates/catalog.html")
	if err != nil {
		slog.Error("❌ Failed to parse template file: " + err.Error())
		utils.ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	posts, err := repository.NewPostRepository(h.DB).GetActivePosts()
	if err != nil {
		slog.Error("❌ Error on GetPost: " + err.Error())
		utils.ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	err = temp.Execute(w, posts)
	if err != nil {
		slog.Error("❌ Failed to execute template: " + err.Error())
		utils.ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
}

// Checks user authorization and returns his ID
func (h *SessionHandler) Authorization(w http.ResponseWriter, r *http.Request) (int, error) {
	var (
		sessionID string
		userID    int
		err       error
	)
	err = h.serv.DeleteExpiredSessions()
	if err != nil {
		slog.Error("❌ Failed to delete expires sessions in system: " + err.Error())
		return 0, err
	}
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie == nil {
		sessionID = utils.GenerateSessionID()
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			Expires:  time.Now().Add(24 * 7 * time.Hour),
			HttpOnly: true,
		})
		userID, err = h.serv.CreateSession(sessionID)
		if err != nil {
			utils.DeleteCookie(w)
			slog.Error("❌ Failed to create session in system: " + err.Error())
			return 0, err
		}
	} else {
		sessionID = cookie.Value
		code, err := h.serv.IsValidSession(sessionID)
		if code == http.StatusUnauthorized {
			userID, err = h.serv.CreateSession(sessionID)
			if err != nil {
				utils.DeleteCookie(w)
				slog.Error("❌ Failed to create session in system: " + err.Error())
				return userID, nil
			}
		}
		if err != nil {
			utils.DeleteCookie(w)
			slog.Error("❌ Session check error: " + err.Error())
			return 0, err
		}
		user, err := h.serv.GetUser(sessionID)
		if err != nil {
			utils.DeleteCookie(w)
			slog.Error("❌ Get User error: " + err.Error())
			return 0, err
		}
		userID = user.ID
	}
	return userID, nil
}
