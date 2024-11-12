package main

import (
	"math/rand"
)

type Player interface {
	Type() string
	Play() int
}

// //////////////////////////
// RandoRex type of player //
// //////////////////////////
type RandoRex struct {
}

// Type returns the type of the player
func (p *RandoRex) Type() string {
	return "RandoRex"
}

// Play returns a move
func (p *RandoRex) Play() int {
	choice := rand.Int() % 3 // 0, 1, or 2
	return choice
}

// //////////////////////////
// Obsessed type of player //
// //////////////////////////
type Obsessed struct {
}

// Type returns the type of the player
func (p *Obsessed) Type() string {
	return "Obsessed"
}

// alwasy return the same value
func (p *Obsessed) Play() int {
	return 1
}

// /////////////////////////
// Flipper type of player //
// /////////////////////////
type Flipper struct {
}

// Type returns the type of the player
func (p *Flipper) Type() string {
	return "Flipper"
}

// random between 0 and 1
func (p *Flipper) Play() int {
	choice := rand.Int() % 2
	return choice
}

// /////////////////////////
// Cyclone type of player //
// /////////////////////////
type Cyclone struct {
	myMove int
}

// Type returns the type of the player
func (p *Cyclone) Type() string {
	return "Cyclone"
}

// cycle between 0, 1 and 2
func (p *Cyclone) Play() int {
	p.myMove = (p.myMove + 1) % 3
	return p.myMove
}
