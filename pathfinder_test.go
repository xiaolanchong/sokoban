
package main

import (
	"strings"
	"testing"
	"fmt"
)

func TestPathFinder_Simple(t *testing.T) {
	stageStr := 
`     
 @ X 
  O  
     `
	
	reader := strings.NewReader(stageStr)
	stage, _ := BuildTextStage(reader)
	
	points := FindPushPoints(stage.warehouse, stage.arrangement)

	if len(points) != 4 {
		t.Errorf("Wrong number of points: %v", points)
		return
	}
	
	var pathStr, directionStr string
	pathStr, directionStr = PathToString(points[0].path), DirectionToString(points[0].direction)
	if pathStr != "_" || directionStr != ">" {
		t.Errorf("Point 0, wrong path or direction: %v, %v", pathStr, directionStr) 
	}
	
	pathStr, directionStr = PathToString(points[1].path), DirectionToString(points[1].direction)
	if pathStr != ">" || directionStr != "_" {
		t.Errorf("Point 1, wrong path or direction: %v, %v", pathStr, directionStr) 
	}
	
	pathStr, directionStr = PathToString(points[2].path), DirectionToString(points[2].direction)
	if pathStr != "__>" || directionStr != "^" {
		t.Errorf("Point 2, wrong path or direction: %v, %v", pathStr, directionStr) 
	}
	
	pathStr, directionStr = PathToString(points[3].path), DirectionToString(points[3].direction)
	if pathStr != ">>_" || directionStr != "<" {
		t.Errorf("Point 3, wrong path or direction: %v, %v", pathStr, directionStr) 
	}
}

func TestPathFinder_Walls(t *testing.T) {
	stageStr := 
` @  X
 ### 
##O  
     `

	reader := strings.NewReader(stageStr)
	stage, _ := BuildTextStage(reader)
	
	points := FindPushPoints(stage.warehouse, stage.arrangement)

	if len(points) != 2 {
		t.Errorf("Wrong number of points: %v", points)
		return
	}

	var pathStr, directionStr string
	pathStr, directionStr = PathToString(points[0].path), DirectionToString(points[0].direction)
	if pathStr != ">>>__<_<" || directionStr != "^" {
		t.Errorf("Point 0, wrong path or direction: %v, %v", pathStr, directionStr) 
	}
	
	pathStr, directionStr = PathToString(points[1].path), DirectionToString(points[1].direction)
	if pathStr != ">>>__<" || directionStr != "<" {
		t.Errorf("Point 1, wrong path or direction: %v, %v", pathStr, directionStr) 
	}
}

func TestPathfindingSolution_5x5_1Crate(t *testing.T) {
	stageStr := 
	"     " + "\n" +
	" @ X " + "\n" +
	"  O  " + "\n" +
	"     "
	
	reader := strings.NewReader(stageStr)
	stage, err := BuildTextStage(reader)
	
	solution, err := SolveWithPathfinding(stage.warehouse, stage.arrangement, stage.slots, DummyFeedback{})
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

func textToSolutionPF(stageStr string) string {
	reader := strings.NewReader(stageStr)
	stage, err := BuildTextStage(reader)
	
	solution, err := SolveWithPathfinding(stage.warehouse, stage.arrangement, stage.slots, DummyFeedback{})
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return PathToString(solution)
}

func TestSolutionPF_5x5_1Crate_Walls(t *testing.T) {
	stageStr := 
`   # 
 @ X 
 #O  
     `
	
	path := textToSolutionPF(stageStr)
	if path != "<__>>^_<<^^>>" {
		t.Errorf("Wrong path found: %v", path)
	}
}

func TestSolutionPF_5x5_2CrateWalls(t *testing.T) {
	stageStr := 
`   # 
 @ X 
 #O  
X  O `
	
	path := textToSolutionPF(stageStr)
	if path != ">>>__<<<>^>^<^<<_>>" {
		t.Errorf("Wrong path found: %v", path)
	}
}

func TestSolutionPF_Level2(t *testing.T) {
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
	
	path := textToSolutionPF(stageStr)
	if path != "" {
		t.Errorf("Wrong path found: %v", path)
	}
}