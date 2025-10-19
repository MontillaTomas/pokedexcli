package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type LocationAreaResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
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

func commandMap() (func() error, func() error) {
	url := "https://pokeapi.co/api/v2/location-area/"
	var prevURL *string

	fetch := func(targetURL string) (*LocationAreaResponse, error) {
		res, err := http.Get(targetURL)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		var response LocationAreaResponse
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			return nil, err
		}

		for _, loc := range response.Results {
			fmt.Println(loc.Name)
		}

		return &response, nil
	}

	forward := func() error {
		if url == "" {
			fmt.Println("No more locations to display.")
			return nil
		}

		response, err := fetch(url)
		if err != nil {
			return err
		}

		prevURL = response.Previous
		url = response.Next

		return nil
	}

	backward := func() error {
		if prevURL == nil || *prevURL == "" {
			fmt.Println("No previous locations to display.")
			return nil
		}

		response, err := fetch(*prevURL)
		if err != nil {
			return err
		}

		prevURL = response.Previous
		url = response.Next

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

	forward, backward := commandMap()
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
