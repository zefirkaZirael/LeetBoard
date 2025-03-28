package repository

import (
	"1337bo4rd/internal/domain"
	"database/sql"
)

type PostRepository struct {
	Db *sql.DB
}

func NewPostRepository(Db *sql.DB) *PostRepository {
	return &PostRepository{Db: Db}
}

var _ domain.PostRepoInt = (*PostRepository)(nil)

func (repo *PostRepository) IsPostExist(id int) (bool, error) {
	var count int
	err := repo.Db.QueryRow(`SELECT COUNT(*) FROM posts
	WHERE ID=$1;`, id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, err
}

func (repo *PostRepository) FindUniquePostID() (int, error) {
	var count int
	err := repo.Db.QueryRow(`SELECT COALESCE(MAX(ID), 0) + 1 FROM Posts;`).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *PostRepository) SavePost(post domain.Post) error {
	tx, err := repo.Db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO posts(Author_id, Title, Content, Created_at, Expires_at, ImageURL) VALUES
	($1, $2, $3, $4, $5, $6)
	`, post.AuthorID, post.Title, post.Content, post.CreatedAt, post.ExpiresAt, post.ImageLink)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (repo *PostRepository) GetActivePosts() ([]domain.Post, error) {
	rows, err := repo.Db.Query(`SELECT ID, Title, Content, Author_id, Created_at, Expires_at, coalesce(ImageURL,'empty') 
	FROM Posts
	WHERE CURRENT_TIMESTAMP<Expires_at;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.ExpiresAt, &post.ImageLink)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// Get all posts
func (repo *PostRepository) GetArchivePosts() ([]domain.Post, error) {
	rows, err := repo.Db.Query(`SELECT ID, Title, Content, Author_id, Created_at, Expires_at, coalesce(ImageURL,'empty') 
	FROM Posts
	WHERE CURRENT_TIMESTAMP>=Expires_at;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.ExpiresAt, &post.ImageLink)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *PostRepository) GetPost(id int) (domain.Post, error) {
	var post domain.Post
	err := repo.Db.QueryRow(`SELECT ID, Title, Content, Author_id, Created_at, Expires_at, coalesce(ImageURL,'empty') 
	FROM Posts
	WHERE ID=$1 AND CURRENT_TIMESTAMP<Expires_at;
	`, id).Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.ExpiresAt, &post.ImageLink)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (repo *PostRepository) GetArchivePost(id int) (domain.Post, error) {
	var post domain.Post
	err := repo.Db.QueryRow(`SELECT ID, Title, Content, Author_id, Created_at, Expires_at, coalesce(ImageURL,'empty') 
	FROM Posts
	WHERE ID=$1 AND CURRENT_TIMESTAMP>=Expires_at;
	`, id).Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.ExpiresAt, &post.ImageLink)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (repo *PostRepository) ArchiveExpiredPosts() error {
	_, err := repo.Db.Exec(`UPDATE Posts SET archived = TRUE WHERE CURRENT_TIMESTAMP >= Expires_at;`)
	return err
}
