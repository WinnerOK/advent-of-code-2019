package main

import (
	"fmt"
	"math"
)

// Sizes
const (
	WIDTH  = 25
	HEIGHT = 6
)

// Colors
const (
	BLACK       = '0'
	WHITE       = '1'
	TRANSPARENT = '2'
)

type line = [WIDTH]byte
type layer = [HEIGHT]line

func splitLines(s string) []line {
	var result []line
	for i := 0; i < len(s)/WIDTH; i++ {
		var tmp line
		copy(tmp[:], s[i*WIDTH:(i+1)*WIDTH])
		result = append(result, tmp)
	}
	return result
}

func splitLayers(lines []line) []layer {
	var result []layer
	for i := 0; i < len(lines)/HEIGHT; i++ {
		var tmp layer
		copy(tmp[:], lines[i*HEIGHT:(i+1)*HEIGHT])
		result = append(result, tmp)
	}
	return result
}

func countByte(l layer, char byte) int {
	count := 0
	for _, i := range l {
		for _, j := range i {
			if j == char {
				count += 1
			}
		}
	}

	return count
}

func printLayer(l layer) {
	for _, i := range l {
		for _, j := range i {
			//print(j-'0', "")
			switch j {
			case WHITE:
				print("█")
			case BLACK:
				print(" ")
			}
		}
		println()
	}
}

func createEmptyLayer() layer {
	result := layer{}
	for rIDX, row := range result{
		for cIDX, _ := range row {
			result[rIDX][cIDX] = TRANSPARENT
		}
	}
	return result
}

func decodeImage(data []layer) layer {
	result := createEmptyLayer()
	for _, lay := range data{
		for rIDX, row := range lay{
			for cIDX, clr := range row{
				if result[rIDX][cIDX] == TRANSPARENT &&
					clr != TRANSPARENT {
					result[rIDX][cIDX] = clr
				}
			}
		}
	}
	return result
}

func main() {
	input := readInput("./in.txt")
	layers := splitLayers(splitLines(input))

	minZeroCnt := math.MaxInt64
	minZeroIDX := -1

	for i, lay := range layers {
		if cnt := countByte(lay, '0'); cnt < minZeroCnt {
			minZeroCnt = cnt
			minZeroIDX = i
		}
	}

	fmt.Printf("Answer for the part 1: %d\n",
		countByte(layers[minZeroIDX], '1')*countByte(layers[minZeroIDX], '2'))

	fmt.Printf("Answer for the part 2:\n")
	printLayer(decodeImage(layers))
}
