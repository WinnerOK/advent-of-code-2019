package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Reagent struct {
	name   string
	amount int
}

type RecipeDescriptor struct {
	outputAmt int
	inputs    []Reagent
}

type RecipeBook map[string]RecipeDescriptor

func makeReagent(str string) Reagent {
	s := strings.Split(strings.TrimSpace(str), " ")
	amount, _ := strconv.Atoi(s[0])
	name := s[1]
	return Reagent{
		name:   name,
		amount: amount,
	}
}

func parseRecipe(str string) (string, RecipeDescriptor) {
	split := strings.Split(str, " => ")
	outputReagent := makeReagent(split[1])
	inputsStr := strings.Split(split[0], ", ")
	var inputs []Reagent
	for _, input := range inputsStr {
		inputs = append(inputs, makeReagent(input))
	}
	return outputReagent.name,
		RecipeDescriptor{
			outputAmt: outputReagent.amount,
			inputs:    inputs,
		}
}

//CalculateOreUsage
//How many ORE is needed to make {targetAmount} of {target}
//using {RecipeBook} knowing {extraReagentsAvailable} are available.
//target > 0
func CalculateOreUsage(recipes RecipeBook, target string, targetAmount int, extraReagentsAvailable map[string]int) int {
	if target == "ORE" {
		return targetAmount
	} else if targetAmount <= getOrDefaultMap(extraReagentsAvailable, target, 0) {
		extraReagentsAvailable[target] -= targetAmount
		return 0
	}

	targetAmount -= getOrDefaultMap(extraReagentsAvailable, target, 0)
	extraReagentsAvailable[target] = 0
	oreUsed := 0
	recipeDesc := recipes[target]
	outputAmount := recipeDesc.outputAmt
	copies := int(math.Ceil(float64(targetAmount) / float64(outputAmount)))
	for _, reagent := range recipeDesc.inputs {
		inputAmount := reagent.amount * copies
		oreUsed += CalculateOreUsage(recipes, reagent.name, inputAmount, extraReagentsAvailable)
	}
	extraReagentsAvailable[target] += outputAmount*copies - targetAmount
	return oreUsed
}

//CalculateFuelProduced calculates how much fuel can be produced from {oreAvailable} using {recipes}
// if 1 fuel require {orePerFuel} ore
func CalculateFuelProduced(recipes RecipeBook, oreAvailable, orePerFuel int) int {
	// The answer is not simply division since additional fuel can be made from extra resources and some ore
	fuelMade := 0
	targetAmount := oreAvailable / orePerFuel
	extraReagentsAvailable := map[string]int{}
	for oreAvailable > 0 && targetAmount > 0 {
		extraReagentsAvailableTmp := mapCopy(extraReagentsAvailable)
		oreUsage := CalculateOreUsage(recipes, "FUEL", targetAmount, extraReagentsAvailableTmp)
		if oreUsage <= oreAvailable {
			fuelMade += targetAmount
			oreAvailable -= oreUsage
			extraReagentsAvailable = extraReagentsAvailableTmp
		} else {
			targetAmount /= 2 // something like binary search
		}
	}
	return fuelMade
}

func main() {
	recipes := RecipeBook{}
	fileInput := readInput("./in.txt")

	for _, s := range fileInput {
		output, desc := parseRecipe(s)
		recipes[output] = desc
	}

	orePerFuel := CalculateOreUsage(recipes, "FUEL", 1, map[string]int{})
	fmt.Printf("Answer for part1: %d\n", orePerFuel, )

	fmt.Printf("Answer for part2: %d\n", CalculateFuelProduced(recipes, 1000000000000, orePerFuel))
}
