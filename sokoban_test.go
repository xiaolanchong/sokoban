
package main

import (
	"testing"
	"strings"
	//"fmt"
)

func TestCoordSteps(t *testing.T) {
	origin := Coord{0, 0}

	if c2 := origin.Left(); c2.x != -1 || c2.y != 0 {
		t.Errorf("(%d, %d) != (-1, 0)", c2.x, c2.y)
	}

	if c2 := origin.Right(); c2.x != 1 || c2.y != 0 {
		t.Errorf("(%d, %d) != (1, 0)", c2.x, c2.y)
	}

	if c2 := origin.Up(); c2.x != 0 || c2.y != -1 {
		t.Errorf("(%d, %d) != (0, -1)", c2.x, c2.y)
	}

	if c2 := origin.Down(); c2.x != 0 || c2.y != 1 {
		t.Errorf("(%d, %d) != (0, 1)", c2.x, c2.y)
	}
}

func TestRoomBuilder(t *testing.T) {
	rb := NewBuilder()

	if e := rb.AddWall(Coord{0, 0}); e != nil {
		t.Errorf("%s", e)
	}

	if e := rb.AddWall(Coord{0, 0}); e == nil {
		t.Errorf("Error expected")
	}

	if e := rb.AddCrate(Coord{1, 1}); e != nil {
		t.Errorf("%s", e)
	}
	if e := rb.AddSlot(Coord{2, 2}); e != nil {
		t.Errorf("%s", e)
	}

	if _, exists := rb.slots[Coord{2, 2}]; !exists {
		t.Errorf("%v", exists)
	}
	if _, exists := rb.crates[Coord{2, 2}]; exists {
		t.Errorf("%v", exists)
	}
	
	rb.SetKeeper(Coord{0, 0})
	if _, err := rb.Finish(); err != nil {
		t.Errorf("No error expected on building finish, error: %v", err)
	}
}

func TestRoomTextBuilding(t *testing.T) {
	room_input := 
	"xxxxx" + "\n" +
	"xA .x" + "\n" +
	"x o x" + "\n" +
	"xxxxx"
	
	reader := strings.NewReader(room_input)
	room, err := BuildTextRoom(reader)
	if err != nil {
		t.Errorf("No error expected on building finish, error: %v", err)
	}
	
	if val, exists := room.walls[Coord{4, 0}]; !exists || !val {
		t.Errorf("%v, %v", val, exists)
	}
	
	if room.keeper.x != 1 || room.keeper.y != 1  {
		t.Errorf("Keeper is placed in %v", room.keeper)
	}
	if len(room.slots) != 1 || !room.slots[Coord{3, 1}]  {
		t.Errorf("Slots are misplaced: %v", room.slots)
	}
	if len(room.crates) != 1 || !room.crates[Coord{2, 2}]  {
		t.Errorf("Crates are misplaced: %v", room.crates)
	}
}