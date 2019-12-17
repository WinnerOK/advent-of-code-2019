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


func (c *Coordinate) str() string {
	return fmt.Sprintf("{%d %d}", c.X, c.Y)
}
