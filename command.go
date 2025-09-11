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

	data, err := api.GetPokedexAPI(url)
	if err != nil {
		return err
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
	data, err := api.GetPokedexAPI(url)
	if err != nil {
		return err
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
func getAreaNames(data []byte, config *Config) ([]string, error) {
	var apiRes api.PokedexAPIResponse
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
