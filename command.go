package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Nightails/pokedexcli/internal/api"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*api.Config) error
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

func commandExit(config *api.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *api.Config) error {
	message := `
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex

map: List 20 location names
mapb: List previous 20 location names
`
	fmt.Println(message)
	return nil
}

func commandMap(config *api.Config) error {
	var url string
	if config.NextURL == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = config.NextURL
	}

	areaNames, err := api.GetAreaNames(url, config)
	if err != nil {
		return err
	}

	for _, name := range areaNames {
		fmt.Println(name)
	}

	return nil
}

func commandMapBack(config *api.Config) error {
	if config.PreviousURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	url := config.PreviousURL
	areaNames, err := api.GetAreaNames(url, config)
	if err != nil {
		return err
	}

	for _, name := range areaNames {
		fmt.Println(name)
	}

	return nil
}
