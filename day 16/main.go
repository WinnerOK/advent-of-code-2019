package main

import "fmt"

const (
	Phases = 100
	Repetitions = 10000
)

var basePattern = []int{0, 1, 0, -1}

func fft(signal []int) []int {

	lastDigit := func(num int) int {
		if num < 0 {
			return (-num) % 10
		} else {
			return num % 10
		}
	}

	var transformedSignal []int
	for pos := 0; pos < len(signal); pos++ {
		cursor := 0
		stayLeft := pos
		partialSum := 0
		for partIDX := 0; partIDX < len(signal); partIDX++ {
			if stayLeft == 0 {
				stayLeft = pos + 1
				cursor = (cursor + 1) % len(basePattern)
				if basePattern[cursor] == 0 {
					partIDX += stayLeft
					cursor = (cursor + 1) % len(basePattern)
				}
			}
			if partIDX >= len(signal){
				break
			}
			var addition int
			if basePattern[cursor] == 1{
				addition = signal[partIDX]
			}
			if basePattern[cursor] == -1{
				addition = -signal[partIDX]
			}
			partialSum += addition
			stayLeft--
		}
		transformedSignal = append(transformedSignal, lastDigit(partialSum))
	}
	return transformedSignal
}

func repeatSignal(signal []int, mult int) []int {
	var res []int
	for i := 0; i < mult; i++ {
		res = append(res, signal...)
	}
	return res
}

func getOffset(signal []int, length int) int {

	pow10 := func(exp int) int {
		res := 1
		for i := 0; i < exp; i++ {
			res *= 10
		}
		return res
	}

	offset := 0
	for i := 1; i <= length; i++ {
		addition := signal[i-1] * pow10(length-i)
		offset += addition
	}
	return offset
}

func main() {
	fileInput := readInput("./in.txt")
	inputSignal := stringSliceToIntSlice(fileInput)
	signal := make([]int, len(inputSignal))
	copy(signal, inputSignal)
	for i := 0; i < Phases; i++ {
		signal = fft(signal)
	}
	fmt.Printf("part 1 answer: %v\n", signal[:8])
	offset := getOffset(inputSignal, 7)
	signal = repeatSignal(inputSignal, Repetitions)
	for i:=0; i<Phases; i++{
		signal = fft(signal)
		//if i%10 == 9 {
			fmt.Printf("Stage %d passed\n", i+1)
		//}
	}
	fmt.Printf("%v\n", signal[offset:offset+8])
}
