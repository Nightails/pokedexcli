package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Nightails/pokedexcli/internal/pokecache"
)

type Config struct {
	NextURL     string
	PreviousURL string
	Argument    string
	Cache       *pokecache.Cache
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func startRepl() {
	config := Config{}
	config.Cache = pokecache.NewCache(5 * time.Second)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		command := input[0]
		cmd, err := getCommand(command)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(input) == 2 {
			config.Argument = input[1]
		}

		if err := cmd.callback(&config); err != nil {
			log.Fatal(err)
			return
		}
	}
}
