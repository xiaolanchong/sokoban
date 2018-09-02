package main

import (
	"fmt"
	"io"
)

const (
	Wall  = '#'
	Crate = 'O'
	Slot  = 'X'
	Keeper = '@'
	Empty = ' '
)

func BuildTextStage(reader io.Reader) (stage, error) {
	sb := NewBuilder()
	buf := make([]byte, 1)
	currentPos := Coord{0, 0}
	for {
		_, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		switch buf[0] {
		case Wall:
			sb.AddWall(currentPos)
			currentPos.x += 1
		case Crate:
			sb.AddCrate(currentPos)
			currentPos.x += 1
		case Slot:
			sb.AddSlot(currentPos)
			currentPos.x += 1
		case Keeper:
			sb.SetKeeper(currentPos)
			currentPos.x += 1
		case Empty:
			sb.AddSpace(currentPos)
			currentPos.x += 1
		case '\n':
			currentPos.x = 0
			currentPos.y += 1
		case '\r':
		default:
			return stage{}, fmt.Errorf("Unexpected char: %x, position: %v", buf, currentPos)
		}
	}

	return sb.Finish()
}
