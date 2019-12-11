package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Coordinate struct {
	X         int
	Y         int
	direction int
}

func parseCoordinate(str string) Coordinate {
	coords := strings.Split(str[1:len(str)-1], " ")
	x, _ := strconv.Atoi(coords[0])
	y, _ := strconv.Atoi(coords[1])
	return Coordinate{
		X:         x,
		Y:         y,
		direction: 0,
	}
}

func (c *Coordinate) addCoords(b Coordinate) {
	c.X += b.X
	c.Y += b.Y
}

func (c *Coordinate) rotate(dir int) {
	switch dir {
	case 1:
		c.direction = (c.direction + 1) % 4
	case 0:
		c.direction = (c.direction + 3) % 4
	}
}

func (c *Coordinate) move(step int) {
	switch c.direction {
	case UP:
		c.addCoords(Coordinate{0, step, UP})
	case RIGHT:
		c.addCoords(Coordinate{step, 0, RIGHT})
	case DOWN:
		c.addCoords(Coordinate{0, -step, DOWN})
	case LEFT:
		c.addCoords(Coordinate{-step, 0, LEFT})
	}
}

func (c *Coordinate) str() string {
	return fmt.Sprintf("{%d %d}", c.X, c.Y)
}
