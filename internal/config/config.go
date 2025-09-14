package config

import "github.com/Nightails/pokedexcli/internal/pokecache"

type Config struct {
	NextURL     string
	PreviousURL string
	Argument    string
	Cache       *pokecache.Cache
}
