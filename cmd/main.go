package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	internal "github.com/asu2sh/pokedex-go/internal/poke"
)

const CACHE_EXPIRATION_DURATION = 10 * time.Second

var pokeCache = internal.NewPokeCache(CACHE_EXPIRATION_DURATION)

type cliCommand struct {
	name        string
	description string
	callback    func(...string)
}

var cliCommandMap map[string]cliCommand

func main() {
	cliCommandMap = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"clear": {
			name:        "clear",
			description: "Clears the Pokedex CLI",
			callback:    clearScreen,
		},
		"help": {
			name:        "help",
			description: "Open the Pokedex Help",
			callback:    helpPokedex,
		},
		"map": {
			name:        "map",
			description: "Get 20 map locations",
			callback:    getNextPokeMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Get 20 previous map locations",
			callback:    getPreviousPokeMap,
		},
		"explore": {
			name:        "explore",
			description: "Explore the Pokedex",
			callback:    explorePokeMap,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback:    catchPokemon,
		},
		"pokedex": {
			name:        "pokedex",
			description: "User's Pokedex",
			callback:    myPokedex,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			text := scanner.Text()
			if len(text) == 0 {
				continue
			}
			textSlice := cleanInput(text)
			cmd, ok := cliCommandMap[textSlice[0]]
			if ok {
				cmd.callback(textSlice[1:]...)
			} else {
				fmt.Println("Unknown command. Type 'help' for available commands.")
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
		}
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	textSlice := strings.Fields(text)
	return textSlice
}

func clearScreen(args ...string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func commandExit(args ...string) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}

func helpPokedex(args ...string) {
	fmt.Println("Welcome to Pokedex!")
	fmt.Printf("Usage:\n\n")
	for key, val := range cliCommandMap {
		fmt.Printf("%v: %v\n", key, val.description)
	}
	fmt.Println("\nAll rights reserved! Created by asu2sh with 💗")
}

func validateArgs(required bool, count int, args ...string) bool {
	if len(args) > count {
		fmt.Println("Ignoring extra arguments!")
		return true
	}
	if required && len(args) < count {
		fmt.Println("Please provide the required arguments!")
		return false
	}
	return true
}

func getNextPokeMap(args ...string) {
	validateArgs(false, 0, args...)
	url := internal.NextPokeMapURL
	internal.GetPokeMap(url, pokeCache)
}

func getPreviousPokeMap(args ...string) {
	validateArgs(false, 0, args...)
	url := internal.PreviousPokeMapURL
	internal.GetPokeMap(url, pokeCache)
}

func explorePokeMap(args ...string) {
	if !validateArgs(true, 1, args...) {
		return
	}
	pokeMapName := args[0]
	internal.ExplorePokeMap(pokeMapName, pokeCache)
}

func catchPokemon(args ...string) {
	if !validateArgs(true, 1, args...) {
		return
	}
	pokemonName := args[0]
	internal.CatchPokemon(pokemonName)
}

func myPokedex(args ...string) {
	if !validateArgs(false, 0, args...) {
		return
	}
	internal.MyPokedex()
}
