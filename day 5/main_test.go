package main_test

import (
	"."
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func captureOutput(f func([]int, int)int, data []int, input int) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f(data, input)
	log.SetOutput(os.Stderr)
	return buf.String()
}

func TestSimulateMachineEqual8(t *testing.T) {
	source := []int{3,9,8,9,10,9,4,9,99,-1,8}
	output := captureOutput( main.SimulateMachine, source, 1)
	assert.Equal(t,"[Machine] 0", output)
}
