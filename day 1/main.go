package main

import (
	"fmt"
)

func fuelCalculate(mass int) int {
	return mass/3.0 - 2
}

func fuelCalculateAll(mass int) int {
	needFuel := fuelCalculate(mass)
	if needFuel > 0 {
		return needFuel + fuelCalculateAll(needFuel)
	} else {
		return 0
	}
}

func part1(massArr []int) {
	var total = 0
	for _, mass := range massArr {
		total += fuelCalculate(mass)
	}
	fmt.Println("Part 1 answer: ", total)
}

func part2(massArr []int) {
	var total = 0
	for _, mass := range massArr {
		fuelForMass := fuelCalculate(mass)
		total += fuelForMass + fuelCalculateAll(fuelForMass)
	}
	fmt.Println("Part 2 answer: ", total)
}

func main() {
	input := readInput("./in.txt")
	intInput := stringSliceToIntSlice(input)
	part1(intInput)
	part2(intInput)
}
