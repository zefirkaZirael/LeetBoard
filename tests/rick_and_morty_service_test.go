package service_test

import (
	"1337bo4rd/internal/infrastructure/external"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// 🔹 Мок API Rick and Morty
type MockRickAndMortyAPI struct {
	Response string
	Fail     bool
}

func (m *MockRickAndMortyAPI) GetCharacter(user *external.RickAndMortyCharacter) error {
	if m.Fail {
		return errors.New("API error")
	}

	// JSON response simulation
	response := `{"id": 1, "name": "Rick Sanchez", "image": "http://example.com/rick.jpg"}`
	if m.Response != "" {
		response = m.Response
	}

	err := json.Unmarshal([]byte(response), user)
	if err != nil {
		return err
	}

	return nil
}

// 🔹 Test: succesfully get character case
func TestGetCharacter_Success(t *testing.T) {
	mockAPI := &MockRickAndMortyAPI{}

	var character external.RickAndMortyCharacter
	err := mockAPI.GetCharacter(&character)
	if err != nil {
		t.Errorf("❌ GetCharacter() error: %v", err)
	}
	if character.Name != "Rick Sanchez" {
		t.Errorf("❌ Expected 'Rick Sanchez', got '%s'", character.Name)
	}
	if !strings.HasPrefix(character.Image, "http") {
		t.Errorf("❌Invalid image URL: %s", character.Image)
	}
}

// 🔹 Test: API error response handling
func TestGetCharacter_APIError(t *testing.T) {
	mockAPI := &MockRickAndMortyAPI{Fail: true}

	var character external.RickAndMortyCharacter
	err := mockAPI.GetCharacter(&character)

	if err == nil {
		t.Errorf("❌ Expected API error, got nil")
	}
}
