package command

import (
	"fmt"

	"github.com/Nightails/pokedexcli/internal/config"
)

func commandPokedex(conf *config.Config) error {
	if len(conf.Pokemon) == 0 {
		fmt.Println("no pokemon in the pokedex")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name, _ := range conf.Pokemon {
		fmt.Printf("- %s\n", name)
	}

	return nil
}
