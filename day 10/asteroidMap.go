package main

const (
	Asteroid = '#'
	Space    = '.'
)

type Coordinate struct {
	X     int
	Y     int
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

func smallestDirectionVector(from, to Coordinate) Coordinate {
	direction := subCoords(to, from)

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
	return Coordinate{direction.X / gcd, direction.Y / gcd}
}

func isDirectlySeen(from, to Coordinate, asteroidMap AsteroidMap) bool {
	directionVector := smallestDirectionVector(from, to)
	reachOther := false
	invisibleAsteroidCoords := addCoords(from, directionVector)
	for ; invisibleAsteroidCoords != to; // increase direction vector unless hit desired asteroid
	invisibleAsteroidCoords = addCoords(invisibleAsteroidCoords, directionVector) {
		if _, ok := asteroidMap[invisibleAsteroidCoords]; ok {
			// if hit some other one, then the desired can't be seen directly
			reachOther = true
			break
		}
	}
	return !reachOther
}

func getVisibleAsteroidsCount(current Coordinate, asteroidMap AsteroidMap) int {
	visibleCnt := 0
	for asteroidCoords, _ := range asteroidMap {
		if asteroidCoords == current {
			continue
		}

		if isDirectlySeen(current, asteroidCoords, asteroidMap) {
			visibleCnt += 1
		}
	}

	return visibleCnt
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
