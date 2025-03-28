package domain

import (
	"io"
	"net/http"
	"time"
)

// CommentService interface (business rules for comments)
type CommentService interface {
	CreateComment(postID, replyID int, content, session_id string, parsedFile []byte) (int, error)
	GetCommentsByPost(postID int) ([]Comment, error)
}

// PostService interface (business rules for posts)
type PostService interface {
	CreatePost(post Post) (int, error)
	ServePost(w http.ResponseWriter, postId string) (int, error)
	ServeArchivePost(w http.ResponseWriter, postIdstr string) (int, error)
	GetArchivePosts() ([]Post, error)
}

// UserSessionService interface (business rules for sessions)
type UserSessionService interface {
	CreateSession(session_id string) (int, error)
	DeleteExpiredSessions() error
	GetSession(id string) (UserSession, error)
	IsValidSession(session_id string) (int, error)
	GetUser(sessionID string) (User, error)
}

type CommentRepoInt interface {
	FindUniqueCommentID() (int, error)
	CreateComment(comment Comment) error
	GetCommentsByPost(postID int) ([]Comment, error)
	IsReplyIdExist(post_id int, reply_id int) (bool, error)
}

type PostRepoInt interface {
	SavePost(post Post) error
	FindUniquePostID() (int, error)
	GetActivePosts() ([]Post, error)
	GetArchivePosts() ([]Post, error)
	GetPost(id int) (Post, error)
	GetArchivePost(id int) (Post, error)
	UpdatePostExpiration(postID int, newExpiration time.Time) error
	IsPostExist(id int) (bool, error)
	ArchiveExpiredPosts() error
}

type SessionRepoInt interface {
	FindUniqueUserID() (int, error)
	SaveUser(user User) error
	ChangeUserName(changed_name string, user_id int) error
	IsNameEqual(changed_name string, user_id int) (bool, error)
	GetUser(session_id string) (User, error)
	GetUserByID(id int) (User, error)
	DeleteExpiredSessions() error
	IsSessionExist(session_id string) (bool, error)
}

type ExternalAPI interface {
	GetCharacter(user *User) error
	GetAvatarCount() (int, error)
}

type S3 interface {
	GetObject(bucket, object string) (io.ReadCloser, error)
	CreateObject(bucket, object string, data io.Reader) (int, error)
	CreateBucket(name string) (int, error)
	InitBuckets() error
}
