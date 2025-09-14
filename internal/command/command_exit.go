package command

import (
	"fmt"
	"os"

	"github.com/Nightails/pokedexcli/internal/config"
)

func commandExit(config *config.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
