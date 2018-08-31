package main

import (
	"fmt"
)

type SolutionNode struct {
	parent *SolutionNode
	left, right, up, down *SolutionNode
}

type Solution = []Direction

func goUpTree(start *SolutionNode) Solution {
	pathToRoot := make(Solution, 0, 30)

	for ;start != nil; start = start.parent {
		if start == nil || start.parent == nil {
			break
		}
		parent := start.parent
		var dir Direction
		switch start {
			case parent.left:  dir = DirLeft
			case parent.right: dir = DirRight
			case parent.up:    dir = DirUp
			case parent.down:  dir = DirDown
			default: panic("Solution node built incorrectly")
		}
		pathToRoot = append(pathToRoot, dir)
	}
	for i := len(pathToRoot)/2-1; i >= 0; i-- {
		opp := len(pathToRoot)-1-i
		pathToRoot[i], pathToRoot[opp] = pathToRoot[opp], pathToRoot[i]
	}
	return pathToRoot
}

type VisitedArrangements = map[string]uint32
type Solutions = []Solution

func solveRecursively(node *SolutionNode, warehouse Warehouse, arr Arrangement, 
					slots Entities, visited VisitedArrangements, allSolutions *Solutions, stepNumber uint32) error {
	//fmt.Printf("Path: %v\n", goUpTree(node))

	hash := arr.GetMd5Hash()
	if value, exists := visited[hash]; exists {
	//	fmt.Printf("Already visited\n")
		if stepNumber < value {
			visited[hash] = stepNumber
		} else {
			return fmt.Errorf("Already visited in fewer steps")
		}
	} else {
		visited[hash] = stepNumber
	}

	solved, err := IsSolved(arr, slots)
	if err != nil {
		return err
	}
	if solved {
		solution := goUpTree(node)
		newSolution := append(*allSolutions, solution)
		//copy(*allSolutions, newSolution)
		*allSolutions = newSolution
		//fmt.Printf("Solution found, len: %v\n", len(solution))
		return fmt.Errorf("Halted because a solution found")
	}
	if IsStuck(warehouse, arr) {
		return fmt.Errorf("Stuck")
	}
	
	var nextArr Arrangement
	var errArr error
	var child *SolutionNode
	stepNumber += 1
	
	child = &SolutionNode{node, nil, nil, nil, nil}
	node.left = child
	nextArr, errArr = arr.Move(warehouse, DirLeft)
	if errArr == nil {
		if err := solveRecursively(child, warehouse, nextArr, slots, visited, allSolutions, stepNumber); err == nil {
			//return solution, nil
		}
	}
	
	child = &SolutionNode{node, nil, nil, nil, nil}
	node.right = child
	nextArr, errArr = arr.Move(warehouse, DirRight)
	if errArr == nil {
		if err := solveRecursively(child, warehouse, nextArr, slots, visited, allSolutions, stepNumber); err == nil {
			//return solution, nil
		}
	}
	
	child = &SolutionNode{node, nil, nil, nil, nil}
	node.up = child
	nextArr, errArr = arr.Move(warehouse, DirUp)
	if errArr == nil {
		if err := solveRecursively(child, warehouse, nextArr, slots, visited, allSolutions, stepNumber); err == nil {
			//return solution, nil
		}
	}
	
	child = &SolutionNode{node, nil, nil, nil, nil}
	node.down = child
	nextArr, errArr = arr.Move(warehouse, DirDown)
	if errArr == nil {
		if err := solveRecursively(child, warehouse, nextArr, slots, visited, allSolutions, stepNumber); err == nil {
			//return solution, nil
		}
	}
	
	return fmt.Errorf("All branches traversed")
}

func Solve(warehouse Warehouse, arr Arrangement, slots Entities) (Solution, error) {
	visited := make(VisitedArrangements)
	root := &SolutionNode{}
	allSolutions := Solutions{}
	solveRecursively(root, warehouse, arr, slots, visited, &allSolutions, 0)
	if len(allSolutions) != 0 {
		//fmt.Printf("%v", allSolutions)
		fmt.Printf("Solution number=%v\n", len(allSolutions))
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

func SolutionToString(solution Solution) string {
	path := ""
	for _, v := range solution {
		var dirStr string
		switch v {
			case DirLeft:  dirStr = "<"
			case DirRight: dirStr = ">"
			case DirUp:    dirStr = "^"
			case DirDown:  dirStr = "_"
			default:       dirStr = "X"
		}
		path = path + dirStr
	}
	return path
}
