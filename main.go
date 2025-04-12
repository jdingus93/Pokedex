package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

var commands map[string]cliCommand

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
	}
	
	scanner := bufio.NewScanner(os.Stdin)

	for {
	fmt.Print("Pokedex >")
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

func cleanInput(text string) []string {
	fields := strings.Fields(text)
	result := make([]string, len(fields))

	for i, word := range fields {
		result[i] = strings.ToLower(word)
	}
	return result
}