package main

import (
	"math/rand"
)

type Player interface {
	Type() string
	Play() int
}

// RandoRex type of player
type RandoRex struct {
}

// Type returns the type of the player
func (p *RandoRex) Type() string {
	return "RandoRex"
}

// Play returns a move
func (p *RandoRex) Play() int {
	choice := rand.Int() % 3
	return choice
}
