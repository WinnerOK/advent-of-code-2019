package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type Coordinate struct {
	X int
	Y int
	Z int
}

var CoordPattern = regexp.MustCompile(`^<x=(-?\d+), y=(-?\d+), z=(-?\d+)>$`)

func parseCoordinate(str string) Coordinate {

	groups := CoordPattern.FindAllStringSubmatch(str, -1)
	x, _ := strconv.Atoi(groups[0][1])
	y, _ := strconv.Atoi(groups[0][2])
	z, _ := strconv.Atoi(groups[0][3])
	return Coordinate{
		X: x,
		Y: y,
		Z: z,
	}
}

func (c *Coordinate) addCoords(b Coordinate) {
	c.X += b.X
	c.Y += b.Y
	c.Z += b.Z

}

func (c *Coordinate) coordsSlice() []int {
	return []int{c.X, c.Y, c.Z}
}

func (c *Coordinate) str() string {
	return fmt.Sprintf("{x=%d y=%d z=%d}", c.X, c.Y, c.Z)
}
