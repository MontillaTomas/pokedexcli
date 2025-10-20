package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"

	"github.com/MontillaTomas/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

func commandExit(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commands map[string]cliCommand) func(args []string) error {
	return func(args []string) error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Println("Usage:")
		fmt.Println("")

		names := make([]string, 0, len(commands))
		for name := range commands {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			cmd := commands[name]
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		return nil
	}
}

func commandMap(client *pokeapi.Client) (func(args []string) error, func(args []string) error) {
	var nextURL string
	var prevURL *string

	forward := func(args []string) error {
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

	backward := func(args []string) error {
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

func commandExploreLocationArea(client *pokeapi.Client) func(args []string) error {
	return func(args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("location area name is required")
		}
		locationAreaName := args[0]

		details, err := client.GetLocationAreaPokemons(locationAreaName)
		if err != nil {
			return err
		}

		fmt.Printf("Exploring %s...\n", details.Name)
		fmt.Println("Found Pokemon:")
		for _, encounter := range details.PokemonEncounters {
			fmt.Println("- ", encounter.Pokemon.Name)
		}

		return nil
	}
}

func commandCatch(client *pokeapi.Client, pokedex *Pokedex) func(args []string) error {
	return func(args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("pokemon name is required")
		}
		pokemonName := args[0]

		pokemon, err := client.GetPokemon(pokemonName)
		if err != nil {
			return err
		}

		fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

		catchChance := 100 - pokemon.BaseExperience/2
		if catchChance < 10 {
			catchChance = 10
		}
		if rand.Intn(100) < catchChance {
			pokedex.Add(*pokemon)
			fmt.Printf("%s was caught!\n", pokemon.Name)
		} else {
			fmt.Printf("%s escaped!\n", pokemon.Name)
		}

		return nil
	}
}

func commandInspect(pokedex *Pokedex) func(args []string) error {
	return func(args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("pokemon name is required")
		}
		pokemonName := args[0]

		pokemon, exists := pokedex.Get(pokemonName)
		if !exists {
			return fmt.Errorf("you have not caught that pokemon")
		}

		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Printf("Stats:\n")
		for _, stat := range pokemon.Stats {
			fmt.Printf("\t-%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _, t := range pokemon.Types {
			fmt.Printf("\t-%s\n", t.Type.Name)
		}

		return nil
	}
}

func initCommands() map[string]cliCommand {
	commands := make(map[string]cliCommand)
	client := pokeapi.NewClient(10 * 1e9) // 10 seconds
	pokedex := NewPokedex()

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
	commands["explore"] = cliCommand{
		name:        "explore <location_area_name>",
		description: "Explore a specific location area by its name to see the Pokemon that can be encountered there.",
		callback:    commandExploreLocationArea(client),
	}
	commands["catch"] = cliCommand{
		name:        "catch <pokemon_name>",
		description: "Attempt to catch a Pokemon by its name.",
		callback:    commandCatch(client, pokedex),
	}
	commands["inspect"] = cliCommand{
		name:        "inspect <pokemon_name>",
		description: "Inspect a caught Pokemon to see its details.",
		callback:    commandInspect(pokedex),
	}

	return commands
}
