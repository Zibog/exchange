package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const pokeapiURL = "http://pokeapi.co/api/v2/pokedex/kanto/"

func CallPokeapi() PokemonResponse {
	response, err := http.Get(pokeapiURL)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject PokemonResponse
	json.Unmarshal(responseData, &responseObject)

	// fmt.Println(responseObject.Name)
	// fmt.Println(len(responseObject.Pokemon))

	// for _, pokemon := range responseObject.Pokemon {
	// 	fmt.Println(pokemon.Species.Name)
	// }

	return responseObject
}

// A PokemonResponse struct to map the Entire PokemonResponse
type PokemonResponse struct {
	Name    string    `json:"name"`
	Pokemon []Pokemon `json:"pokemon_entries"`
}

// A Pokemon Struct to map every pokemon to.
type Pokemon struct {
	EntryNo int            `json:"entry_number"`
	Species PokemonSpecies `json:"pokemon_species"`
}

// A struct to map our Pokemon's Species which includes it's name
type PokemonSpecies struct {
	Name string `json:"name"`
}
