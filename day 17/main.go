package main

import "fmt"

func buildBoard(source []int) [][]string {
	board := make([][]string, 1)
	board[0] = make([]string, 0)
	line := 0
	state := CreateState(source, []int{}, 0)
	output, state := SimulateMachine(state)
	//for _, out := range output {
	//	print(string(out.Int64()))
	//}
	//println("\n\n")
	for _, out := range output {
		symbol := string(out.Int64())
		if symbol == "\n" {
			board = append(board, make([]string, 0))
			line++
		} else {
			board[line] = append(board[line], symbol)
		}
	}
	return board
}

func calculateAlignmentParameterSum(board [][]string) int {
	result := 0
	height := len(board)-2
	width := len(board[0])
	for i := 1; i < height-1; i++ {
		for j := 1; j < width-1; j++ {
			if board[i][j] == "#" && board[i+1][j] == "#" &&
				board[i-1][j] == "#" && board[i][j+1] == "#" &&
				board[i][j-1] == "#" {
				result += i*j
				board[i][j]="O"
			}
		}
	}

	for _, i := range board {
		for _, j := range i {
			print(j)
		}
		println()
	}

	return result
}

func main() {
	fileInput := readInput("./in.txt")
	source := stringSliceToIntSlice(fileInput)
	board := buildBoard(source)
	fmt.Printf("%d\n", calculateAlignmentParameterSum(board))
}
