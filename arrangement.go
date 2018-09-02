package main

import (
	"fmt"
	"hash/crc32"
	"crypto/md5"
	"encoding/base64"
	)
	
type Direction = uint8

const (
	DirLeft Direction = iota
	DirRight
	DirDown
	DirUp
)

type Box struct {
	left	CoordUnit
	top		CoordUnit
	right	CoordUnit  // rigtest position
	bottom	CoordUnit  // lowest position of an entity
}

type Arrangement struct {
	crates	Entities
	keeper	Coord
}

// Static entities: walls and bounding box
type Warehouse struct {
	walls Entities
	bound Box
}

func moveWithinBox(pos Coord, bound Box, direction Direction) (Coord, error) {
	var dx, dy CoordUnit
	switch direction {
	case DirLeft:
		if pos.x <= bound.left {
			return Coord{}, fmt.Errorf("Hit the left bound %v at %v", bound.left, pos)
		}
		dx, dy = -1, 0
	case DirRight:
		if pos.x >= bound.right {
			return Coord{}, fmt.Errorf("Hit the right bound %v at %v", bound.right, pos)
		}
		dx, dy = 1, 0
	case DirUp:
		if pos.y <= bound.top {
			return Coord{}, fmt.Errorf("Hit the top bound %v at %v", bound.top, pos)
		}
		dx, dy = 0, -1
	case DirDown:
		if pos.y >= bound.bottom {
			return Coord{}, fmt.Errorf("Hit the bottom bound %v at %v", bound.bottom, pos)
		}
		dx, dy = 0, 1
	default:
		return Coord{}, fmt.Errorf("Unknown direction %v", direction)
	}
	return Coord{pos.x + dx, pos.y + dy}, nil
}

func moveWithinWalls(pos Coord, warehouse Warehouse, direction Direction) (Coord, error) {
	nextPos, err := moveWithinBox(pos, warehouse.bound, direction)
	if err != nil {
		return Coord{}, err
	}
	if _, exists := warehouse.walls[nextPos]; exists {
		return Coord{}, fmt.Errorf("Hit a wall, position=%v", nextPos)
	}
	return nextPos, nil
}

func moveWithinWallsAndCrates(pos Coord, warehouse Warehouse, 
							  crates Entities, direction Direction) (Coord, error) {
	nextPos, err := moveWithinWalls(pos, warehouse, direction)
	if err != nil {
		return nextPos, err
	}
	if _, exists := crates[nextPos]; exists {
		return Coord{}, fmt.Errorf("Hit a crate, position=%v", nextPos)
	}
	return nextPos, nil
}

func (a Arrangement) GetHash() uint32 {
	hasher := crc32.NewIEEE()
	for k, _ := range a.crates {
		hasher.Write([]byte{byte(k.x), byte(k.y)})
	}
	hasher.Write([]byte {byte(a.keeper.x), byte(a.keeper.y)})
	return hasher.Sum32()
}

func (a Arrangement) GetMd5Hash() string {
	hasher := md5.New()
	for k, _ := range a.crates {
		hasher.Write([]byte{byte(k.x), byte(k.y)})
	}
	hasher.Write([]byte {byte(a.keeper.x), byte(a.keeper.y)})
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func (a Arrangement) Move(warehouse Warehouse, direction Direction) (Arrangement, error) {
	var err error
	nextKeeperPos, err := moveWithinWalls(a.keeper, warehouse, direction)
	if err != nil {
		return Arrangement{}, err
	}
	
	if _, exists := a.crates[nextKeeperPos]; exists {
		// bump into a crate
		nextCratePos, err := moveWithinWallsAndCrates(nextKeeperPos, warehouse, a.crates, direction)
		if err != nil {
			return Arrangement{}, fmt.Errorf("Cannot move bumped crate in %v", nextKeeperPos)
		}
		
		aNext := Arrangement{make(Entities), a.keeper}
		for k, v := range a.crates {
			if k != nextKeeperPos {
				aNext.crates[k] = v
			}
		}
		aNext.crates[nextCratePos] = true
		aNext.keeper = nextKeeperPos
		return aNext, nil
	} else {
		// vacant place
		aNext := Arrangement{make(Entities), a.keeper}
		for k, v := range a.crates {
			aNext.crates[k] = v
		}
		aNext.keeper = nextKeeperPos
		return aNext, nil
	}
}

func isCrateStuckInCorner(warehouse Warehouse, arrangement Arrangement, slots Entities) bool {
	for cratePos, _ := range arrangement.crates {
		if _, slotExists := slots[cratePos]; slotExists {
			continue         // crate in a slot, no stucking
		}
		canMoveInPrevDir := true
		for _, dir := range []Direction {DirLeft, DirUp, DirRight, DirDown, DirLeft} {
			_, err := moveWithinWalls(cratePos, warehouse, dir)
			if !canMoveInPrevDir && err != nil {
				return true
			}
			canMoveInPrevDir = (err == nil)
		}
	}
	return false
}

// Checks the arrangement is terminal, i.e. no crate can be moved vertically or horizontally
func IsStuck(w Warehouse, a Arrangement) bool {
	for k, _ := range a.crates {
		canMoveInDirection := make([]bool, 4)
		for _, dir := range []Direction {DirLeft, DirRight, DirDown, DirUp} {
			_, err := moveWithinWallsAndCrates(k, w, a.crates, dir)
			canMoveInDirection[dir] = (err == nil)
		}
		canMoveHorz := canMoveInDirection[0] && canMoveInDirection[1]
		canMoveVert := canMoveInDirection[2] && canMoveInDirection[3]
		if canMoveVert || canMoveHorz {
			return false
		}
	}
	return true
}

func IsStuckExt(warehouse Warehouse, arrangement Arrangement, slots Entities) bool {
	return isCrateStuckInCorner(warehouse, arrangement, slots)
}

func IsSolved(a Arrangement, slots Entities) (bool, error) {
	if len(a.crates) < len(slots) {
		return false, fmt.Errorf("Number of slots %v greater than that of crates %v", len(slots), len(a.crates))
	}
	for pos, _ := range a.crates {
		if _, exists := slots[pos]; !exists {
			return false, nil
		}
	}
	
	return true, nil
}
