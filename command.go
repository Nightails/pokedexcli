package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Nightails/pokedexcli/internal/api"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

// LocationsResponse is the response from the PokeAPI, with a list of location areas
type LocationsResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

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

func getCommand(cmd string) (cliCommand, error) {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Show how to use Pokedexs",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Show 20 location names",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map back",
			description: "Show the previous 20 location names",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore the area",
			callback:    commandExplore,
		},
	}

	if c, exist := commands[cmd]; exist {
		return c, nil
	}
	return cliCommand{}, errors.New("unknown command")
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	message := `
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex

map: List 20 location names
mapb: List previous 20 location names
`
	fmt.Printf("%v/n", message)
	return nil
}

func commandMap(config *Config) error {
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

func commandMapBack(config *Config) error {
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

func commandExplore(config *Config) error {
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

// getAreaNames returns a slice of area names from the given data
func getAreaNames(data []byte, config *Config) ([]string, error) {
	var apiRes LocationsResponse
	if err := json.Unmarshal(data, &apiRes); err != nil {
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

func getEncounterPokemons(data []byte, config *Config) ([]string, error) {
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
