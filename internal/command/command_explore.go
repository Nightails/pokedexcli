package command

import (
	"encoding/json"
	"fmt"

	"github.com/Nightails/pokedexcli/internal/api"
	"github.com/Nightails/pokedexcli/internal/config"
)

func commandExplore(conf *config.Config) error {
	if conf.Argument == "" {
		fmt.Println("please specify an area")
		return nil
	}

	url := "https://pokeapi.co/api/v2/location-area/" + conf.Argument
	var data []byte
	// Check if we have the data in the cache
	entry, exist := conf.Cache.Get(url)
	if !exist {
		var err error
		data, err = api.GetPokedexAPI(url)
		if err != nil {
			return err
		}
		// Cache the data
		conf.Cache.Add(url, data)
	} else {
		data = entry
	}

	names, err := getEncounterPokemons(data, conf)
	if err != nil {
		return err
	}
	fmt.Println("Exploring " + conf.Argument + "...")
	fmt.Println("Found Pokemon:")
	for _, name := range names {
		fmt.Println(name)
	}

	return nil
}

// getEncounterPokemons returns a slice of pokemon names from the given data
func getEncounterPokemons(data []byte, conf *config.Config) ([]string, error) {
	var apiRes config.LocationArea
	if err := json.Unmarshal(data, &apiRes); err != nil {
		return []string{}, err
	}

	var names []string
	for _, encounter := range apiRes.PokemonEncounters {
		names = append(names, encounter.Pokemon.Name)
	}

	return names, nil
}
