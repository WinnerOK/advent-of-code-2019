package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	fileInput := readInput("./in.txt")
	asteroids := makeAsteroidMap(fileInput)

	var stationCoords Coordinate
	stationView := 0

	for asteroidCoords, _ := range asteroids {
		if view := getVisibleAsteroidsCount(asteroidCoords, copyAsteroidMap(asteroids)); view > stationView {
			stationCoords = asteroidCoords
			stationView = view
		}
	}

	fmt.Printf("Best station is on %v, see %d asteroids\n", stationCoords, stationView)

	slopesSet := map[Coordinate]bool{} // set

	width := len(fileInput[0])
	height := len(fileInput)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x == stationCoords.X && y == stationCoords.Y {
				continue
			} else {
				// Due to float truncation errors, we should work with coordinates,
				// otherwise number of distinct slopes decreases much
				point := Coordinate{x, y}
				slope := smallestDirectionVector(stationCoords, point)
				slopesSet[slope] = true
			}
		}
	}

	println(len(slopesSet)) //Answer is 517
	var slopesSlice []Coordinate
	for k, _ := range slopesSet {
		slopesSlice = append(slopesSlice, k)
	}
	sort.Slice(slopesSlice, func(i, j int) bool {
		comparisonKey := func(dir Coordinate) (int, float64) {
			// Remember that Y-axis is inverted in comparison to Cartesian coordinates
			// So let angle alpha be an angle between a direction vector and y+ axis
			// In inverted coordinates laser rotates CCW (alpha increases), so the slope of direction vector is
			// cot(alpha) = x / y
			switch {
			case dir.Y < 0 && dir.X == 0:
				return 10, 1.
			case dir.Y < 0 && dir.X > 0:
				return 9, math.Abs(float64(dir.Y) / float64(dir.X))
			case dir.X > 0 && dir.Y == 0:
				return 8, 1.
			case dir.X > 0 && dir.Y > 0:
				return 7, math.Abs(float64(dir.X) / float64(dir.Y))
			case dir.Y > 0 && dir.X ==0	:
				return 6,1.
			case dir.Y > 0 && dir.X < 0:
				return 5, float64(dir.Y)/float64(dir.X)
			case dir.Y == 0 && dir.X < 0:
				return 4, 1.
			default:
				return 3, math.Abs(float64(dir.X) / float64(dir.Y))
			}
		}
		arg1Priority, arg1Slope := comparisonKey(slopesSlice[i])
		arg2Priority, arg2Slope := comparisonKey(slopesSlice[j])
		if arg1Priority != arg2Priority {
			return arg1Priority > arg2Priority
		}else{
			return arg1Slope > arg2Slope
		}
	},
	)

	//fmt.Printf("%v\n", slopesSlice)

	cnt := 0
	for true{
		for _, slope := range slopesSlice {
			shotCoords := addCoords(stationCoords, slope)
			for 0 <= shotCoords.X && shotCoords.X < width &&
				0 <= shotCoords.Y && shotCoords.Y < height {
				if _, ok := asteroids[shotCoords]; ok {
					cnt += 1
					delete(asteroids, shotCoords)
					if cnt == 200 {
						fmt.Printf("200th asteroid to be vaporized is %v. Answer: %d",
							shotCoords, shotCoords.X*100+shotCoords.Y)
						return
					}
					break
				}
				shotCoords = addCoords(shotCoords, slope)
			}
		}
	}
}
