package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Asfolny/pokedex-boot-dev/internal/pokecache"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	startIdx := 1
	state := state{pokecache.New(5 * time.Minute), &startIdx, nil}
	for {
		fmt.Print("Pokedex > ")
		cmd, err := stdin.ReadString('\n')
		if err != nil {
			log.Fatalln("Input parsing broken")
		}
		cmd = strings.ReplaceAll(cmd, "\n", "")
		state.cmdParts = strings.Split(cmd, " ")
		cmd = state.cmdParts[0]

		cliDef, ok := getCommands()[cmd]
		if !ok {
			cliDef = getCommands()["help"]
		}

		err = cliDef.callback(state)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(state) error
}

func getCommands() map[string]cliCommand {
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
			callback:    mapCommand,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the next 20 maps",
			callback:    mapBackCommand,
		},
		"explore": {
			name:        "explore",
			description: "Explore a map",
			callback:    exploreCommand,
		},
	}
}

func commandHelp(state state) error {
	fmt.Print("Usage:\n\n")
	for _, cliCmd := range getCommands() {
		fmt.Printf("%s: %s\n", cliCmd.name, cliCmd.description)
	}
	fmt.Println()

	return nil
}

func commandExit(state state) error {
	os.Exit(0)
	return nil
}

type state struct {
	cache           *pokecache.Cache
	mapCurrentIndex *int
	cmdParts        []string
}
