package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"math/rand"
	"errors"

	"github.com/jdingus/Pokedex/pokeapi"
)

var cache = make(map[string][]byte)

type cliCommand struct {
	name string
	description string
	callback func(param string) error
}

var commands = make(map[string]cliCommand)
var config = struct {
	Next 	*string
	Previous	*string
	}{
		Next: nil,
		Previous: nil,
	}

var caughtPokemon = make(map[string]Pokemon)
type Pokemon struct {
	Name string
	Height int
	Weight int
	Stats []Stat
	Types []string
	BaseExperience int
}

type Stat struct {
	Name string
	Value int
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
		"explore": {
			name: "explore",
			description: "Shows what Pokemon are in the area",
			callback: commandExplore,
		},
		"catch": {
			name: "catch",
			description: "Attempt to catch a Pokemon",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect",
			description: "Inspect a pokemon's stats",
			callback: commandInspect,
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
			var param string
			if len(words) > 1 {
				param = words[1]
			}
			err := command.callback(param)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unkown command")
		}
	}
}

func commandExit(param string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(param string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(param string) error {
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

func commandMapBack(param string) error {
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

func commandExplore(param string) error {
	if param == "" {
		return fmt.Errorf("please provide a location area name")
	}
	fmt.Printf("Exploring %s...\n", param)

	var locationArea pokeapi.LocationAreaDetail 

	if cachedData, ok := cache[param]; ok {
		err := json.Unmarshal(cachedData, &locationArea)
		if err != nil {
			return err
		}
	} else {
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", param)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	cache[param] = body

	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return err
		}
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(param string) error {
	args := strings.Split(param, " ")
	if len(args) == 0 {
		return errors.New("you must provide a pokemon name")
	}
	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("pokemon %s not found", pokemonName)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var pokemonData struct {
		BaseExperience int `json:"base_experience"`
		Height int `json:"height"`
		Weight int `json:"weight"`
		Name string `json:"name"`
		Stats []struct {
			BaseStat int `json:"base_stat"`
			Stat struct {
				Name string `json:"name"`
			} `json:"stat"`
		} `json:"stats"`
		Types []struct {
			Type struct {
				Name string `json:"name"`
			 } `json:"type"`
		} `json:"types"`
	}

	if err := json.Unmarshal(body, &pokemonData); err != nil {
		return err
	}

	baseExperience := pokemonData.BaseExperience
	catchProbability := 100.0 / float64(baseExperience)
	caught := rand.Float64() < catchProbability

	if caught {
		fmt.Printf("%s was caught!\n", pokemonName)

		var stats []Stat
		for _, statData := range pokemonData.Stats {
			stats = append(stats, Stat{
				Name: statData.Stat.Name,
				Value: statData.BaseStat,
			})
		}
		var types []string
		for _, typeData := range pokemonData.Types {
			types = append(types, typeData.Type.Name)
		}

		caughtPokemon[pokemonName] = Pokemon{
			Name: pokemonName,
			Height: pokemonData.Height,
			Weight: pokemonData.Weight,
			Stats: stats,
			Types: types,
			BaseExperience: pokemonData.BaseExperience,
		}
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	
	return nil
}

func commandInspect(param string) error {
	args := strings.Split(param, " ")
	if len(args) != 1 {
		return fmt.Errorf("you must provide a pokemon name")
	}
	
	pokemonName := args[0]

	pokemon, ok := caughtPokemon[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("	-%s: %d\n", stat.Name, stat.Value)
	}
	fmt.Println("Types:")
	for _, typeName := range pokemon.Types {
		fmt.Printf("	-%s\n", typeName)
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