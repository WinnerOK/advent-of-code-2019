package main

import (
	"io/ioutil"
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
	return strings.Split(s, "-")
}
