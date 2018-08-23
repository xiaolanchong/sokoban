
package main

import (
	"fmt"
)

type Entities map[Coord]bool

type IRoom interface {
	Walls() Entities
	Crates() Entities
	Slots() Entities
	Keeper() Entities
}

type Disposition struct {
	crates	Entities
	slots	Entities
	keeper	Coord
}

type room struct {
	walls	Entities
	crates	Entities
	slots	Entities
	keeper	Coord
}

type roomBuilder struct {
	walls	Entities
	crates	Entities
	slots	Entities
	keeper	*Coord

	topLeft		Coord
	bottomRight	Coord
}

func NewBuilder () roomBuilder {
	return roomBuilder{ make(Entities), make(Entities), make(Entities), nil, Coord{0, 0}, Coord{0, 0} }
}

func (rb roomBuilder) checkSpace(c Coord) error {
	_, exists := rb.walls[c]
	if (exists) {
		return fmt.Errorf("There is a wall in (%d, %d)", c.x, c.y)
	}

	_, exists = rb.crates[c]
	if (exists) {
		return fmt.Errorf("There is a crate in (%d, %d)", c.x, c.y)
	}

	_, exists = rb.slots[c]
	if (exists) {
		return fmt.Errorf("There is a slot in (%d, %d)", c.x, c.y)
	}
	
	return nil
}

func (rb *roomBuilder) adjustRect(c Coord) {
	if rb.topLeft.x > c.x {
		rb.topLeft.x = c.x
	}

	if rb.topLeft.y > c.y {
		rb.topLeft.y = c.y
	}

	if rb.bottomRight.x < c.x {
		rb.bottomRight.x = c.x
	}

	if rb.bottomRight.y < c.y {
		rb.bottomRight.y = c.y
	}
}

func (rb *roomBuilder) addEntity(e *Entities, c Coord) error {
	if e := rb.checkSpace(c); e != nil {
		return e
	}

	rb.adjustRect(c)
	(*e)[c] = true

	return nil
}

func (rb *roomBuilder) AddWall(c Coord) error {
	return rb.addEntity(&rb.walls, c)
}

func (rb *roomBuilder) AddCrate(c Coord) error {
	return rb.addEntity(&rb.crates, c)
}

func (rb *roomBuilder) AddSlot(c Coord) error {
	return rb.addEntity(&rb.slots, c)
}

func (rb *roomBuilder) SetKeeper(c Coord) error {
	rb.keeper = &c
	return nil
}

func (rb roomBuilder) Finish() (room, error) {
	if rb.keeper == nil {
		return room{}, fmt.Errorf("Keeper position not set")
	}
	return room{rb.walls, rb.crates, rb.slots, *rb.keeper}, nil
}