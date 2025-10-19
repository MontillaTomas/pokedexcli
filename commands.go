package main

import (
	"fmt"
	"os"

	"github.com/MontillaTomas/pokedexcli/internal/pokeapi"
)

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
		fmt.Println("Usage:")
		fmt.Println("")
		for _, cmd := range commands {
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		return nil
	}
}

func commandMap(client *pokeapi.Client) (func() error, func() error) {
	var nextURL string
	var prevURL *string

	forward := func() error {
		resp, err := client.GetLocationAreas(nextURL)
		if err != nil {
			return err
		}

		for _, loc := range resp.Results {
			fmt.Println(loc.Name)
		}

		nextURL = resp.Next
		prevURL = resp.Previous
		return nil
	}

	backward := func() error {
		if prevURL == nil || *prevURL == "" {
			fmt.Println("No previous locations to display.")
			return nil
		}

		resp, err := client.GetLocationAreas(*prevURL)
		if err != nil {
			return err
		}

		for _, loc := range resp.Results {
			fmt.Println(loc.Name)
		}

		nextURL = resp.Next
		prevURL = resp.Previous
		return nil
	}

	return forward, backward
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

	client := pokeapi.NewClient(10 * 1e9) // 10 seconds
	forward, backward := commandMap(client)
	commands["map"] = cliCommand{
		name:        "map",
		description: "Displays the names of 20 location areas in the Pokemon world. Each time you run this command, it shows the next set of 20 locations.",
		callback:    forward,
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Displays the previous 20 location areas in the Pokemon world.",
		callback:    backward,
	}

	return commands
}
