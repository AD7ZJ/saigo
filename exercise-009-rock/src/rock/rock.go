package main

func main() {

	// Create Game
	game := &Game{}

	// Add Players
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
