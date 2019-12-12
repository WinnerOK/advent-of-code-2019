package main

const Steps = 1000

func main() {
	fileInput := readInput("./in.txt")
	var moons []Moon
	for _, str := range fileInput{
		moons = append(moons, MakeMoon(parseCoordinate(str)))
	}

	system := makeMoonSystem(moons)
	for i:=0; i <=Steps; i++{
		println(system.String())
		system.simulateStep()
		println()
	}

}
