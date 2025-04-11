package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("'Hello, World!'")
}

func cleanInput(text string) []string {
	fields := strings.Fields(text)
	result := make([]string, len(fields))

	for i, word := range fields {
		result[i] = strings.ToLower(word)
	}
	return result
}