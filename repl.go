package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Nightails/pokedexcli/internal/command"
	"github.com/Nightails/pokedexcli/internal/config"
	"github.com/Nightails/pokedexcli/internal/pokecache"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func startRepl() {
	conf := config.Config{}
	conf.Cache = pokecache.NewCache(5 * time.Second)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		arg := input[0]
		cmd, err := command.GetCommand(arg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(input) == 2 {
			conf.Argument = input[1]
		}

		if err := cmd.Callback(&conf); err != nil {
			log.Fatal(err)
			return
		}
	}
}
