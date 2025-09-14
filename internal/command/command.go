package command

import (
	"errors"

	"github.com/Nightails/pokedexcli/internal/config"
)

type CliCommand struct {
	name        string
	description string
	Callback    func(*config.Config) error
}

func GetCommand(cmd string) (CliCommand, error) {
	commands := map[string]CliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Show how to use Pokedex",
			Callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Show 20 location names",
			Callback:    commandMap,
		},
		"mapb": {
			name:        "map back",
			description: "Show the previous 20 location names",
			Callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore the area",
			Callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try catching a pokemon",
			Callback:    commandCatch,
		},
	}

	if c, exist := commands[cmd]; exist {
		return c, nil
	}
	return CliCommand{}, errors.New("unknown command")
}
