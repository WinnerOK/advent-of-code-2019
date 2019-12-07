package main

import (
	"fmt"
	"strings"
)

const ROOT = "COM"
const YOU = "YOU"
const SANTA = "SAN"

type set = map[string]bool

func calculateDepth(obj string, mapping map[string]string) int {
	if obj == ROOT {
		return 0
	} else {
		return 1 + calculateDepth(mapping[obj], mapping)
	}
}

func symmetricDiff(s1, s2 set) set {
	result := make(set)
	for k, _ := range s1 {
		if _, ok := s2[k]; !ok {
			result[k] = true
		}
	}
	for k, _ := range s2 {
		if _, ok := s1[k]; !ok {
			result[k] = true
		}
	}
	return result
}

func main() {
	input := readInput("./in.txt")
	orbitMapping := map[string]string{} // orbitMapping[A]=B iff B is on orbit of A

	for _, line := range input {
		orbiter := strings.Split(line, ")")[0]
		orbitee := strings.Split(line, ")")[1]
		orbitMapping[orbitee] = orbiter
	}

	total := 0
	for k, _ := range orbitMapping {
		total += calculateDepth(k, orbitMapping)
	}

	fmt.Printf("Answer for part 1: %d\n", total)
	youObj := YOU
	santaObj := SANTA
	youWay := make(set)
	santaWay := make(set)

	for youObj != ROOT {
		youWay[youObj] = true
		youObj = orbitMapping[youObj]
	}

	for santaObj != ROOT {
		santaWay[santaObj] = true
		santaObj = orbitMapping[santaObj]
	}

	// subtract 2 since for n+2 nodes there are n transitions
	fmt.Printf("Answer for the part 2: %d\n", len(symmetricDiff(youWay, santaWay))-2)
}
