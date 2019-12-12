package main

import "fmt"

type MoonSystem struct {
	moons []Moon
	steps int
}

func makeMoonSystem(moons []Moon) MoonSystem {
	return MoonSystem{
		moons: moons,
		steps: 0,
	}
}

func (ms *MoonSystem) addMoons(moons []Moon) {
	ms.moons = append(ms.moons, moons...)
}

func (ms *MoonSystem) applyGravity() {
	for i := range ms.moons {
		for j := i; j < len(ms.moons); j++ {
			applyGravityMoon(&ms.moons[i], &ms.moons[j])
		}
	}
}

func (ms *MoonSystem) applyVelocity() {
	for i := range ms.moons{
		ms.moons[i].move()
	}
}

func (ms MoonSystem) totalEnergy() int {
	total := 0
	for _, moon := range ms.moons {
		total += moon.totalEnergy()
	}
	return total
}

func (ms *MoonSystem) simulateStep() {
	//ms.resetVelocities()
	ms.applyGravity()
	ms.applyVelocity()
	ms.steps += 1
}

func (ms *MoonSystem) resetVelocities()  {
	for i := range ms.moons{
		ms.moons[i].resetVelocity()
	}
}

func (ms MoonSystem) String() string{
	result := fmt.Sprintf("After step %d\n", ms.steps)
	for _, moon := range ms.moons {
		result += fmt.Sprintf("%s\n", moon.String())
	}
	result += fmt.Sprintf("Total energy: %d", ms.totalEnergy())
	return result
}
