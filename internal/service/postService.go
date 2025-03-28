package service

import (
	"1337bo4rd/internal/domain"
	"1337bo4rd/internal/infrastructure/repository"
	"1337bo4rd/utils"
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// PostService struct
type PostService struct {
	Repo        domain.PostRepoInt
	SessionRepo domain.SessionRepoInt
	StorageRepo domain.S3
	DB          *sql.DB
}

// Constructor
func NewPostService(repo domain.PostRepoInt, sessionRepo domain.SessionRepoInt, storageRepo domain.S3, db *sql.DB) *PostService {
	return &PostService{
		Repo:        repo,
		SessionRepo: sessionRepo,
		StorageRepo: storageRepo,
		DB:          db,
	}
}

var _ domain.PostService = (*PostService)(nil)

// Create a new post
func (s *PostService) CreatePost(post domain.Post) (int, error) {
	var err error
	var equal bool
	// Give Unique Post ID
	post.ID, err = s.Repo.FindUniquePostID()
	if err != nil {
		slog.Error("❌ Find Unique Id error: " + err.Error())
		return http.StatusInternalServerError, err
	}
	if post.Changed_name != "" {
		// Use injected session repository to check if the user name is equal
		equal, err = s.SessionRepo.IsNameEqual(post.Changed_name, post.AuthorID)
		if err != nil {
			slog.Error("❌ Name equal check error: " + err.Error())
			return http.StatusInternalServerError, err
		}
		if !equal {
			// Change the user name if not equal
			err := s.SessionRepo.ChangeUserName(post.Changed_name, post.AuthorID)
			if err != nil {
				slog.Error("❌ Change User Name error: " + err.Error())
				return http.StatusInternalServerError, err
			}
		}
	}
	// Set post time
	post.CreatedAt = time.Now()
	post.ExpiresAt = post.CreatedAt.Add(10 * time.Minute)

	// Handle file upload (image processing)
	if post.File != nil {
		extension, err := utils.DetectType(post.File)
		if err != nil {
			return http.StatusBadRequest, err
		}
		body := bytes.NewReader(post.File)
		objectname := strconv.Itoa(post.ID) + extension
		code, err := s.StorageRepo.CreateObject(domain.DefPostDir, objectname, body)
		if err != nil {
			slog.Error("❌ Create object error: " + err.Error())
			return code, err
		}
		post.ImageLink = domain.S3url + "/" + domain.DefPostDir + "/" + objectname
	}

	// Save post to database
	err = s.Repo.SavePost(post)
	if err != nil {
		slog.Error("❌ Save post error: " + err.Error())
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *PostService) ServeArchivePost(w http.ResponseWriter, postIdstr string) (int, error) {
	postId, err := strconv.Atoi(postIdstr)
	if err != nil {
		slog.Error("❌ Failed to convert post id: " + "invalid syntax: " + postIdstr)
		return http.StatusBadRequest, errors.New("invalid syntax: " + postIdstr)
	}

	// Checks is Post exist by ID
	exist, err := s.Repo.IsPostExist(postId)
	if err != nil {
		slog.Error("❌ Failed to Check is post exist: " + err.Error())
		return http.StatusInternalServerError, err
	}
	if !exist {
		NotexistErr := errors.New("post is not exist")
		slog.Error("❌ Failed to Serve Post: " + NotexistErr.Error())
		return http.StatusBadRequest, NotexistErr
	}

	// Get the post from the repository
	post, err := s.Repo.GetArchivePost(postId)
	if err == sql.ErrNoRows {
		slog.Error("❌ Archive Post is not exist")
		return http.StatusNotFound, errors.New("archive post is not found")
	} else if err != nil {
		slog.Error("❌ Failed to Get Archive Post information: " + err.Error())
		return http.StatusInternalServerError, err
	}

	// Use injected session repository to get user data
	author, err := s.SessionRepo.GetUserByID(post.AuthorID)
	if err != nil {
		slog.Error("❌ Failed to get user: " + err.Error())
		return http.StatusInternalServerError, err
	}

	// Parse the post HTML template
	temp, err := template.ParseFiles("web/templates/archive-post.html")
	if err != nil {
		slog.Error("❌ Failed to parse template file: " + err.Error())
		return http.StatusInternalServerError, err
	}

	// Get comments related to the post
	comments, err := repository.NewCommentRepository(s.DB).GetCommentsByPost(postId)
	if err != nil {
		slog.Error("❌ Failed to get comments by post: " + err.Error())
		return http.StatusInternalServerError, err
	}

	// Prepare the data to be passed into the template
	data := struct {
		ImageURL  string
		Name      string
		CreatedAt string
		ID        int
		ImageLink string
		Title     string
		Content   string
		Comments  []domain.Comment
	}{
		ImageURL:  author.ImageURL,
		Name:      author.Name,
		CreatedAt: post.CreatedAt.Format("02 January 2006, 15:04:05 UTC"),
		ID:        postId,
		ImageLink: post.ImageLink,
		Title:     post.Title,
		Content:   post.Content,
		Comments:  comments,
	}

	// Execute the template with the data
	err = temp.Execute(w, data)
	if err != nil {
		slog.Error("❌ Failed to execute template: " + err.Error())
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// Serve a post
func (s *PostService) ServePost(w http.ResponseWriter, postIdstr string) (int, error) {
	// Convert postIdstr to an integer
	postId, err := strconv.Atoi(postIdstr)
	if err != nil {
		slog.Error("❌ Failed to convert post id: " + "invalid syntax: " + postIdstr)
		return http.StatusBadRequest, errors.New("invalid syntax: " + postIdstr)
	}

	// Checks is Post exist by ID
	exist, err := s.Repo.IsPostExist(postId)
	if err != nil {
		slog.Error("❌ Failed to Check is post exist: " + err.Error())
		return http.StatusInternalServerError, err
	}
	if !exist {
		NotexistErr := errors.New("post is not exist")
		slog.Error("❌ Failed to Serve Post: " + NotexistErr.Error())
		return http.StatusBadRequest, NotexistErr
	}

	// Get the post from the repository
	post, err := s.Repo.GetPost(postId)
	if err == sql.ErrNoRows {
		slog.Error("❌ Post is not exist")
		return http.StatusNotFound, errors.New("post is not found")
	} else if err != nil {
		slog.Error("❌ Failed to Get Post information: " + err.Error())
		return http.StatusInternalServerError, err
	}

	// Use injected session repository to get user data
	author, err := s.SessionRepo.GetUserByID(post.AuthorID)
	if err != nil {
		slog.Error("❌ Failed to get user: " + err.Error())
		return http.StatusInternalServerError, err
	}

	// Parse the post HTML template
	temp, err := template.ParseFiles("web/templates/post.html")
	if err != nil {
		slog.Error("❌ Failed to parse template file: " + err.Error())
		return http.StatusInternalServerError, err
	}

	// Get comments related to the post
	comments, err := repository.NewCommentRepository(s.DB).GetCommentsByPost(postId)
	if err != nil {
		slog.Error("❌ Failed to get comments by post: " + err.Error())
		return http.StatusInternalServerError, err
	}

	// Prepare the data to be passed into the template
	data := struct {
		ImageURL  string
		Name      string
		CreatedAt string
		ID        int
		ImageLink string
		Title     string
		Content   string
		Comments  []domain.Comment
	}{
		ImageURL:  author.ImageURL,
		Name:      author.Name,
		CreatedAt: post.CreatedAt.Format("02 January 2006, 15:04:05 UTC"),
		ID:        postId,
		ImageLink: post.ImageLink,
		Title:     post.Title,
		Content:   post.Content,
		Comments:  comments,
	}

	// Execute the template with the data
	err = temp.Execute(w, data)
	if err != nil {
		slog.Error("❌ Failed to execute template: " + err.Error())
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *PostService) GetArchivePosts() ([]domain.Post, error) {
	return s.Repo.GetArchivePosts()
}
