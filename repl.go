package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Nightails/pokedexcli/internal/api"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func startRelp() {
	config := api.Config{}
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
		cmd.callback(&config)
	}
}
