package command

import (
	"fmt"

	"github.com/Nightails/pokedexcli/internal/config"
)

func commandInspect(conf *config.Config) error {
	if conf.Argument == "" {
		fmt.Println("please specify a pokemon")
	}
	pokemon, exist := conf.Pokemon[conf.Argument]
	if !exist {
		fmt.Println("you have not caught this pokemon yet")
		return nil
	}

	/* Example output:
	Name: pidgey
	Height: 3
	Weight: 18
	Stats:
		-hp: 40
		-attack: 45
		-defense: 40
		-special-attack: 35
		-special-defense: 35
		-speed: 56
	Types:
		- normal
		- flying
	*/

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	displayStats(&pokemon)
	displayTypes(&pokemon)

	return nil
}

// displayStats displays the stats of a pokemon
func displayStats(pokemon *config.Pokemon) {
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
}

// displayTypes displays the types of a pokemon
func displayTypes(pokemon *config.Pokemon) {
	fmt.Println("Types:")
	for _, type_ := range pokemon.Types {
		fmt.Printf("-%s\n", type_.Type.Name)
	}
}
