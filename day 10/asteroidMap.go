package main

import (
	"math"
)

const (
	Asteroid = '#'
	Space    = '.'
)

type Coordinate struct {
	X int
	Y int
}

func addCoords(a, b Coordinate) Coordinate {
	return Coordinate{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

func subCoords(a, b Coordinate) Coordinate {
	return Coordinate{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}

func scalarMulCoords(a Coordinate, i int) Coordinate {
	return Coordinate{
		X: a.X * i,
		Y: a.Y * i,
	}
}

func getAsteroidMapBoundaries(m AsteroidMap) (int, int, int, int) {
	xMax := math.MinInt64
	yMax := math.MinInt64

	xMin := math.MaxInt64
	yMin := math.MaxInt64

	for k, _ := range m {
		if k.X > xMax {
			xMax = k.X
		}

		if k.X < xMin {
			xMin = k.X
		}

		if k.Y > yMax {
			yMax = k.Y
		}

		if k.Y < yMin {
			yMin = k.Y
		}
	}

	return xMin, yMin, xMax, yMax
}

func smallestDirectionVector(direction Coordinate) Coordinate {
	var (
		xDir = 1
		yDir = 1
	)

	if direction.X < 0 {
		xDir = -1
	}

	if direction.Y < 0 {
		yDir = -1
	}

	gcd := GCD(direction.X*xDir, direction.Y*yDir) // make both coordinates positive
	return Coordinate{  direction.X / gcd, direction.Y / gcd}
}

func getVisibleAsteroidsCount(current Coordinate, asteroidMap AsteroidMap) int {
	xMin, yMin, xMax, yMax := getAsteroidMapBoundaries(asteroidMap)
	for asteroidCoords, visible := range asteroidMap {
		if asteroidCoords == current {
			continue
		}
		if visible {
			directionVector := smallestDirectionVector(subCoords(asteroidCoords, current))
			canBlock := false
			for i := 1; true; i++ {
				invisibleAsteroidCoords := addCoords(current, scalarMulCoords(directionVector, i))
				if invisibleAsteroidCoords.X < xMin || invisibleAsteroidCoords.Y < yMin || // make sure you are not
					invisibleAsteroidCoords.X > xMax || invisibleAsteroidCoords.Y > yMax { // out of bounds
					break
				} else if invisibleAsteroidCoords == asteroidCoords { // only from now the asteroid can block others
					canBlock = true
					continue
				} else {
					if visible, ok := asteroidMap[invisibleAsteroidCoords]; ok && visible && canBlock {
						asteroidMap[invisibleAsteroidCoords] = false
						//fmt.Printf("From %v: %v blockes %v\n", current, asteroidCoords, invisibleAsteroidCoords)
					}
				}
			}
		}
	}
	visibleCnt := 0
	for _, visible := range asteroidMap {
		if visible {
			visibleCnt += 1
		}
	}
	return visibleCnt - 1
}

type AsteroidMap map[Coordinate]bool

func makeAsteroidMap(mapStr []string) AsteroidMap {
	resultMap := AsteroidMap{}
	for Y, row := range mapStr {
		for X, entity := range row {
			if entity == Asteroid {
				resultMap[Coordinate{X, Y}] = true
			}
		}
	}
	return resultMap
}

func copyAsteroidMap(m AsteroidMap) AsteroidMap {
	newMap := AsteroidMap{}
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
