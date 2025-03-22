package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var NextPokeMapURL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
var PreviousPokeMapURL = ""

type PokeMapLocation struct {
	Name string `json:"name"`
}

type PokeMapResult struct {
	Results  []PokeMapLocation `json:"results"`
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
}

func PrintMapLocations(map_results PokeMapResult) {
	for _, map_location := range map_results.Results {
		fmt.Println(map_location.Name)
	}

	NextPokeMapURL = map_results.Next
	PreviousPokeMapURL = map_results.Previous
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

	var map_results PokeMapResult

	if err := json.NewDecoder(res.Body).Decode(&map_results); err != nil {
		fmt.Println("Map Locations Decode Failure", err)
	}

	PrintMapLocations(map_results)
	pokecache.Add(url, map_results, 5*time.Second)
}
