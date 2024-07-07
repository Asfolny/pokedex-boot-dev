package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const mapSize = 20

func mapCommand(state state) error {
	for i := *state.mapCurrentIndex; i < *state.mapCurrentIndex+mapSize; i++ {
		curMap := getMap(state, i)
		fmt.Printf("%v\n", curMap.Name)
	}

	*state.mapCurrentIndex += mapSize

	return nil
}

func mapBackCommand(state state) error {
	if *state.mapCurrentIndex < mapSize {
		return errors.New("Already on first page")
	}
	*state.mapCurrentIndex -= mapSize

	for i := *state.mapCurrentIndex; i < *state.mapCurrentIndex+mapSize; i++ {
		curMap := getMap(state, i)
		fmt.Printf("%v\n", curMap.Name)
	}

	return nil
}

func getMap(state state, index int) pokeMap {
	cacheKey := fmt.Sprintf("pokeapi-map-%v", index)
	cacheItem, ok := state.cache.Get(cacheKey)
	if !ok {
		res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", index))
		if err != nil {
			log.Fatalln(err)
		}

		newPokeMap := pokeMap{index, "Not Found", nil}
		if res.StatusCode == 200 {
			decoder := json.NewDecoder(res.Body)
			if err := decoder.Decode(&newPokeMap); err != nil {
				log.Fatalln(err)
			}
		}

		cacheItem, err = json.Marshal(newPokeMap)
		if err != nil {
			log.Fatalln(err)
		}

		state.cache.Add(cacheKey, cacheItem)
	}

	var item pokeMap
	err := json.Unmarshal(cacheItem, &item)
	if err != nil {
		log.Fatalln(err)
	}

	return item
}

func exploreCommand(state state) error {
	location := state.cmdParts[1]
	fmt.Printf("Exploring %v...\n", location)
	curMap := getMapByName(state, location)
	fmt.Println("Found Pokemon:")
	for i, pokeEncounter := range curMap.PokemonEncounters {
		fmt.Printf("%v - %v\n", i, pokeEncounter.Pokemon.Name)
	}
	return nil
}

func getMapByName(state state, name string) pokeMap {
	cacheKey := fmt.Sprintf("pokeapi-map-%v", name)
	cacheItem, ok := state.cache.Get(cacheKey)
	if !ok {
		res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", name))
		if err != nil {
			log.Fatalln(err)
		}

		var newPokeMap pokeMap
		if res.StatusCode == 200 {
			decoder := json.NewDecoder(res.Body)
			if err := decoder.Decode(&newPokeMap); err != nil {
				log.Fatalln(err)
			}
		}

		cacheItem, err = json.Marshal(newPokeMap)
		if err != nil {
			log.Fatalln(err)
		}

		state.cache.Add(cacheKey, cacheItem)
	}

	var item pokeMap
	err := json.Unmarshal(cacheItem, &item)
	if err != nil {
		log.Fatalln(err)
	}

	return item
}

func catchCommand(state state) error {
	pokemonName := state.cmdParts[1]
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)
	pokemon := getPokemonByName(state, pokemonName)
	rander := rand.New(rand.NewSource(time.Now().UnixNano()))
	if rander.Intn(pokemon.BaseExp*2) < pokemon.BaseExp {
		fmt.Printf("%v escaped!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%v was caught!\n", pokemon.Name)
	state.pokedex[pokemon.Name] = pokemon
	return nil
}

func inspectCommand(state state) error {
	name := state.cmdParts[1]
	pokemon, ok := state.pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)

	fmt.Print("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("    -%v:%v\n", stat.Stat.Name, stat.StatBase)
	}

	fmt.Print("Types:\n")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("    -%v\n", pokemonType.Type.Name)
	}

	return nil
}

func getPokemonByName(state state, name string) pokemon {
	cacheKey := fmt.Sprintf("pokeapi-pokemon-%v", name)
	cacheItem, ok := state.cache.Get(cacheKey)
	if !ok {
		res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", name))
		if err != nil {
			log.Fatalln(err)
		}

		var newPokemon pokemon
		if res.StatusCode == 200 {
			decoder := json.NewDecoder(res.Body)
			if err := decoder.Decode(&newPokemon); err != nil {
				log.Fatalln(err)
			}
		}

		cacheItem, err = json.Marshal(newPokemon)
		if err != nil {
			log.Fatalln(err)
		}

		state.cache.Add(cacheKey, cacheItem)
	}

	var item pokemon
	err := json.Unmarshal(cacheItem, &item)
	if err != nil {
		log.Fatalln(err)
	}

	return item
}

type pokeMap struct {
	Id                int                `json:"id"`
	Name              string             `json:"name"`
	PokemonEncounters []pokemonEncounter `json:"pokemon_encounters"`
}

type pokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
	} `json:"pokemon"`
}

type pokemon struct {
	Name    string       `json:"name"`
	BaseExp int          `json:"base_experience"`
	Height  int          `json:"height"`
	Weight  int          `json:"weight"`
	Stats   []statStruct `json:"stats"`
	Types   []typeStruct `json:"types"`
}

type statStruct struct {
	StatBase int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type typeStruct struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}
