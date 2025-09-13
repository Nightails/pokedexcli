package api

import (
	"errors"
	"io"
	"net/http"
)

func GetPokedexAPI(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = errors.New("failed to close body")
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
