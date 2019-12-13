package main

import (
	"fmt"
)

const (
	Empty  = 0
	Wall   = 1
	Block  = 2
	Paddle = 3
	Ball   = 4
)

func blocksCnt(source []int) int {
	blocks := 0
	state := CreateState(source, []int{}, 3)
	var output []bigInt
	for state.PC.running {
		output, state = SimulateMachine(state)
		if len(output) == 3 {
			ent := int(output[2].Int64())
			if ent == Block {
				blocks += 1
			}

		}
	}
	return blocks
}

func noPlay(source []int) int {
	state := CreateState(source, []int{0}, 3)
	lastScore := 0
	var output []bigInt
	for state.PC.running {
		output, state = SimulateMachine(state)
		if len(output) == 3 {
			lastScore = int(output[2].Int64())
			state.addInput([]int{0})
		} else {
			return lastScore
		}
	}
	return 0 // will never reach this, if works properly
}

func main() {
	fileInput := readInput("./in.txt")
	source := stringSliceToIntSlice(fileInput)
	blocks := blocksCnt(source)
	fmt.Printf("Part 1 answer: %d\n", blocks)

	// Some cheating strategies
	// Find register corresponding to the paddle,
	//and fill this row with walls. Then make no moves. No matter, where the ball will fall, it will
	// bounce off the bottom wall
	var ballPos int
	for i := 0; i < len(source)-2; i++ {
		if source[i] == Empty && source[i+1] == Paddle && source[i+2] == Empty {
			ballPos = i + 1
			break
		}
	}

	for i := ballPos - 17; i <= ballPos+17; i++ {
		source[i] = Wall
	}
	source[0] = 2
	fmt.Printf("Part 2 answer: %d\n", noPlay(source))
}
