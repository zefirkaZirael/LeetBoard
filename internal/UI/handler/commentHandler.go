package handler

import (
	"1337bo4rd/internal/UI/middleware"
	"1337bo4rd/internal/domain"
	"1337bo4rd/utils"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	serv        domain.CommentService
	sessionServ domain.UserSessionService
	DB          *sql.DB
}

func NewCommentHandler(serv domain.CommentService, sessionServ domain.UserSessionService, db *sql.DB) *CommentHandler {
	return &CommentHandler{serv: serv, DB: db, sessionServ: sessionServ}
}

func (h *CommentHandler) SubmitCommentHandler(w http.ResponseWriter, r *http.Request) {
	sessionHandler := middleware.NewSessionHandler(h.sessionServ, h.DB)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		var code int = http.StatusInternalServerError
		slog.Error("❌ Failed to parse form: " + err.Error())
		if http.ErrNotMultipart == err {
			code = http.StatusBadRequest
		}
		utils.ErrorPage(w, err, code)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("postID"))
	if err != nil {
		slog.Error("❌ Invalid post ID: " + err.Error())
		utils.ErrorPage(w, errors.New("post id invalid syntax"), http.StatusBadRequest)
		return
	}

	var replyToID int
	if r.FormValue("ReplyTo") != "" {
		replyToID, err = strconv.Atoi(r.FormValue("ReplyTo"))
		if err != nil {
			slog.Error("❌ Invalid reply ID: " + err.Error())
			utils.ErrorPage(w, errors.New("reply id invalid syntax"), http.StatusBadRequest)
			return
		}
	}
	var parsedFile []byte

	file, _, err := r.FormFile("File")
	if err != nil {
		if err != http.ErrMissingFile {
			slog.Error("❌ Form file reading error: " + err.Error())
			utils.ErrorPage(w, errors.New("form file reading error: "+err.Error()), http.StatusBadRequest)
			return
		}
	} else {
		defer file.Close()
		parsedFile, err = io.ReadAll(file)
		if err != nil {
			slog.Error("❌ Form file reading error: " + err.Error())
			utils.ErrorPage(w, errors.New("Form file reading error: "+err.Error()), http.StatusBadRequest)
			return
		}
	}

	content := r.FormValue("Content")
	if content == "" {
		slog.Error("❌ Form Value check error: Content cannot be empty")
		utils.ErrorPage(w, errors.New("content cannot be empty"), http.StatusBadRequest)
		return
	}

	_, err = sessionHandler.Authorization(w, r)
	if err != nil {
		slog.Error("Authorization error: " + err.Error())
		utils.ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		var code int = http.StatusInternalServerError
		if err == http.ErrNoCookie {
			code = http.StatusUnauthorized
		}
		slog.Error("Get session id error: " + err.Error())
		utils.ErrorPage(w, err, code)
		return
	}
	session_id := cookie.Value
	code, err := h.serv.CreateComment(postID, replyToID, content, session_id, parsedFile)
	if err != nil {
		slog.Error("Failed to submit comment: " + err.Error())
		utils.ErrorPage(w, err, code)
		return
	}

	sessionHandler.MainPage(w, r)
}
