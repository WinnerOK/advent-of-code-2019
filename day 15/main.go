package main

import (
	"container/list"
	"fmt"
)

const (
	NORTH = 1
	SOUTH = 2
	WEST  = 3
	EAST  = 4
)

const (
	WallHit     = 0
	OK          = 1
	OxygenFound = 2
)

const (
	Empty  = ' '
	Wall   = '#'
	Oxygen = 'O'
)

type StepDescriptor struct {
	state         ProgramState
	stepDirection int
	exploringPos  Coordinate
	stepHistory   []int
}

var directions = []int{NORTH, SOUTH, WEST, EAST}

func findOxygen(memory []int) []int {
	//board := make([][]rune, 100)
	//for i := range board {
	//	board[i] = make([]rune, 100)
	//	for j := range board[i] {
	//		board[i][j] = Empty
	//	}
	//}

	isOxygenFound := false
	queue := list.New()
	state := CreateState(memory, []int{}, 1)
	for _, dir := range directions {
		//exploring := Coordinate{startX, startY}
		//exploring.move(dir)
		queue.PushBack(
			StepDescriptor{
				state:         CopyState(state),
				stepDirection: dir,
				stepHistory:   []int{dir},
				//exploringPos:  exploring,
			},
		)
	}

	var goodHistory []int

	for !isOxygenFound && queue.Len() > 0 {
		head := queue.Front()
		step := head.Value.(StepDescriptor)
		step.state.addInput([]int{step.stepDirection})
		output, newState := SimulateMachine(step.state)
		result := int(output[0].Int64())
		switch result {
		case OxygenFound:
			isOxygenFound = true
			goodHistory = step.stepHistory
		case OK:
			for _, dir := range directions {
				if lastStep := step.stepDirection;
					lastStep == SOUTH && dir == NORTH ||
						lastStep == NORTH && dir == SOUTH ||
						lastStep == WEST && dir == EAST ||
						lastStep == EAST && dir == WEST {
					// prohibit stepping back
					continue
				} else {
					//exploring := Coordinate{step.exploringPos.X,step.exploringPos.Y}
					//exploring.move(dir)
					tmp := make([]int, len(step.stepHistory))
					copy(tmp, step.stepHistory)
					tmp = append(tmp, dir)
					queue.PushBack(
						StepDescriptor{
							state:         CopyState(newState),
							stepDirection: dir,
							stepHistory:   tmp,
						},
					)
				}
			}

		}
		queue.Remove(head)
	}

	if !isOxygenFound {
		panic("Search has no result!")
	}

	return goodHistory
}

func main() {
	fileInput := readInput("./in.txt")
	source := stringSliceToIntSlice(fileInput)
	fmt.Printf("%v\n", findOxygen(source))
	fmt.Printf("Part 1 answer: %d\n", len(findOxygen(source)))
}

