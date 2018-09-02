
package main

import (
	//"fmt"
	"sort"
)

type Path = []Direction

type PushPoint struct {
	path Path
	direction  Direction
}

type PushPoints = [] PushPoint

type SearchArea = map[Coord]Path


func getOppositeDir(direction Direction) Direction {
	switch direction {
		case DirLeft:  return DirRight
		case DirRight: return DirLeft
		case DirUp:    return DirDown
		case DirDown:  return DirUp
		default: panic("Unknown direction")
	}
}


type PointAndDir struct { 
	point Coord
	direction Direction
}

type ResultMap = map[PointAndDir] Path

func copyPath(src Path) Path {
	dst := make(Path, len(src))
	copy(dst, src)
	return dst
}

func findPushPointsRecursively(warehouse Warehouse, arr Arrangement, area SearchArea, path Path, keeperPos Coord, direction Direction, resultMap ResultMap) {
	nextPos, err := moveWithinWalls(keeperPos, warehouse, direction)
	if err != nil {
		return // can't move this way
	}
	
	if _, exists := arr.crates[nextPos]; exists { // bumped into a crate
		key := PointAndDir{keeperPos, direction}
		if value, pathExists := resultMap[key]; !pathExists || (len(path) < len(value)) {
			resultMap[key] = copyPath(path)
		}
		return
	}

	if existingPath, exists := area[nextPos]; !exists || ((len(path) + 1) < len(existingPath)) {
		newPath := make(Path, len(path) + 1)
		copy(newPath, path)
		newPath[len(newPath) - 1] = direction
		
		area[nextPos] = copyPath(newPath)
		oppositeDir := getOppositeDir(direction)
		for _, nextDir := range []Direction { DirLeft, DirRight, DirDown, DirUp} {
			if nextDir != oppositeDir {
				findPushPointsRecursively(warehouse, arr, area, newPath, nextPos, nextDir, resultMap)
			}
		}
	}
}

func stableSortResultMap(resultMap ResultMap) PushPoints {
	results := make(PushPoints, 0, len(resultMap))
	keys := make([]Coord, 0, len(resultMap))
	for k, _ := range resultMap {
		keys = append(keys, k.point)
	}
	sort.Slice(keys, func(i, j int) bool { 
		if keys[i].x < keys[j].x {
			return true
		} else if(keys[i].x > keys[j].x) {
			return false
		} else {
			return keys[i].y < keys[j].y
		}
	})
	
	for _, key := range keys {
		for _, direction := range []Direction {DirLeft, DirRight, DirDown, DirUp} {
			if value, exists := resultMap[PointAndDir{key, direction}]; exists {
				results = append(results, PushPoint{value, direction})
			}
		}
	}
	return results
}

func FindPushPoints(warehouse Warehouse, arr Arrangement) (PushPoints) {

	path := Path{}//make(Path, 0, 10)
	area := SearchArea{}
	area[arr.keeper] = Path{}
	
	resultMap := ResultMap{}
	for _, direction := range []Direction {DirLeft, DirRight, DirDown, DirUp} {
		findPushPointsRecursively(warehouse, arr, area, path, arr.keeper, direction, resultMap)
	}
	
	results := stableSortResultMap(resultMap)
	return results
}
