package main

import (
	"fmt"
)

type stageBuilder struct {
	walls	Entities
	crates	Entities
	slots	Entities
	keeper	*Coord

	box		Box
}

func NewBuilder () stageBuilder {
	return stageBuilder{ make(Entities), make(Entities), make(Entities), nil, Box{} }
}

func (sb stageBuilder) checkSpace(c Coord) error {
	_, exists := sb.walls[c]
	if (exists) {
		return fmt.Errorf("There is a wall in (%d, %d)", c.x, c.y)
	}

	_, exists = sb.crates[c]
	if (exists) {
		return fmt.Errorf("There is a crate in (%d, %d)", c.x, c.y)
	}

	_, exists = sb.slots[c]
	if (exists) {
		return fmt.Errorf("There is a slot in (%d, %d)", c.x, c.y)
	}
	
	return nil
}

func (sb *stageBuilder) adjustRect(c Coord) {
	if sb.box.left > c.x {
		sb.box.left = c.x
	}

	if sb.box.top > c.y {
		sb.box.top = c.y
	}

	inflated := Coord{c.x, c.y}
	if sb.box.right < inflated.x {
		sb.box.right = inflated.x
	}

	if sb.box.bottom < inflated.y {
		sb.box.bottom = inflated.y
	}
}

func (sb *stageBuilder) addEntity(e *Entities, c Coord) error {
	if e := sb.checkSpace(c); e != nil {
		return e
	}

	sb.adjustRect(c)
	(*e)[c] = true

	return nil
}

func (sb *stageBuilder) AddWall(c Coord) error {
	return sb.addEntity(&sb.walls, c)
}

func (sb *stageBuilder) AddCrate(c Coord) error {
	return sb.addEntity(&sb.crates, c)
}

func (sb *stageBuilder) AddSlot(c Coord) error {
	return sb.addEntity(&sb.slots, c)
}

func (sb *stageBuilder) AddSpace(c Coord) error {
	if e := sb.checkSpace(c); e != nil {
		return e
	}
	sb.adjustRect(c)
	return nil
}

func (sb *stageBuilder) SetKeeper(c Coord) error {
	sb.keeper = &c
	return nil
}

func (sb stageBuilder) Finish() (stage, error) {
	if sb.keeper == nil {
		return stage{}, fmt.Errorf("Keeper position not set")
	}
	
	//box := Box{sb.topLeft.x, sb.topLeft.y, sb.}
	warehouse := Warehouse{sb.walls, sb.box}
	arrangement := Arrangement{sb.crates, *sb.keeper}
	return stage{warehouse, arrangement, sb.slots}, nil
}
