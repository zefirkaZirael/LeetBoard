package external

import (
	"1337bo4rd/internal/domain"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

var _ domain.ExternalAPI = (*ExternalAdapter)(nil)

var rickMortiURL string = "https://rickandmortyapi.com/api/character/"

type RickAndMortyCharacter struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type ExternalAdapter struct {
	URL string
}

func DefaultExternalAPI() *ExternalAdapter {
	return &ExternalAdapter{URL: rickMortiURL}
}

func (e *ExternalAdapter) GetCharacter(user *domain.User) error {
	resp, err := http.Get(e.URL + strconv.Itoa(user.ID))
	if err != nil {
		return err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Temporary variable for parsing
	var character RickAndMortyCharacter
	err = json.Unmarshal(data, &character)
	if err != nil {
		return err
	}

	// Writing data in user
	user.Name = character.Name
	user.ImageURL = character.Image

	return nil
}

func (e *ExternalAdapter) GetAvatarCount() (int, error) {
	resp, err := http.Get(e.URL)
	if err != nil {
		return 0, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var data struct {
		Info struct {
			Count int `json:"count"`
		} `json:"info"`
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return 0, err
	}
	return data.Info.Count, nil
}
