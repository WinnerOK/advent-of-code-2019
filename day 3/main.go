package main

import (
	"fmt"
	"github.com/paulsmith/gogeos/geos"
	"math"
	"sort"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

func pointSum(p1 point, x, y int) point {
	return point{
		x: p1.x + x,
		y: p1.y + y,
	}
}

type lineSegment struct {
	begin point
	end   point
}

func printToGeos(s lineSegment) string {
	return fmt.Sprintf("LINESTRING(%d %d,%d %d)", s.begin.x, s.begin.y, s.end.x, s.end.y)
}

func segmentsIntersection(s1, s2 lineSegment) *geos.Geometry {
	line1 := geos.Must(geos.FromWKT(printToGeos(s1)))
	line2 := geos.Must(geos.FromWKT(printToGeos(s2)))
	return geos.Must(line1.Intersection(line2))
}

func findAllIntersections(layout1, layout2 []lineSegment) []point {
	var result []point
	for _, s1 := range layout1 {
		for _, s2 := range layout2 {
			intersection := segmentsIntersection(s1, s2)
			if b, _ := intersection.IsEmpty(); !b {
				x, _ := intersection.X()
				y, _ := intersection.Y()
				result = append(result, point{
					x: int(x),
					y: int(y),
				})
			}
		}
	}

	return result
}

func splitWireToSegments(wire []string) []lineSegment {
	var wireSegments []lineSegment
	current := point{
		x: 0,
		y: 0,
	}
	dx := map[string]int{"U": 0, "D": 0, "R": 1, "L": -1}
	dy := map[string]int{"U": 1, "D": -1, "R": 0, "L": 0}

	for _, step := range wire {
		num, _ := strconv.Atoi(step[1:])
		direction := step[:1]
		newPoint := pointSum(current, dx[direction]*num, dy[direction]*num)

		wireSegments = append(wireSegments, lineSegment{
			begin: current,
			end:   newPoint,
		})
		current = newPoint
	}

	return wireSegments
}

func makeWireDistMap(wirePlan []string) map[point]int {
	distMap := map[point]int{}
	current := point{
		x: 0,
		y: 0,
	}
	dx := map[string]int{"U": 0, "D": 0, "R": 1, "L": -1}
	dy := map[string]int{"U": 1, "D": -1, "R": 0, "L": 0}

	stepsDone := 0

	for _, move := range wirePlan {
		direction := move[:1]
		step, _ := strconv.Atoi(move[1:])

		for i := 0; i < step; i++ {
			current = pointSum(current, dx[direction], dy[direction])
			stepsDone += 1

			if _, ok := distMap[current]; !ok {
				distMap[current] = stepsDone
			}
		}

	}

	return distMap
}

func CalculateMinTotalLengthIntersection(wire1, wire2 map[point]int, intersections []point) int {
	min := math.MaxInt64
	for _, v := range intersections {
		dist := wire1[v] + wire2[v]
		if dist < min {
			min = dist
		}
	}
	return min
}

func main() {
	input := readInput("./in.txt")
	wire1DistMap := makeWireDistMap(strings.Split(input[0], ","))
	wire2DistMap := makeWireDistMap(strings.Split(input[1], ","))

	wire1Segments := splitWireToSegments(strings.Split(input[0], ","))
	wire2Segments := splitWireToSegments(strings.Split(input[1], ","))

	intersections := findAllIntersections(wire1Segments, wire2Segments)
	sort.Slice(intersections, func(i, j int) bool {
		return math.Abs(float64(intersections[i].x))+math.Abs(float64(intersections[i].y)) <
			math.Abs(float64(intersections[j].x))+math.Abs(float64(intersections[j].y))
	})
	fmt.Printf("Answer for part 1: %d\n", intersections[1].y+intersections[1].x)

	fmt.Printf("Answer for part 2: %d\n",
		CalculateMinTotalLengthIntersection(wire1DistMap, wire2DistMap, intersections[1:]))

}
