package main

import (
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readInput(filename string) []string {
	data, err := ioutil.ReadFile(filename)
	check(err)
	s := string(data)
	return strings.Split(s, "\n")
}

func stringSliceToIntSlice(strs []string) []int {
	var nums []int
	for _, s := range strs {
		if len(s) > 0 {
			n, _ := strconv.Atoi(s)
			nums = append(nums, n)
		}
	}
	return nums
}

func IntAbs(num int) int {
	return int(math.Abs(float64(num)))
}

func all(bools ...bool) bool {
	for _, exp := range bools {
		if !exp {
			return false
		}
	}
	return true
}

func GCD(a, b int) int {

	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}

	return a
}

func LCM(a, b int) int {
	return (a * b) / GCD(a, b)
}
