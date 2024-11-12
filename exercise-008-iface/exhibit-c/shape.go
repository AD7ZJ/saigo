package main

import (
	"fmt"
	"math"
)

////////////
// Square //
////////////

type Square struct {
	side float64
}

func (s *Square) Name() string {
	return "Square"
}

func (s *Square) Perimeter() float64 {
	return 4 * s.side
}

func (s *Square) Area() float64 {
	return s.side * s.side
}

////////////
// Circle //
////////////

type Circle struct {
	radius float64
}

func (c *Circle) Name() string {
	return "Circle"
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func (c *Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

//////////////
// Triangle //
//////////////

type Triangle struct {
	side float64
}

func (t *Triangle) Name() string {
	return "Triangle"
}

func (t *Triangle) Perimeter() float64 {
	return 3 * t.side
}

func (t *Triangle) Area() float64 {
	return 0.5 * t.side * t.side
}

/////////////
// Hexagon //
/////////////

type Hexagon struct {
	side float64
}

func (h *Hexagon) Name() string {
	return "Hexagon"
}

func (h *Hexagon) Perimeter() float64 {
	return 6 * h.side
}

// Area calculates the area of the hexagon using the formula (3 * √3 / 2) * side ^ 2
func (h *Hexagon) Area() float64 {
	return (3 * math.Sqrt(3) / 2) * h.side * h.side
}

////////////////
// Efficiency //
////////////////

type Shape interface {
	Name() string
	Perimeter() float64
	Area() float64
}

func Efficiency(s Shape) {
	name := s.Name()
	area := s.Area()
	rope := s.Perimeter()

	efficiency := 100.0 * area / (rope * rope)
	fmt.Printf("Efficiency of a %s is %f\n", name, efficiency)
}

// Build function to create an instance of type Shape. Every return value should implement the Shape interface
func Build(shape string, parameters ...float64) Shape {
	switch shape {
	case "Triangle":
		if len(parameters) < 1 {
			fmt.Println("Not enough parameters for Triangle")
			return nil
		}
		return &Triangle{side: parameters[0]}

	case "Square":
		if len(parameters) < 1 {
			fmt.Println("Not enough parameters for Square")
			return nil
		}
		return &Square{side: parameters[0]}

	case "Hexagon":
		if len(parameters) < 1 {
			fmt.Println("Not enough parameters for Hexagon")
			return nil
		}
		return &Hexagon{side: parameters[0]}

	case "Circle":
		if len(parameters) < 1 {
			fmt.Println("Not enough parameters for Circle")
			return nil
		}
		return &Circle{radius: parameters[0]}

	default:
		fmt.Println("Unknown shape type")
		return nil
	}
}

func main() {
	//s := Square{side: 10.0}
	Efficiency(Build("Square", 10.0))

	//c := Circle{radius: 10.0}
	Efficiency(Build("Circle", 10.0))

	//t := Triangle{side: 10.0}
	Efficiency(Build("Triangle", 10.0))

	//h := Hexagon{side: 10.0}
	Efficiency(Build("Hexagon", 10.0))
}
