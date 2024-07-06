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
	for {
		fmt.Print("Pokedex > ")
		cmd, err := stdin.ReadString('\n')
		if err != nil {
			log.Fatalln("Input parsing broken")
		}
		cmd = strings.ReplaceAll(cmd, "\n", "")

		cliDef, ok := getCommands()[cmd]
		if !ok {
			cliDef = getCommands()["help"]
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
	}
}

func commandHelp() error {
	fmt.Println("Usage:\n")
	for _, cliCmd := range getCommands() {
		fmt.Printf("%s: %s\n", cliCmd.name, cliCmd.description)
	}
	fmt.Println()

	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

