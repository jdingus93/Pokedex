package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
	fmt.Print("Pokedex >")
	scanner.Scan()
	userInput := scanner.Text()
	words := strings.Fields(userInput)
	if len(words) > 0 {
		fmt.Println("Your command was:", strings.ToLower(words[0]))
		}
	}
}

func cleanInput(text string) []string {
	fields := strings.Fields(text)
	result := make([]string, len(fields))

	for i, word := range fields {
		result[i] = strings.ToLower(word)
	}
	return result
}