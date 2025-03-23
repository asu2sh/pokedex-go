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

var pokeMapName = ""

type cliCommand struct {
	name        string
	description string
	callback    func()
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
			callback:    getPokeMap("map"),
		},
		"mapb": {
			name:        "mapb",
			description: "Get 20 previous map locations",
			callback:    getPokeMap("mapb"),
		},
		"explore": {
			name:        "explore",
			description: "Explore the Pokedex",
			callback:    explorePokeMap,
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
				if cmd.name == cliCommandMap["explore"].name {
					pokeMapName = textSlice[1]
				}
				cmd.callback()
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

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func commandExit() {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}

func helpPokedex() {
	fmt.Println("Welcome to Pokedex!")
	fmt.Printf("Usage:\n\n")
	for key, val := range cliCommandMap {
		fmt.Printf("%v: %v\n", key, val.description)
	}
	fmt.Println("\nAll rights reserved! Created by asu2sh with 💗")
}

func getPokeMap(command string) func() {
	return func() {
		var url string
		if command == "map" {
			url = internal.NextPokeMapURL
		} else {
			url = internal.PreviousPokeMapURL
		}
		internal.GetPokeMap(url, pokeCache)
	}
}

func explorePokeMap() {
	internal.ExplorePokeMap(pokeMapName, pokeCache)
}
