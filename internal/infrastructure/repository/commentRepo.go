package repository

import (
	"1337bo4rd/internal/domain"
	"database/sql"
	"time"
)

type CommentRepository struct {
	DB *sql.DB
}

var _ domain.CommentRepoInt = (*CommentRepository)(nil)

func NewCommentRepository(DB *sql.DB) *CommentRepository {
	return &CommentRepository{DB: DB}
}

func (repo *CommentRepository) CreateComment(comment domain.Comment) error {
	_, err := repo.DB.Exec(`
		INSERT INTO comments (post_id, Reply_to, content, Author_id, ImageURL) 
		VALUES ($1, $2, $3, $4, $5)`,
		comment.PostID, comment.ReplyToID, comment.Content, comment.Author_id, comment.ImageLink,
	)
	return err
}

func (repo *CommentRepository) GetCommentsByPost(postID int) ([]domain.Comment, error) {
	rows, err := repo.DB.Query(`
		SELECT c.ID, c.Post_id, COALESCE(c.Reply_to, 0)::INTEGER, c.Content, u.Avatar_URL, u.Name, c.Created_at, COALESCE(c.ImageURL,'')  FROM Comments c
		INNER JOIN Users u On c.Author_id=u.ID
		WHERE c.Post_id = $1 ORDER BY c.Created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.ReplyToID, &comment.Content, &comment.AvatarURL, &comment.Username, &comment.CreatedAt, &comment.ImageLink)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (repo *PostRepository) UpdatePostExpiration(postID int, newExpiration time.Time) error {
	_, err := repo.Db.Exec("UPDATE posts SET expires_at = $1 WHERE id = $2", newExpiration, postID)
	return err
}

func (repo *CommentRepository) FindUniqueCommentID() (int, error) {
	var id int
	err := repo.DB.QueryRow("SELECT COALESCE(MAX(ID), 0) + 1 FROM comments").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *CommentRepository) IsReplyIdExist(post_id int, reply_id int) (bool, error) {
	var count int
	err := repo.DB.QueryRow(`SELECT COUNT(c.ID) FROM comments c 
	INNER JOIN posts p 
	ON p.ID=c.Post_id
	WHERE c.ID=$1 AND p.ID=$2
	`, reply_id, post_id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
