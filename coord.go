
package main

import (
	"math"
)

type CoordUnit = int8

const (
	MinCoordUnit = math.MinInt8
	MaxCoordUnit = math.MaxInt8
)

type Coord struct {
	x, y CoordUnit
}
