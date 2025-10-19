package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	startREPL()
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func startREPL() {
	var commands = initCommands()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)

		if len(cleanedInput) == 0 {
			continue
		}

		commandName := cleanedInput[0]
		command, exists := commands[commandName]

		if !exists {
			fmt.Printf("Unknown command: %s\n", commandName)
			continue
		}

		if err := command.callback(); err != nil {
			fmt.Printf("Error executing command %s: %v\n", commandName, err)
		}
	}
}
