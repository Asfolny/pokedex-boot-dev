package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	state := state{
		mapStore: pokeMapStore{
			currentPageId: 1,
			pokeMaps:      make(map[int]pokeMap),
		},
	}
	for {
		fmt.Print("Pokedex > ")
		cmd, err := stdin.ReadString('\n')
		if err != nil {
			log.Fatalln("Input parsing broken")
		}
		cmd = strings.ReplaceAll(cmd, "\n", "")

		cliDef, ok := getCommands(state)[cmd]
		if !ok {
			cliDef = getCommands(state)["help"]
		}

		err = cliDef.callback()
		if err != nil {
			fmt.Printf("%v\n", err.Error())
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands(s state) map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 maps",
			callback:    s.mapStore.mapCommand,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the next 20 maps",
			callback:    s.mapStore.mapBackCommand,
		},
	}
}

func commandHelp() error {
	fmt.Println("Usage:\n")
	for _, cliCmd := range getCommands(state{}) {
		fmt.Printf("%s: %s\n", cliCmd.name, cliCmd.description)
	}
	fmt.Println()

	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

type state struct {
	mapStore pokeMapStore
}
