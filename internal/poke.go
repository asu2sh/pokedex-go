package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var NextPokeMapURL = "https://pokeapi.co/api/v2/location-area/"
var PreviousPokeMapURL = ""

type PokeMapLocation struct {
	Name string `json:"name"`
}

type PokeMapResult struct {
	Results  []PokeMapLocation `json:"results"`
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
}

func GetPokeMapAPI(url string) {
	if url == "" {
		fmt.Println("you're on the first page")
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

	for _, map_location := range map_results.Results {
		fmt.Println(map_location.Name)
	}

	NextPokeMapURL = map_results.Next
	PreviousPokeMapURL = map_results.Previous
}
