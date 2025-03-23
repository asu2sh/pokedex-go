package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var PokeMapURL = "https://pokeapi.co/api/v2/location-area/"
var NextPokeMapURL = PokeMapURL
var PreviousPokeMapURL = ""

type PokeMapLocation struct {
	Name string `json:"name"`
}

type PokeMapResult struct {
	Results  []PokeMapLocation `json:"results"`
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
}

func PrintMapLocations(map_results []byte) {
	var poke_map_results PokeMapResult
	json.Unmarshal(map_results, &poke_map_results)

	for _, map_location := range poke_map_results.Results {
		fmt.Println(map_location.Name)
	}

	NextPokeMapURL = poke_map_results.Next
	PreviousPokeMapURL = poke_map_results.Previous
}

func GetPokeMapAPI(url string, pokecache *PokeCache) {
	if url == "" {
		fmt.Println("you're on the first page")
		return
	}

	if map_results, ok := pokecache.Get(url); ok {
		PrintMapLocations(map_results)
		return
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("Map fetching FAILED!")
	}

	map_results, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Map Locations Decode Failure", err)
	}

	PrintMapLocations(map_results)
	pokecache.Add(url, map_results)
}

type Pokemon struct {
	Name string `json:"name"`
}

type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}

type ExplorePokeMapResult struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

func PrintPokemonsLocations(map_results []byte) {
	var poke_map_results ExplorePokeMapResult
	json.Unmarshal(map_results, &poke_map_results)

	for _, pokemon := range poke_map_results.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}
}

func ExplorePokeMapAPI(map_name string, pokecache *PokeCache) {
	fmt.Printf("Exploring %v...\n", map_name)

	if map_results, ok := pokecache.Get(map_name); ok {
		PrintPokemonsLocations(map_results)
		return
	}

	url := fmt.Sprintf("%v%v", PokeMapURL, map_name)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("Map fetching FAILED!")
	}

	map_results, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Map Locations Decode Failure", err)
	}

	PrintPokemonsLocations(map_results)
	pokecache.Add(map_name, map_results)
}
