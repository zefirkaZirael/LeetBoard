package service

import (
	"1337bo4rd/internal/domain"
	"1337bo4rd/internal/infrastructure/s3"
	"1337bo4rd/utils"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type CommentService struct {
	Repo        domain.CommentRepoInt
	PostRepo    domain.PostRepoInt
	StorageRepo *s3.S3
	SessionRepo domain.SessionRepoInt
}

// Constructor
func NewCommentService(repo domain.CommentRepoInt, postRepo domain.PostRepoInt, sessionRepo domain.SessionRepoInt, StorageRepo *s3.S3) *CommentService {
	return &CommentService{
		Repo:        repo,
		PostRepo:    postRepo,
		StorageRepo: StorageRepo,
		SessionRepo: sessionRepo,
	}
}

var _ domain.CommentService = (*CommentService)(nil)

// Create a new comment
func (s *CommentService) CreateComment(postID, replyID int, content, session_id string, parsedFile []byte) (int, error) {
	_, err := s.PostRepo.GetPost(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("❌ Error: post is not found in database")
			return http.StatusBadRequest, errors.New("post is not exist")
		}
		slog.Error("❌ Error on GetPost: " + err.Error())
		return http.StatusInternalServerError, err
	}
	if content == "" {
		slog.Error("❌ Error: comment content is empty")
		return http.StatusBadRequest, errors.New("content cannot be empty")
	}
	if replyID != 0 {
		exist, err := s.Repo.IsReplyIdExist(postID, replyID)
		if err != nil {
			slog.Error("❌ Error on IsReplyIdExist:" + err.Error())
			return http.StatusInternalServerError, err
		}
		if !exist {
			slog.Error("❌ Error: Comment reply id is not exist")
			return http.StatusBadRequest, errors.New("comment reply id is not exist")
		}
	}
	user, err := s.SessionRepo.GetUser(session_id)
	if err != nil {
		slog.Error("❌ Error on GetUser", "session_id", session_id, "error", err)
		return http.StatusInternalServerError, err
	}

	// Generate id before creating comment
	commentID, err := s.Repo.FindUniqueCommentID()
	if err != nil {
		slog.Error("❌ Error on FindUniqueCommentID", "error", err)
		return http.StatusInternalServerError, err
	}

	comment := domain.Comment{
		ID:        commentID,
		PostID:    postID,
		Author_id: user.ID,
		Content:   content,
		AvatarURL: user.ImageURL,
		Username:  user.Name,
		ReplyToID: replyID,
		CreatedAt: time.Now(),
	}

	if parsedFile != nil {
		extension, err := utils.DetectType(parsedFile)
		if err != nil {
			return http.StatusBadRequest, err
		}
		body := bytes.NewReader(parsedFile)
		objectname := strconv.Itoa(postID) + "_" + strconv.Itoa(commentID) + extension
		code, err := s.StorageRepo.CreateObject(s.StorageRepo.DefCommentDir, objectname, body)
		if err != nil {
			slog.Error("❌ Create object error: " + err.Error())
			return code, err
		}
		comment.ImageLink = s.StorageRepo.S3url + "/" + s.StorageRepo.DefCommentDir + "/" + objectname
	}

	// Store comment
	err = s.Repo.CreateComment(comment)
	if err != nil {
		fmt.Println("❌ Error on CreateComment:", err)
		slog.Error("❌ Error on CreateComment", "comment", comment, "error", err)
		return http.StatusInternalServerError, err
	}

	// Update post expiration (+15 min)
	err = s.PostRepo.UpdatePostExpiration(postID, time.Now().Add(15*time.Minute))
	if err != nil {
		fmt.Println("❌ Error on UpdatePostExpiration:", err)
		slog.Error("❌ Error on UpdatePostExpiration", "postID", postID, "error", err)
	}

	return comment.ID, nil
}

// Business rule: Get comments for a post
func (s *CommentService) GetCommentsByPost(postID int) ([]domain.Comment, error) {
	repoComments, err := s.Repo.GetCommentsByPost(postID)
	if err != nil {
		return nil, err
	}

	var comments []domain.Comment
	for _, repoComment := range repoComments {
		comments = append(comments, domain.Comment{
			ID:        repoComment.ID,
			PostID:    repoComment.PostID,
			ReplyToID: repoComment.ReplyToID,
			Content:   repoComment.Content,
			AvatarURL: repoComment.AvatarURL,
			Username:  repoComment.Username,
			CreatedAt: repoComment.CreatedAt,
		})
	}
	return comments, nil
}
