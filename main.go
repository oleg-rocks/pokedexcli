package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/oleg-rocks/pokedexcli/internal/pokeapi"
)

var registry map[string]cliCommand

func init() {
	registry = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedox",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations",
			callback:    commandMapb,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := Config{}

	for {
		fmt.Print("Pokedox > ")
		scanner.Scan()
		input := scanner.Text()
		hasCommand := false
		for key, value := range registry {
			if input == key {
				err := value.callback(&config)
				if err != nil {
					fmt.Println("Error: ", err)
				}
				hasCommand = true
				break
			}
		}
		if !hasCommand {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	lowText := strings.ToLower(text)
	words := strings.Fields(lowText)
	return words
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config) error
}

type Config struct {
	Next     *string
	Previous *string
}

func commandExit(config *Config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, command := range registry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *Config) error {
	resp, err := pokeapi.MakeLocationsRequest(config.Next)
	if err != nil {
		return err
	}

	config.Next = resp.Next
	config.Previous = resp.Previous

	printAreas(resp.Results)
	return nil
}

func commandMapb(config *Config) error {
	resp, err := pokeapi.MakeLocationsRequest(config.Previous)
	if err != nil {
		return err
	}

	config.Next = resp.Next
	config.Previous = resp.Previous

	printAreas(resp.Results)
	return nil
}

func printAreas(results []pokeapi.LocationAreaResult) {
	for _, area := range results {
		fmt.Println(area.Name)
	}
}
