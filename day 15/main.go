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
	Empty  = " "
	Wall   = "#"
	Oxygen = "O"
)

type StepDescriptor struct {
	state         ProgramState
	stepDirection int
	exploringPos  Coordinate
	stepHistory   []int
}

var directions = []int{NORTH, SOUTH, WEST, EAST}

func printBoard(board [][]string) {
	for i := range board {
		for j := range board[i] {
			fmt.Printf("%s", board[i][j])
		}
		println()
	}
}

func exploreBoard(source []int, board [][]string, startX, startY int) ([]int, Coordinate) {
	used := make([][]bool, len(board))
	for i := range used {
		used[i] = make([]bool, len(board[i]))
	}
	boardGetUsed := func(coordinate Coordinate) bool { return used[coordinate.Y][coordinate.X] }
	boardSetUsed := func(coordinate Coordinate) { used[coordinate.Y][coordinate.X] = true }
	isOxygenFound := false
	queue := list.New()
	state := CreateState(source, []int{}, 1)
	for _, dir := range directions {
		exploring := Coordinate{startX, startY}
		exploring.move(dir)
		queue.PushBack(
			StepDescriptor{
				state:         CopyState(state),
				stepDirection: dir,
				stepHistory:   []int{dir},
				exploringPos:  exploring,
			},
		)
		boardSetUsed(exploring)
	}

	var goodHistory []int
	var OxygenPosition Coordinate

	for /*!isOxygenFound && */ queue.Len() > 0 {
		head := queue.Front()
		step := head.Value.(StepDescriptor)
		step.state.addInput([]int{step.stepDirection})
		output, newState := SimulateMachine(step.state)
		result := int(output[0].Int64())
		switch result {
		case OxygenFound:
			isOxygenFound = true
			goodHistory = step.stepHistory
			OxygenPosition = step.exploringPos
			board[step.exploringPos.Y][step.exploringPos.X] = Oxygen
		case WallHit:
			board[step.exploringPos.Y][step.exploringPos.X] = Wall
		case OK:
			for _, dir := range directions {
				exploring := Coordinate{step.exploringPos.X, step.exploringPos.Y}
				exploring.move(dir)
				if !boardGetUsed(exploring) {
					boardSetUsed(exploring)
					tmp := make([]int, len(step.stepHistory))
					copy(tmp, step.stepHistory)
					tmp = append(tmp, dir)
					queue.PushBack(
						StepDescriptor{
							state:         CopyState(newState),
							stepDirection: dir,
							stepHistory:   tmp,
							exploringPos:  exploring,
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
	return goodHistory, OxygenPosition
}

func bfs(board [][]string, start Coordinate) int {
	used := make([][]bool, len(board))
	for i := range used {
		used[i] = make([]bool, len(board[i]))
	}
	boardGet := func(coordinate Coordinate) string { return board[coordinate.Y][coordinate.X] }
	boardSet := func(coordinate Coordinate, c string) { board[coordinate.Y][coordinate.X] = c }
	boardGetUsed := func(coordinate Coordinate) bool { return used[coordinate.Y][coordinate.X] }
	boardSetUsed := func(coordinate Coordinate) { used[coordinate.Y][coordinate.X] = true }

	step := 0
	queue := list.New()
	queue.PushBack(start)
	boardSetUsed(start)
	stepsToFinishInterval := 1
	nextIntervalLength := 0
	for queue.Len() > 0 {
		head := queue.Front()
		pos := head.Value.(Coordinate)
		boardEnt := boardGet(pos)
		if boardEnt == Empty || boardEnt == Oxygen {
			for _, dir := range directions {
				exploring := Coordinate{pos.X, pos.Y}
				exploring.move(dir)
				exploringEnt := boardGet(exploring)
				if !boardGetUsed(exploring) && exploringEnt != Wall {
					boardSetUsed(exploring)
					boardSet(exploring, Oxygen)
					queue.PushBack(exploring)
					nextIntervalLength += 1
				}
			}
		}
		queue.Remove(head)
		stepsToFinishInterval -= 1
		if stepsToFinishInterval == 0 {
			step += 1
			stepsToFinishInterval = nextIntervalLength
			nextIntervalLength = 0
		}
	}
	return step - 1
}

func main() {
	fileInput := readInput("./in.txt")
	source := stringSliceToIntSlice(fileInput)
	board := make([][]string, 51)
	for i := range board {
		board[i] = make([]string, 51)
		for j := range board[i] {
			board[i][j] = Empty
		}
	}
	pathToOxygen, oxygenPos := exploreBoard(source, board, 25, 25)
	printBoard(board)
	fmt.Printf("Part 1 answer: %d\n", len(pathToOxygen))
	fmt.Printf("Part 2 answer: %d\n", bfs(board, oxygenPos))
}
