package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

		newPokeMap := pokeMap{index, "Not Found"}
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
	curMap := getMapByName(state, location)
	for _, poke := range curMap.PokemonEncounters {
		fmt.Printf(" - %v\n", poke.Name)
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

type pokeMap struct {
	Id                int       `json:"id"`
	Name              string    `json:"name"`
	PokemonEncounters []pokemon `json:"pokemon_encounters"`
}

type pokemon struct {
	Name string `json:"name"`
}
