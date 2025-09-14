package command

import (
	"encoding/json"
	"fmt"

	"github.com/Nightails/pokedexcli/internal/api"
	"github.com/Nightails/pokedexcli/internal/config"
)

func commandMap(config *config.Config) error {
	var url string
	if config.NextURL == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = config.NextURL
	}

	var data []byte
	// Check if we have the data in the cache
	entry, exist := config.Cache.Get(url)
	if !exist {
		var err error
		data, err = api.GetPokedexAPI(url)
		if err != nil {
			return err
		}
		// Cache the data
		config.Cache.Add(url, data)
	} else {
		data = entry
	}

	names, err := getAreaNames(data, config)
	if err != nil {
		return err
	}
	for _, name := range names {
		fmt.Println(name)
	}

	return nil
}

func commandMapBack(config *config.Config) error {
	if config.PreviousURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	url := config.PreviousURL
	var data []byte
	// Check if we have the data in the cache
	entry, exist := config.Cache.Get(url)
	if !exist {
		var err error
		data, err = api.GetPokedexAPI(url)
		if err != nil {
			return err
		}
		// Cache the data
		config.Cache.Add(url, data)
	} else {
		data = entry
	}

	names, err := getAreaNames(data, config)
	if err != nil {
		return err
	}
	fmt.Println("previous page")
	for _, name := range names {
		fmt.Println(name)
	}

	return nil
}

// getAreaNames returns a slice of area names from the given data
func getAreaNames(data []byte, conf *config.Config) ([]string, error) {
	var apiRes config.Locations
	if err := json.Unmarshal(data, &apiRes); err != nil {
		return []string{}, err
	}

	conf.NextURL = apiRes.Next
	conf.PreviousURL = apiRes.Previous

	var names []string
	for _, str := range apiRes.Results {
		names = append(names, str.Name)
	}
	return names, nil
}
