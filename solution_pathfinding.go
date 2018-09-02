
package main

import (
	"fmt"
)

type SolutionNodePF struct {
	parent		*SolutionNodePF
	pathFromParent []Direction
	children	[]*SolutionNodePF
}

//type VisitedArrangements = map[string]uint32

func goUpTreePF(start *SolutionNodePF) Solution {
	nodesToRoot := make([]*SolutionNodePF, 0, 30)
	result := make(Solution, 0, 100)

	for ;start != nil; start = start.parent {
		nodesToRoot = append(nodesToRoot, start)
	}
	for i := int32(len(nodesToRoot)) - 1 ; i >= 0; i-- {
		result = append(result, nodesToRoot[i].pathFromParent...)
	}
	return result
}

func applyMoves(pos Coord, path []Direction, bound Box) Coord {
	for _, move := range path {
		pos, _ = moveWithinBox(pos, bound, move)
	}
	return pos
}

func solveWithPathfinding(arr Arrangement, warehouse Warehouse,
						slots Entities, thisNode *SolutionNodePF, feedback IFeedback,
						pathLength uint32,
						allSolutions *Solutions,
						visited VisitedArrangements) (error) {
	
	if feedback != nil {
		feedback.Render(arr)
	}
	
	hash := arr.GetMd5Hash()
	if value, exists := visited[hash]; exists {
		if pathLength < value {
			visited[hash] = pathLength
		} else {
			return fmt.Errorf("Already visited in fewer steps")
		}
	} else {
		visited[hash] = pathLength
	}
	
	solved, err := IsSolved(arr, slots)
	if err != nil {
		return err
	}
	if solved {
		solution := goUpTreePF(thisNode)
		newSolution := append(*allSolutions, solution)
		*allSolutions = newSolution
		return fmt.Errorf("Halted because a solution found")
	}
	if IsStuckExt(warehouse, arr, slots) {
		return fmt.Errorf("Stuck")
	}
	
	pushPoints := FindPushPoints(warehouse, arr)
	for _, pushPoint := range pushPoints {
		newKeeperPos := applyMoves(arr.keeper, pushPoint.path, warehouse.bound) // todo, get from FindPushPoints directly
		newArr := Arrangement{arr.crates, newKeeperPos}
		movedArr, err := newArr.Move(warehouse, pushPoint.direction)
		
		if err == nil {
			fullPath := append(pushPoint.path, pushPoint.direction)
			newNode := SolutionNodePF{thisNode, fullPath, make([]*SolutionNodePF, 0, 10)}
			thisNode.children = append(thisNode.children, &newNode)
			newPathLength := pathLength + uint32(len(fullPath))
			solveWithPathfinding(movedArr, warehouse, slots, &newNode, feedback, newPathLength, allSolutions, visited)
		}
	}
	
	return nil
}


func SolveWithPathfinding(warehouse Warehouse, arr Arrangement, 
						slots Entities, feedback IFeedback) (Solution, error) {
	visited := VisitedArrangements{}
	allSolutions := make(Solutions, 0, 10)
	root := &SolutionNodePF{nil, make([]Direction, 0, 0), make([]*SolutionNodePF, 0, 10)}
	solveWithPathfinding(arr, warehouse, slots, root, feedback, 0, &allSolutions, visited)
	
	if len(allSolutions) != 0 {
		//fmt.Printf("%v", allSolutions)
		//fmt.Printf("Solution number=%v\n", len(allSolutions))
		minIndex, minLen := 0, len(allSolutions[0])
		for i := 1; i < len(allSolutions); i += 1 {
			if minLen > len(allSolutions[i]) {
				minIndex, minLen = i, len(allSolutions[i])
			}
		}
		return allSolutions[minIndex], nil
	}
	
	return nil, fmt.Errorf("No solution found")
}