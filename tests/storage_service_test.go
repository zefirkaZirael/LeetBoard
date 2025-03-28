package service_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

// 🔹 Мок-хранилище для S3
type MockStorageRepo struct {
	objects map[string][]byte
	buckets map[string]bool
}

func NewMockStorageRepo() *MockStorageRepo {
	return &MockStorageRepo{
		objects: make(map[string][]byte),
		buckets: make(map[string]bool),
	}
}

func (m *MockStorageRepo) CreateObject(bucket, object string, data io.Reader) (int, error) {
	if !m.buckets[bucket] {
		return http.StatusNotFound, errors.New("bucket not found")
	}
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, data)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	m.objects[bucket+"/"+object] = buf.Bytes()
	return http.StatusOK, nil
}

func (m *MockStorageRepo) GetObject(bucket, object string) (io.ReadCloser, error) {
	data, exists := m.objects[bucket+"/"+object]
	if !exists {
		return nil, errors.New("object not found")
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (m *MockStorageRepo) CreateBucket(name string) (int, error) {
	if _, exists := m.buckets[name]; exists {
		return http.StatusConflict, errors.New("bucket already exists")
	}
	m.buckets[name] = true
	return http.StatusOK, nil
}

// 🔹 Test: Creating an object
func TestCreateObject(t *testing.T) {
	mockStorage := NewMockStorageRepo()
	mockStorage.CreateBucket("test-bucket")

	data := bytes.NewBufferString("test data")
	status, err := mockStorage.CreateObject("test-bucket", "test-object", data)
	if err != nil {
		t.Errorf("❌ CreateObject() error: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("❌ Expected status 200, got %d", status)
	}
}

// 🔹 Test: Get an object
func TestGetObject(t *testing.T) {
	mockStorage := NewMockStorageRepo()
	mockStorage.CreateBucket("test-bucket")
	mockStorage.CreateObject("test-bucket", "test-object", bytes.NewBufferString("test data"))

	obj, err := mockStorage.GetObject("test-bucket", "test-object")
	if err != nil {
		t.Errorf("❌ GetObject() error: %v", err)
	}

	data, _ := io.ReadAll(obj)
	if string(data) != "test data" {
		t.Errorf("❌ Expected 'test data', got '%s'", string(data))
	}
}

// 🔹 Test: Creating a bucket
func TestCreateBucket(t *testing.T) {
	mockStorage := NewMockStorageRepo()
	status, err := mockStorage.CreateBucket("new-bucket")
	if err != nil {
		t.Errorf("❌ CreateBucket() error: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("❌ Expected status 200, got %d", status)
	}
}
