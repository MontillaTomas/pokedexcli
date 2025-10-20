package main

import (
	"github.com/MontillaTomas/pokedexcli/internal/pokeapi"
)

type Pokedex struct {
	Pokemons map[string]pokeapi.Pokemon
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		Pokemons: make(map[string]pokeapi.Pokemon),
	}
}

func (p *Pokedex) Add(pokemon pokeapi.Pokemon) {
	p.Pokemons[pokemon.Name] = pokemon
}

func (p *Pokedex) Get(name string) (pokeapi.Pokemon, bool) {
	pokemon, exists := p.Pokemons[name]
	return pokemon, exists
}
