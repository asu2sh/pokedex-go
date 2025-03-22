package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"github.com/asu2sh/pokedex-go/internal/poke"
)

type cliCommand struct {
	name        string
	description string
	callback    func()
}

var cliCommand_map map[string]cliCommand

var pokecache = internal.NewPokeCache()

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
			callback:    func() { getPokeMap("map")() },
		},
		"mapb": {
			name:        "mapb",
			description: "Get 20 previous map locations",
			callback:    func() { getPokeMap("mapb")() },
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
			cmd, ok := cliCommand_map[text]
			if ok {
				cmd.callback()
			} else {
				cleanInput(text)
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
	if len(text_slice) != 0 {
		fmt.Println("Your command was: ", text_slice[0])
	}
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
		internal.GetPokeMapAPI(url, pokecache)
	}
}
