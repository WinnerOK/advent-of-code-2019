package main

import (
	"fmt"
	"math"
)

const (
	BLACK = 0
	WHITE = 1
)

const (
	UP    = 0
	RIGHT = 1
	DOWN  = 2
	LEFT  = 3
)

func getOrDefaultMap(m map[string]int, key string, def int) int {
	if v, ok := m[key]; ok {
		return v
	} else {
		return def
	}
}

func paint(source []int, startClr int) map[string]int {
	painted := map[string]int{}
	position := Coordinate{0, 0, UP}
	state := CreateState(source, []int{startClr}, 2)
	var output []bigInt
	for state.PC.running {
		output, state = SimulateMachine(state)
		if len(output) == 2 {
			paintedTo := int(output[0].Int64())
			rotation := int(output[1].Int64())
			painted[position.str()] = paintedTo

			position.rotate(rotation)
			position.move(1)
			state.addInput([]int{getOrDefaultMap(painted, position.str(), BLACK)})
		}
	}
	return painted
}

func out(color int){
	switch color {
	case BLACK:
		fmt.Printf(" ")
	case WHITE:
		fmt.Printf("█")
	}
}

func draw(painted map[string]int) {
	xMin := math.MaxInt64
	yMin := math.MaxInt64
	xMax := math.MinInt64
	yMax := math.MinInt64
	for k, _ := range painted {
		coords := parseCoordinate(k)
		if coords.X < xMin {
			xMin = coords.X
		}
		if coords.X > xMax {
			xMax = coords.X
		}
		if coords.Y < yMin {
			yMin = coords.Y
		}
		if coords.Y > yMax {
			yMax = coords.Y
		}
	}
	width := int(math.Abs(float64(xMin)) + math.Abs(float64(xMax)) + 1)
	heigh := int(math.Abs(float64(yMin)) + math.Abs(float64(yMax)) + 1)
	canvas := make([][]int, heigh)
	for y := range canvas {
		canvas[y] = make([]int, width)
		for x := range canvas[y] {
			correctedCoord := Coordinate{
				X:         xMin + x,
				Y:         yMin + y,
				direction: 0,
			}
			canvas[y][x] = getOrDefaultMap(painted,
				correctedCoord.str(),
				BLACK)

			out(canvas[y][x])
		}
		println()
	}
}

func main() {
	fileInput := readInput("./in.txt")
	source := stringSliceToIntSlice(fileInput)
	painted := paint(source, BLACK)
	fmt.Printf("Part 1 answer: %d\n", len(painted))

	painted = paint(source, WHITE)
	draw(painted) // draws mirrored and upside-down
}
