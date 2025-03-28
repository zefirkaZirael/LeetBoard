package handler

import (
	"1337bo4rd/internal/UI/middleware"
	"1337bo4rd/internal/domain"
	"1337bo4rd/utils"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

type PostHandler struct {
	serv        domain.PostService
	sessionServ domain.UserSessionService
	DB          *sql.DB
}

func NewPostHandler(serv domain.PostService, sessionServ domain.UserSessionService, db *sql.DB) *PostHandler {
	return &PostHandler{serv: serv, DB: db, sessionServ: sessionServ}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	sessionhandle := middleware.NewSessionHandler(h.sessionServ, h.DB)

	var post domain.Post
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
	post.Changed_name = r.FormValue("Name")
	post.Title = r.FormValue("Title")
	post.Content = r.FormValue("Content")

	// File Handling
	file, _, err := r.FormFile("File")
	if err != nil {
		if err != http.ErrMissingFile {
			slog.Error("❌ Form file reading error: " + err.Error())
			utils.ErrorPage(w, errors.New("form file reading error: "+err.Error()), http.StatusBadRequest)
			return
		}
	} else {
		defer file.Close()
		post.File, err = io.ReadAll(file)
		if err != nil {
			slog.Error("❌ Form file reading error: " + err.Error())
			utils.ErrorPage(w, errors.New("file reading error: "+err.Error()), http.StatusBadRequest)
			return
		}
	}

	userId, err := sessionhandle.Authorization(w, r)
	if err != nil {
		slog.Error("❌ Authorization error: " + err.Error())
		utils.ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	// Post Validation
	err = h.PostValidation(post)
	if err != nil {
		slog.Error("❌ Post Validation error: " + err.Error())
		utils.ErrorPage(w, err, http.StatusBadRequest)
		return
	}
	post.AuthorID = userId
	// Calling business logic
	code, err := h.serv.CreatePost(post)
	if err != nil {
		slog.Error("❌ Create Post error: " + err.Error())
		utils.ErrorPage(w, err, code)
		return
	}
	sessionhandle.MainPage(w, r)
}

func (h *PostHandler) ServePost(w http.ResponseWriter, r *http.Request) {
	err := h.sessionServ.DeleteExpiredSessions()
	if err != nil {
		slog.Error("❌ Failed to delete expires sessions in system: " + err.Error())
		utils.ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	postId := r.PathValue("id")
	code, err := h.serv.ServePost(w, postId)
	if err != nil {
		slog.Error("❌ Serve Post error: " + err.Error())
		utils.ErrorPage(w, err, code)
		return
	}
}

func (h *PostHandler) ServeArchivePost(w http.ResponseWriter, r *http.Request) {
	err := h.sessionServ.DeleteExpiredSessions()
	if err != nil {
		slog.Error("❌ Failed to delete expires sessions in system: " + err.Error())
		utils.ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	postId := r.PathValue("id")
	code, err := h.serv.ServeArchivePost(w, postId)
	if err != nil {
		slog.Error("❌ Serve Archive Post error: " + err.Error())
		utils.ErrorPage(w, err, code)
		return
	}
}

func (h *PostHandler) PostValidation(post domain.Post) error {
	if post.ID != 0 {
		return fmt.Errorf("post ID must be empty")
	}
	if post.Title == "" {
		return fmt.Errorf("post title is empty")
	}
	if post.Content == "" {
		return fmt.Errorf("post content is empty")
	}
	if post.AuthorID != 0 {
		return fmt.Errorf("author_id field must be empty")
	}
	if !post.CreatedAt.IsZero() {
		return fmt.Errorf("createdAt field must be empty")
	}
	if !post.ExpiresAt.IsZero() {
		return fmt.Errorf("expires_at field must be empty")
	}
	if post.ImageLink != "" {
		return fmt.Errorf("image URL field must be empty")
	}
	return nil
}

func (h *PostHandler) ServeArchivePage(w http.ResponseWriter, r *http.Request) {
	err := h.sessionServ.DeleteExpiredSessions()
	if err != nil {
		slog.Error("❌ Failed to delete expires sessions in system: " + err.Error())
		utils.ErrorPage(w, err, http.StatusInternalServerError)
		return
	}

	// Get all archived posts from the repository
	posts, err := h.serv.GetArchivePosts()
	if err != nil {
		slog.Error("❌ Failed to fetch archive posts: " + err.Error())
		utils.ErrorPage(w, err, http.StatusInternalServerError)
		return
	}

	// Load the HTML template
	temp, err := template.ParseFiles("web/templates/archive.html")
	if err != nil {
		slog.Error("❌ Failed to parse archive template: " + err.Error())
		utils.ErrorPage(w, err, http.StatusInternalServerError)
		return
	}

	// Pass the posts to the template
	err = temp.Execute(w, posts)
	if err != nil {
		slog.Error("❌ Failed to render archive page: " + err.Error())
		utils.ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
}
