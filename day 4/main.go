package main

import (
	"fmt"
	"strconv"
)

func isPasswordValid(pass string) bool {
	return hasLengthValid(pass) && hasDoubleDigit(pass) && hasNonDecreasingDigits(pass)
}

func hasNonDecreasingDigits(pass string) bool {
	for idx := 0; idx < len(pass)-1; idx++ {
		if int(pass[idx])-int('0') > int(pass[idx+1])-int('0') {
			return false
		}
	}
	return true
}

func hasDoubleDigit(pass string) bool {
	for idx := 0; idx < len(pass)-1; idx++ {
		if pass[idx] == pass[idx+1] &&
			(idx+2 >= len(pass) || pass[idx] != pass[idx+2]) &&
			(idx-1 < 0 || pass[idx-1] != pass[idx]) {
			return true
		}
	}
	return false
}

func hasLengthValid(pass string) bool {
	return len(pass) == 6
}

func main() {
	input := readInput("./in.txt")
	start, _ := strconv.Atoi(input[0])
	end, _ := strconv.Atoi(input[1])

	var passwords []int

	for i := start; i <= end; i++ {
		if isPasswordValid(strconv.FormatInt(int64(i), 10)) {
			passwords = append(passwords, i)
		}
	}

	fmt.Printf("Answer for the part 1: %d\n", len(passwords))
	// For answer fo the part 2, remove additional conditions on lines 24-25

}
