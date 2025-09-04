package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type Config struct {
	NextURL     string
	PreviousURL string
}

type PokedexAPIResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetAreaNames(url string, config *Config) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []string{}, err
	}

	var apiRes PokedexAPIResponse
	if err := json.Unmarshal(body, &apiRes); err != nil {
		return []string{}, err
	}

	config.NextURL = apiRes.Next
	config.PreviousURL = apiRes.Previous

	var names []string
	for _, str := range apiRes.Results {
		names = append(names, str.Name)
	}
	return names, nil
}
