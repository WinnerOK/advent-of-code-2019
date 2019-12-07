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

func pointSum(p1, p2 point) point {
	return point{
		x: p1.x + p2.x,
		y: p1.y + p2.y,
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

func parseWire(wire []string) []lineSegment {
	var result []lineSegment
	current := point{
		x: 0,
		y: 0,
	}
	var delta point

	for _, step := range wire {
		num, _ := strconv.Atoi(step[1:])
		switch step[0] {
		case 'U':
			delta = point{
				y: num,
			}
		case 'D':
			delta = point{
				y: -num,
			}
		case 'R':
			delta = point{
				x: num,
			}
		case 'L':
			delta = point{
				x: -num,
			}
		}
		result = append(result, lineSegment{
			begin: current,
			end:   pointSum(current, delta),
		})

		current = pointSum(current, delta)

	}

	return result
}

func findAllIntersections(layout1, layout2 []lineSegment) []point {
	var result []point
	for _, s1 := range layout1 {
		for _, s2 := range layout2 {
			intersection := segmentsIntersection(s1, s2)
			if b, _ := intersection.IsEmpty(); !b {
				//fmt.Printf("Intersection: %v %v\n", s1, s2)
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

func rank(intersections []point) {
	sort.Slice(intersections, func(i, j int) bool {
		return math.Abs(float64(intersections[i].x))+math.Abs(float64(intersections[i].y)) < math.Abs(float64(intersections[j].x))+math.Abs(float64(intersections[j].y))
	})

	fmt.Printf("%v\nAnswer for part 1: %d", intersections, intersections[1].y+intersections[1].x)
}

func part1(wire1, wire2 []lineSegment) {
	rank(findAllIntersections(wire1, wire2))
}

//All intersections
//[{0 0} {258 0} {433 0} {433 408} {639 408} {990 153} {738 408} {496 683} {731 683} {890 708} {890 754} {890 915} {890 1208} {890 1260} {939 1367} {1166 1367} {1327 1787} {1327 1848} {1327 1888} {1327 1942} {1938 1348} {1327 1961} {1938 1454} {1327 2147} {2242 1348} {2194 1454} {2400 1348} {2290 1469} {2290 1474} {1435 2340} {2290 1530} {2744 1188} {1435 2757} {2290 2018} {3272 1220} {2419 2151} {2371 2231} {3272 1352} {2065 2757} {1389 3807} {2862 2719} {2727 2855} {2727 2990} {3026 2719} {3031 2719} {2104 3700} {2119 3700} {2119 3779} {2478 3463} {2478 3499} {2119 3988} {2119 3996} {3476 2719} {2308 3999} {2313 3999} {2906 3463} {3240 3670} {3714 3889} {4093 3912} {4888 3326} {5581 3167} {5702 3178} {5702 3531}]

func main() {
	input := readInput("./in.txt")
	wire1 := parseWire(strings.Split(input[0], ","))
	wire2 := parseWire(strings.Split(input[1], ","))

	part1(wire1, wire2)

}
