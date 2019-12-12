package main

import "fmt"

const Steps = 100

func main() {
	fileInput := readInput("./in.txt")
	var moons []Moon
	for _, str := range fileInput {
		moons = append(moons, MakeMoon(parseCoordinate(str)))
	}
	initialState := make([]Moon, len(moons))
	copy(initialState, moons)
	system := makeMoonSystem(moons)
	for i := 0; i < Steps; i++ {
		system.simulateStep()
	}
	println(system.String())

	//	---------------------------------------------------------------------------------------------------
	//	Part 2
	// each state has a unique parent, since can be obtained by subtracting current velocities from current position
	// each axis exists independently

	// We must find first full repetition of each axis independently and then find LCM

	repeatedAfter := []int{-1, -1, -1}
	axisRepeated := []bool{false, false, false}

	for !all(axisRepeated...) {
		axisVels := system.axisVel()
		for i, alreadyRepeated := range axisRepeated {
			if !alreadyRepeated && axisVels[i] == 0 {
				repeated := true
				for j := range initialState {
					if initialState[j].coordinates.coordsSlice()[i] != system.moons[j].coordinates.coordsSlice()[i] {
						repeated = false
						break
					}
				}
				if repeated {
					axisRepeated[i] = true
					repeatedAfter[i] = system.steps + 1
				}
			}
		}
		system.simulateStep()
	}

	// To be sure, that answer for part 2 is correct, disable part 1 (set Steps constants to 0)
	fmt.Printf("Answer for part 2: %d\n", LCM(repeatedAfter[0], LCM(repeatedAfter[1], repeatedAfter[2])))

}
