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

type cliCommand struct {
	name        string
	description string
	callback    func()
}

var cliCommand_map map[string]cliCommand

var CacheExpirationDuration = 10 * time.Second
var pokecache = internal.NewPokeCache(CacheExpirationDuration)

var poke_map_name = ""

func main() {
	cliCommand_map = map[string]cliCommand{
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
			text_slice := cleanInput(text)
			cmd, ok := cliCommand_map[text_slice[0]]
			if ok {
				if cmd.name == cliCommand_map["explore"].name {
					poke_map_name = text_slice[1]
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
	text_slice := strings.Fields(text)
	return text_slice
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
	for key, val := range cliCommand_map {
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
		internal.GetPokeMap(url, pokecache)
	}
}

func explorePokeMap() {
	internal.ExplorePokeMap(poke_map_name, pokecache)
}
