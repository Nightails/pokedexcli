package command

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/Nightails/pokedexcli/internal/api"
	"github.com/Nightails/pokedexcli/internal/config"
)

func commandCatch(conf *config.Config) error {
	if conf.Argument == "" {
		fmt.Println("please specify which pokemon to catch")
		return nil
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + conf.Argument
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

	fmt.Printf("Throwing a Pokeball at %s...\n", conf.Argument)
	pokemon, err := getPokemon(data)
	if err != nil {
		return err
	}
	rate := 1000 - pokemon.BaseExperience
	randNum := rand.IntN(1000)
	if randNum <= rate {
		fmt.Printf("%s was caught!\n", conf.Argument)
		conf.Pokemon[conf.Argument] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", conf.Argument)
	}

	return nil
}

// getPokemonCatchRate returns the catch rate of the given pokemon
func getPokemon(data []byte) (config.Pokemon, error) {
	var pokemon config.Pokemon
	if err := json.Unmarshal(data, &pokemon); err != nil {
		return config.Pokemon{}, err
	}
	return pokemon, nil
}
