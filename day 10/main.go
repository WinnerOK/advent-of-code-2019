package main

import "fmt"

func main() {
	fileInput := readInput("./in.txt")
	asteroids := makeAsteroidMap(fileInput)

	var bestCoord Coordinate
	bestView := 0

	for asteroidCoords, _ := range asteroids {
		if view := getVisibleAsteroidsCount(asteroidCoords, copyAsteroidMap(asteroids)); view > bestView {
			bestCoord = asteroidCoords
			bestView = view
		}
	}

	fmt.Printf("Best station is on %v, see %d asteroids\n", bestCoord, bestView)

}
