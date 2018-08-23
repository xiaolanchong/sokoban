
package main

type CoordUnit = int8

type Coord struct {
	x, y CoordUnit
}

func (c Coord) Left() Coord {
	return Coord{c.x - 1, c.y}
}

func (c Coord) Right() Coord {
	return Coord{c.x + 1, c.y}
}

func (c Coord) Up() Coord {
	return Coord{c.x, c.y - 1}
}

func (c Coord) Down() Coord {
	return Coord{c.x, c.y + 1}
}
