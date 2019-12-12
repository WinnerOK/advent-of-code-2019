package main

import "fmt"

type Moon struct {
	coordinates Coordinate
	velocity    Coordinate
}

func MakeMoon(coordinate Coordinate) Moon {
	return Moon{
		coordinates: coordinate,
		velocity:    Coordinate{0, 0, 0},
	}
}

func (m Moon) potentialEnergy() int {
	return IntAbs(m.coordinates.X) + IntAbs(m.coordinates.Y) + IntAbs(m.coordinates.Z)
}

func (m Moon) kineticEnergy() int {
	return IntAbs(m.velocity.X) + IntAbs(m.velocity.Y) + IntAbs(m.velocity.Z)
}

func (m Moon) totalEnergy() int {
	return m.potentialEnergy() * m.kineticEnergy()
}

func (m *Moon) move() {
	m.coordinates.addCoords(m.velocity)
}

func (m *Moon) resetVelocity(){
	m.velocity = Coordinate{0,0,0}
}

func (m Moon) String() string {
	return fmt.Sprintf("pos=<x=%2d, y=%2d, z=%2d> vel=<x=%2d, y=%2d, z=%2d>",
		m.coordinates.X, m.coordinates.Y, m.coordinates.Z,
		m.velocity.X, m.velocity.Y, m.velocity.Z)
}

func calculateChange(a, b int) int {
	switch {
	case a == b:
		return 0
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		panic("Unexpected behaviour")
	}
}

func applyGravityMoon(m1, m2 *Moon) {
	dx := calculateChange(m1.coordinates.X, m2.coordinates.X)
	dy := calculateChange(m1.coordinates.Y, m2.coordinates.Y)
	dz := calculateChange(m1.coordinates.Z, m2.coordinates.Z)

	m1.velocity.X -= dx
	m2.velocity.X += dx
	m1.velocity.Y -= dy
	m2.velocity.Y += dy
	m1.velocity.Z -= dz
	m2.velocity.Z += dz
}
