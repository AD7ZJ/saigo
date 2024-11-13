package main

func main() {

	// Create Game
	game := &Game{}

	// Add Players. The syntax here is a bit odd, since Add() is defined to take a Player type, not a Player pointer type.
	// But since the interface is using pointer reveivers on the Type() and Play() function, we need to pass in a
	// pointer here. Go allows a lot of leeway in matching pointers vs the dereferenced value.
	game.Add(&RandoRex{})
	game.Add(&Obsessed{})
	game.Add(&Flipper{})
	game.Add(&Cyclone{})

	// A Thousand Round-Robins!
	for i := 0; i < 1000; i++ {
		game.RoundRobin()
	}

	// Display Results
	game.Display()
}
