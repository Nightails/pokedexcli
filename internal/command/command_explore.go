package command

import (
	"encoding/json"
	"fmt"

	"github.com/Nightails/pokedexcli/internal/api"
	"github.com/Nightails/pokedexcli/internal/config"
)

// LocationAreaResponse is the response from the PokeAPI, with details about a location area
type LocationAreaResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func commandExplore(config *config.Config) error {
	if config.Argument == "" {
		fmt.Println("please specify an area")
		return nil
	}

	url := "https://pokeapi.co/api/v2/location-area/" + config.Argument
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

	names, err := getEncounterPokemons(data, config)
	if err != nil {
		return err
	}
	fmt.Println("Exploring " + config.Argument + "...")
	fmt.Println("Found Pokemon:")
	for _, name := range names {
		fmt.Println(name)
	}

	return nil
}

// getEncounterPokemons returns a slice of pokemon names from the given data
func getEncounterPokemons(data []byte, config *config.Config) ([]string, error) {
	var apiRes LocationAreaResponse
	if err := json.Unmarshal(data, &apiRes); err != nil {
		return []string{}, err
	}

	var names []string
	for _, encounter := range apiRes.PokemonEncounters {
		names = append(names, encounter.Pokemon.Name)
	}

	return names, nil
}
