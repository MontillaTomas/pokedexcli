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

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commands map[string]cliCommand) func() error {
	return func() error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Println("Usage:\n")
		for _, cmd := range commands {
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		return nil
	}
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

func initCommands() map[string]cliCommand {
	commands := make(map[string]cliCommand)
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp(commands),
	}

	return commands
}
