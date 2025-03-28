package repository

import (
	"1337bo4rd/internal/domain"
	"database/sql"
)

type SessionRepo struct {
	Db *sql.DB
}

var _ domain.SessionRepoInt = (*SessionRepo)(nil)

func NewSessionRepo(Db *sql.DB) *SessionRepo {
	return &SessionRepo{Db: Db}
}

func (repo *SessionRepo) GetUserByID(id int) (domain.User, error) {
	var user domain.User
	err := repo.Db.QueryRow(`SELECT ID, name, Token_ID, TokenDate, Expires_at, Avatar_URL 
	FROM Users
	WHERE ID=$1;`, id).Scan(&user.ID, &user.Name, &user.Token_ID, &user.TokenDate, &user.Expires_at, &user.ImageURL)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (repo *SessionRepo) FindUniqueUserID() (int, error) {
	var count int
	err := repo.Db.QueryRow(`SELECT COALESCE(MAX(ID), 0) + 1 FROM Users;`).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *SessionRepo) SaveUser(user domain.User) error {
	tx, err := repo.Db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO Users (ID, Name, Token_ID, TokenDate, Expires_at, Avatar_URL) 
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING ID
	`, user.ID, user.Name, user.Token_ID, user.TokenDate, user.Expires_at, user.ImageURL)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (repo *SessionRepo) ChangeUserName(changed_name string, user_id int) error {
	tx, err := repo.Db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE Users
	SET Name=$1
	WHERE ID=$2;`, changed_name, user_id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (repo *SessionRepo) IsNameEqual(changed_name string, user_id int) (bool, error) {
	var exist string
	err := repo.Db.QueryRow(`SELECT Name=$1 FROM Users
	WHERE ID=$2;`, changed_name, user_id).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist == "t", nil
}

func (repo *SessionRepo) GetUser(session_id string) (domain.User, error) {
	var user domain.User
	err := repo.Db.QueryRow(`SELECT ID, name, Token_ID, TokenDate, Expires_at, Avatar_URL 
	FROM Users
	WHERE Token_ID=$1;`, session_id).Scan(&user.ID, &user.Name, &user.Token_ID, &user.TokenDate, &user.Expires_at, &user.ImageURL)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (repo *SessionRepo) DeleteExpiredSessions() error {
	_, err := repo.Db.Exec(`DELETE FROM Users
	WHERE CURRENT_TIMESTAMP>Expires_at;`)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SessionRepo) IsSessionExist(session_id string) (bool, error) {
	var count int
	err := repo.Db.QueryRow(`SELECT COUNT(*) FROM Users
	WHERE Token_ID=$1;`, session_id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, err
}
