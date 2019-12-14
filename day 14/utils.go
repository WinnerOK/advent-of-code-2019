package main

import (
	"io/ioutil"
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

func getOrDefaultMap(m map[string]int, key string, def int) int {
	if v, ok := m[key]; ok {
		return v
	} else {
		return def
	}
}

func mapCopy(m map[string]int) map[string]int {
	newMap := map[string]int{}
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
