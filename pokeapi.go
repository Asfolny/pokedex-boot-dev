package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const mapSize = 20

func (store *pokeMapStore) mapCommand() error {
	for i := store.currentPageId; i < store.currentPageId+mapSize; i++ {
		curMap := store.get(i)
		fmt.Printf("%v\n", curMap.Name)
	}

	store.currentPageId += mapSize

	return nil
}

func (store *pokeMapStore) mapBackCommand() error {
	if store.currentPageId < mapSize {
		return errors.New("Already on first page")
	}
	store.currentPageId -= mapSize

	for i := store.currentPageId; i < store.currentPageId+mapSize; i++ {
		curMap := store.get(i)
		fmt.Printf("%v\n", curMap.Name)
	}

	return nil
}

func (store *pokeMapStore) get(index int) pokeMap {
	item, ok := store.pokeMaps[index]
	if !ok {
		res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", index))
		if err != nil {
			log.Fatalln(err)
		}

		if res.Status == "404" {
			item = pokeMap{index, "Not Found"}
			store.pokeMaps[index] = item
			return item
		}

		item = pokeMap{}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&item); err != nil {
			log.Fatal(err)
		}

		store.pokeMaps[index] = item
	}

	return item
}

type pokeMap struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type pokeMapStore struct {
	currentPageId int
	pokeMaps      map[int]pokeMap
}
