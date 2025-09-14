package command

import (
	"fmt"

	"github.com/Nightails/pokedexcli/internal/config"
)

func commandHelp(config *config.Config) error {
	message := `
Welcome to the Pokedex!
Usage:
> help: Displays a help message
> exit: Exit the Pokedex

> map: List 20 location names
> mapb: List previous 20 location names

> explore: Explore the area

`
	fmt.Printf("%v", message)
	return nil
}
