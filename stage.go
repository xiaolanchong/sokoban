
package main

type Entities map[Coord]bool

type IStage interface {
	Walls() Entities
	Crates() Entities
	Slots() Entities
	Keeper() Entities
}

type stage struct {
	warehouse Warehouse
	arrangement	Arrangement
	slots	Entities
}
