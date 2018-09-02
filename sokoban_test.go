
package main

import (
	"testing"
	"strings"
	"fmt"
)

func TestStageBuilder(t *testing.T) {
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

func TestStageTextBuilding(t *testing.T) {
	stageStr := 
	"#####" + "\n" +
	"#@ X#" + "\n" +
	"# O #" + "\n" +
	"#####"
	
	reader := strings.NewReader(stageStr)
	stage, err := BuildTextStage(reader)
	if err != nil {
		t.Errorf("No error expected on building finish, error: %v", err)
	}
	
	if val, exists := stage.warehouse.walls[Coord{4, 0}]; !exists || !val {
		t.Errorf("%v, %v", val, exists)
	}
	expectedBox := Box{0, 0, 4, 3}
	if stage.warehouse.bound != expectedBox {
		t.Errorf("Crates wrong bound: %v. %v expected", stage.warehouse.bound, expectedBox)
	}
	
	if stage.arrangement.keeper.x != 1 || stage.arrangement.keeper.y != 1  {
		t.Errorf("Keeper is placed in %v", stage.arrangement.keeper)
	}
	if len(stage.slots) != 1 || !stage.slots[Coord{3, 1}]  {
		t.Errorf("Slots are misplaced: %v", stage.slots)
	}
	if len(stage.arrangement.crates) != 1 || !stage.arrangement.crates[Coord{2, 2}]  {
		t.Errorf("Crates are misplaced: %v", stage.arrangement.crates)
	}

}

func TestArrangementHash(t *testing.T) {
	a := Arrangement{ make(Entities), Coord{0, 0} }
	crc := a.GetHash()
	md5 := a.GetMd5Hash()
	if crc != 0x41d912ff || md5 == "" {
		t.Errorf("No hash calculated: %v, %v", crc, md5)
	}
	t.Logf("Hash: crc32=%x, md5=%x", crc, md5)
}

func TestMovementInVacantPlace(t *testing.T) {
	bound := Box{0, 0, 5, 5}
	warehouse := Warehouse{ make(Entities), bound }
	
	a := Arrangement{ make(Entities), Coord{1, 1} }
	nextArr, _ := a.Move(warehouse, DirLeft)
	expected := Coord{0, 1}
	if expected != nextArr.keeper {
		t.Errorf("Failed to move left: %v", nextArr)
	}
	
	nextArr, _ = a.Move(warehouse, DirRight)
	expected = Coord{2, 1}
	if expected != nextArr.keeper {
		t.Errorf("Failed to move left: %v", nextArr)
	}

	nextArr, _ = a.Move(warehouse, DirUp)
	expected = Coord{1, 0}
	if expected != nextArr.keeper {
		t.Errorf("Failed to move up: %v", nextArr)
	}

	nextArr, _ = a.Move(warehouse, DirDown)
	expected = Coord{1, 2}
	if expected != nextArr.keeper {
		t.Errorf("Failed to move down: %v", nextArr)
	}
	
	a.keeper = Coord{0, 0}
	if _, err := a.Move(warehouse, DirLeft); err == nil {
		t.Errorf("Failed to restrict moving left")
	}
	a.keeper = Coord{0, 0}
	if _, err := a.Move(warehouse, DirUp); err == nil {
		t.Errorf("Failed to restrict moving up")
	}
	a.keeper = Coord{bound.right, bound.bottom}
	if _, err := a.Move(warehouse, DirDown); err == nil {
		t.Errorf("Failed to restrict moving down, keeper=%v, bound=%v", a.keeper, warehouse.bound)
	}
	a.keeper = Coord{bound.right, bound.bottom}
	if _, err := a.Move(warehouse, DirRight); err == nil {
		t.Errorf("Failed to restrict moving right, keeper=%v, bound=%v", a.keeper, warehouse.bound)
	}
}

func TestBumpIntoWalls(t *testing.T) {
	bound := Box{0, 0, 5, 5}
	walls := Entities{Coord{2, 2}: true, Coord{3, 2}: true}
	warehouse := Warehouse{ walls, bound }
	a := Arrangement{ make(Entities), Coord{2, 1} }

	if _, err := a.Move(warehouse, DirDown); err == nil {
		t.Errorf("Failed to restrict moving down")
	}
}

func TestMoveCrates(t *testing.T) {
	bound := Box{0, 0, 5, 5}
	walls := Entities{}
	warehouse := Warehouse{ walls, bound }
	
	crates := Entities{Coord{3, 2}: true, Coord{4, 2}: true}
	a := Arrangement{ crates, Coord{3, 1} }
	nextArr, err := a.Move(warehouse, DirDown)
	if err != nil {
		t.Errorf("Failed to restrict moving down: '%v'", err)
	}
	expected := Coord{3, 2}
	if nextArr.keeper != expected {
		t.Errorf("Keeper is in wrong position, %v", nextArr.keeper)
	}
	if _, exists := nextArr.crates[Coord{3, 3}]; !exists {
		t.Errorf("Crate not moved to %v", Coord{3, 3})
	}
}

func TestStuckState(t *testing.T) {
	bound := Box{0, 0, 15, 15}
	walls := Entities{Coord{3, 3}: true, Coord{4, 3}: true}
	warehouse := Warehouse{walls, bound}
	
	crates := Entities{Coord{3, 2}: true, Coord{4, 2}: true}
	arrangement := Arrangement{crates, Coord{0, 0}}
	if isTerminal := IsStuck(warehouse, arrangement); !isTerminal {
		t.Errorf("State is not terminal")	
	}
	
	crates = Entities{Coord{3, 2}: true}
	arrangement = Arrangement{crates, Coord{0, 0}}
	if isTerminal := IsStuck(warehouse, arrangement); isTerminal {
		t.Errorf("State is terminal")
	}
}

func TestIsSolved(t *testing.T) {
	crates := Entities{Coord{3, 2}: true, Coord{4, 2}: true}
	arrangement := Arrangement{crates, Coord{0, 0}}
	slots := Entities{Coord{3, 2}: true, Coord{4, 2}: true}
	
	if solved, _ := IsSolved(arrangement, slots); !solved {
		t.Errorf("Stage not solved")
	}
	
	slots = Entities{Coord{3, 2}: true, Coord{4, 3}: true}
	if solved, _ := IsSolved(arrangement, slots); solved {
		t.Errorf("Stage solved")
	}
}

type DummyFeedback struct {}

func (f DummyFeedback) Render(arr Arrangement) {}

func textToSolution(stageStr string) string {
	reader := strings.NewReader(stageStr)
	stage, err := BuildTextStage(reader)
	
	solution, err := Solve(stage.warehouse, stage.arrangement, stage.slots, DummyFeedback{})
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return PathToString(solution)
}

func TestSolution5x5_1Crate(t *testing.T) {
	stageStr := 
	"     " + "\n" +
	" @ X " + "\n" +
	"  O  " + "\n" +
	"     "
	
	reader := strings.NewReader(stageStr)
	stage, err := BuildTextStage(reader)
	
	solution, err := Solve(stage.warehouse, stage.arrangement, stage.slots, DummyFeedback{})
	if err != nil {
		t.Errorf("Error occured: %v", err)
		return
	}
	path := PathToString(solution)
	t.Logf("%v", path)
	if path != "_>_>^" {
		t.Errorf("Wrong path found: %v", path)
	}
}

func TestSolution_5x5_2Crates(t *testing.T) {
	return
	bound := Box{0, 0, 5, 5}
	walls := Entities{}
	warehouse := Warehouse{ walls, bound }
	
	crates := Entities{Coord{3, 2}: true, Coord{4, 2}: true}
	keeper := Coord{0, 0}
	arrangement := Arrangement{crates, keeper}

	slots  := Entities{Coord{3, 3}: true, Coord{4, 3}: true}
	
	solution, err := Solve(warehouse, arrangement, slots, DummyFeedback{})
	if err != nil {
		t.Errorf("Error occured: %v", err)
		return
	}
	path := PathToString(solution)
	t.Logf("%v", path)
	if path != ">>>__^>_" {
		t.Errorf("Wrong path found: %v", path)
	}
}

func TestSolution5x5_1Crate_Walls(t *testing.T) {
	stageStr := 
`   # 
 @ X 
 #O  
     `
	
	path := textToSolution(stageStr)
	if path != "<__>>^_<<^^>>" {
		t.Errorf("Wrong path found: %v", path)
	}
}

func TestSolution5x5_2CrateWalls(t *testing.T) {
	return
	stageStr := 
`   # 
 @ X 
 #O  
X  O `
	
	path := textToSolution(stageStr)
	if path != ">>>__<<<>^>^<^<<_>>" {
		t.Errorf("Wrong path found: %v", path)
	}
}

func TestSolution_Level1(t *testing.T) {
	return
	stageStr := 
`#########
#  @    #
#  OXXO #
#       #
#########`
	
	path := textToSolution(stageStr)
	if path != "<_>^>>>>_<" {
		t.Errorf("Wrong path found: %v", path)
	}
}

func TestSolution_Level2(t *testing.T) {
	return
	stageStr := 
`    #####
    #   #
    #O  #
  ###  O##
  #  O O #
### # ## #   ######
#   # ## #####  XX#
# O  O          XX#
##### ### #@##  XX#
    #     #########
    #######`
	
	path := textToSolution(stageStr)
	if path != "" {
		t.Errorf("Wrong path found: %v", path)
	}
}
