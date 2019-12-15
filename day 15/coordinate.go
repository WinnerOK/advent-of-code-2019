package main

import (
	"fmt"
)

type Coordinate struct {
	X int
	Y int
}

func (c *Coordinate) addCoords(b Coordinate) {
	c.X += b.X
	c.Y += b.Y
}

func (c *Coordinate) move(direction int) {
	switch direction {
	case NORTH:
		c.addCoords(Coordinate{Y: -1})
	case SOUTH:
		c.addCoords(Coordinate{Y: 1})
	case WEST:
		c.addCoords(Coordinate{X: -1})
	case EAST:
		c.addCoords(Coordinate{X: 1})
	}
}

func (c *Coordinate) str() string {
	return fmt.Sprintf("{%d %d}", c.X, c.Y)
}
