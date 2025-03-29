package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

const pokemonURL = "https://pokeapi.co/api/v2/pokemon/"
const pokeMapURL = "https://pokeapi.co/api/v2/location-area/"

var NextPokeMapURL = pokeMapURL
var PreviousPokeMapURL = ""

type FieldName struct {
	Name string `json:"name"`
}

type PokeMapResult struct {
	Results  []FieldName `json:"results"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
}

type Stats struct {
	BaseStat int       `json:"base_stat"`
	StatName FieldName `json:"stat"`
}

type Types struct {
	TypeName FieldName `json:"type"`
}

type Pokemon struct {
	Name           string  `json:"name"`
	BaseExperience int     `json:"base_experience"`
	Height         int     `json:"height"`
	Weight         int     `json:"weight"`
	Stats          []Stats `json:"stats"`
	Types          []Types `json:"types"`
}

type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}

type ExplorePokeMapResult struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

type Pokedex struct {
	Pokemons map[string]Pokemon
}

var userPokedex = Pokedex{Pokemons: map[string]Pokemon{}}

func pokeMapAPI(key string, pokeCache *PokeCache) []byte {
	// Try to get the map results from cache first
	if mapResults, ok := pokeCache.Get(key); ok {
		return mapResults
	}

	// Fetch map results if not found in cache
	res, err := http.Get(key)
	if err != nil {
		fmt.Println("error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("Map fetching FAILED!")
	}

	mapResults, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Map Locations Decode Failure", err)
	}

	// Cache the results for future requests
	pokeCache.Add(key, mapResults)
	return mapResults
}

func printMapLocations(mapResults []byte) {
	var pokeMapResults PokeMapResult
	json.Unmarshal(mapResults, &pokeMapResults)

	// Print the names of all map locations
	for _, mapLocation := range pokeMapResults.Results {
		fmt.Println(mapLocation.Name)
	}

	// Update the next and previous URLs for pagination
	NextPokeMapURL = pokeMapResults.Next
	PreviousPokeMapURL = pokeMapResults.Previous
}

func GetPokeMap(url string, pokeCache *PokeCache) {
	if url == "" {
		fmt.Println("You're on the first page")
		return
	}

	mapResults := pokeMapAPI(url, pokeCache)
	printMapLocations(mapResults)
}

func printPokemonsLocations(mapResults []byte) {
	var pokeMapResults ExplorePokeMapResult
	json.Unmarshal(mapResults, &pokeMapResults)

	fmt.Println("Found Pokémon:")

	// Print the names of all Pokémon for the map location
	for _, pokemon := range pokeMapResults.PokemonEncounters {
		fmt.Printf("- %v\n", pokemon.Pokemon.Name)
	}
}

func ExplorePokeMap(mapName string, pokeCache *PokeCache) {
	fmt.Printf("Exploring %v...\n", mapName)

	url := fmt.Sprintf("%v%v", pokeMapURL, mapName)
	mapResults := pokeMapAPI(url, pokeCache)
	printPokemonsLocations(mapResults)
}

func pokemonDetailsAPI(pokemonName string) ([]byte, error) {
	url := fmt.Sprintf("%v%v", pokemonURL, pokemonName)
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return []byte{}, fmt.Errorf("%v doesn't exist", pokemonName)
	}

	pokemonResults, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("pokemon Decode Failure: %w", err)
	}

	return pokemonResults, nil
}

func tryCatchPokemon(pokemonName string, pokemon Pokemon) {
	randomChance := rand.Intn(500)
	fmt.Println(randomChance, pokemon.BaseExperience)

	if randomChance >= pokemon.BaseExperience {
		fmt.Printf("%v was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")
		userPokedex.Pokemons[pokemonName] = pokemon
	} else {
		fmt.Printf("%v escaped!\n", pokemonName)
	}
}

func CatchPokemon(pokemonName string) {
	if _, exists := userPokedex.Pokemons[pokemonName]; exists {
		fmt.Printf("%v already exists in Pokedex!\n", pokemonName)
		return
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)
	pokemonResults, err := pokemonDetailsAPI(pokemonName)

	if err != nil {
		fmt.Println(err)
		return
	}

	var pokemon Pokemon
	err = json.Unmarshal(pokemonResults, &pokemon)

	if err != nil {
		fmt.Println(err)
		return
	}

	tryCatchPokemon(pokemonName, pokemon)
}

func MyPokedex() {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range userPokedex.Pokemons {
		fmt.Printf(" - %v\n", pokemon.Name)
	}
}

func printPokemonDetails(pokemon Pokemon) {
	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, data := range pokemon.Stats {
		fmt.Printf(" - %v: %v\n", data.StatName.Name, data.BaseStat)
	}
	fmt.Println("Types:")
	for _, data := range pokemon.Types {
		fmt.Printf(" - %v\n", data.TypeName.Name)
	}
}

func InspectPokemon(pokemonName string) {
	pokemon_details, exists := userPokedex.Pokemons[pokemonName]
	if !exists {
		fmt.Printf("you have not caught %v\n", pokemonName)
		return
	}
	printPokemonDetails(pokemon_details)
}
