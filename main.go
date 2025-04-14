package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"github.com/jdingus/Pokedex/pokeapi"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

var commands = make(map[string]cliCommand)
var config = struct {
	Next 	*string
	Previous	*string
	}{
		Next: nil,
		Previous: nil,
	}

func main() {

	var commands = map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"map": {
			name:	"map",
			description: 	"Fetch and display the next 20 location areas",
			callback: 	commandMap,
		},
		"mapb": {
			name: 	"mapb",
			description: 	"Fetch and display the prevoius 20 location areas",
			callback: 	commandMapBack,
		},
	}
	
	fmt.Println("Welcome to the Pokedex!")
	scanner := bufio.NewScanner(os.Stdin)

	for {
	fmt.Print("Pokedex > ")
	scanner.Scan()
	userInput := scanner.Text()
	
	words := cleanInput(userInput)
	if len(words) == 0 {
		continue
		}

		commandName := words[0]
		if command, exists := commands[commandName]; exists {
			err := command.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unkown command")
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap() error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if config.Next != nil {
		url = *config.Next
	}
	fmt.Println("fgetching from URL:" , url)

	data, err := pokeapi.FetchLocationAreas(url)
	if err != nil {
		return fmt.Errorf("error fetching location areas: %w", err)
	}
	fmt.Printf("got response: Nest=%v, Previous=%v, Results=%d\n",
				data.Next, data.Previous, len(data.Results))

	fmt.Println("location areas:")
	if len(data.Results) == 0 {
		fmt.Println("no locations found!")
	}
	for i, location := range data.Results {
		fmt.Printf("%d. %s\n", i+1, location.Name)
	}
	config.Next = data.Next
	config.Previous = data.Previous
	return nil
}

func commandMapBack() error {
	if config.Previous == nil {
		fmt.Println("You're on the first page.")
		return nil
	}

	data, err := pokeapi.FetchLocationAreas(*config.Previous)
	if err != nil {
		return fmt.Errorf("error fetching location areas: %w", err)
	}
	
	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	config.Next = data.Next
	config.Previous = data.Previous
	return nil
}


func cleanInput(text string) []string {
	fields := strings.Fields(text)
	result := make([]string, len(fields))

	for i, word := range fields {
		result[i] = strings.ToLower(word)
	}
	return result
}