package service_test

import (
	"1337bo4rd/internal/domain"
	"1337bo4rd/internal/infrastructure/repository"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = setupTestDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to test database: %v", err)
	}
	defer testDB.Close()

	exitCode := m.Run()
	os.Exit(exitCode)
}

// Подключаемся к тестовой БД
func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=0.0.0.0 user=hacker password=password dbname=1337bo4rd port=5432 sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Clear table before the test
func clearTable(db *sql.DB) {
	db.Exec("DELETE FROM posts WHERE title LIKE 'TEST_%'")
}

func TestPostRepository_SavePost(t *testing.T) {
	repo := repository.NewPostRepository(testDB)
	clearTable(testDB)
	var err error
	post := domain.Post{
		AuthorID:  1,
		Title:     "TEST_Title",
		Content:   "Test Content",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	err = repo.SavePost(post)
	if err != nil {
		t.Fatalf("❌ SavePost() error: %v", err)
	}
	post.ID, err = repo.FindUniquePostID()
	if err != nil {
		t.Fatalf("❌ FindUniquePostID error: %v", err)
	}
	post.ID--
	savedPost, err := repo.GetPost(post.ID)
	if err != nil {
		t.Fatalf("❌ GetPost() error: %v", err)
	}

	if savedPost.Title != post.Title {
		t.Errorf("❌ Expected title %s, got %s", post.Title, savedPost.Title)
	}
}

func TestPostRepository_FindUniquePostID(t *testing.T) {
	repo := repository.NewPostRepository(testDB)
	clearTable(testDB)
	clearTable(testDB)      // ✅ Clear table before the test
	resetIDSequence(testDB) // ✅ Reset sequence to start fresh

	id, err := repo.FindUniquePostID()
	if err != nil {
		t.Fatalf("❌ FindUniquePostID() error: %v", err)
	}

	if id != 1 {
		t.Errorf("❌ Expected post ID 1, got %d", id)
	}
}

func resetIDSequence(db *sql.DB) {
	db.Exec("ALTER SEQUENCE posts_id_seq RESTART WITH 1") // Adjust sequence name if needed
}
