package s3

import (
	"1337bo4rd/internal/domain"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type S3 struct {
	S3url         string
	DefCommentDir string // Pictures from Comments stored here
	DefPostDir    string // Pictures from Posts stored here
}

func NewS3Repo() *S3 {
	return &S3{S3url: domain.S3url, DefCommentDir: domain.DefCommentDir, DefPostDir: domain.DefPostDir}
}

var _ domain.S3 = (*S3)(nil)

// Returns binary file of object
func (repo *S3) GetObject(bucket, object string) (io.ReadCloser, error) {
	url := repo.S3url + "/" + bucket + "/" + object
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("response from server: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return resp.Body, nil
}

// Creates new object in storage
func (repo *S3) CreateObject(bucket, object string, data io.Reader) (int, error) {
	url := repo.S3url + "/" + bucket + "/" + object
	req, err := http.NewRequest("PUT", url, data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	fmt.Println("Response from server: ", string(body))

	return resp.StatusCode, nil
}

// Creates bucket in s3 storage
func (repo *S3) CreateBucket(name string) (int, error) {
	url := repo.S3url + "/" + name
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return 0, err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	slog.Info("Response from server: " + string(body))
	return resp.StatusCode, nil
}

func (repo *S3) InitBuckets() error {
	_, err := repo.CreateBucket(repo.DefPostDir)
	if err != nil {
		return err
	}
	_, err = repo.CreateBucket(repo.DefCommentDir)
	if err != nil {
		return err
	}
	return nil
}
