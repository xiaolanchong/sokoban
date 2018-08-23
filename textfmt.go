package main

import (
	"fmt"
	"io"
)

const (
	Wall  = 'x'
	Crate = 'o'
	Slot  = '.'
	Keeper = 'A'
	Empty = ' '
)

func BuildTextRoom(reader io.Reader) (room, error) {
	rb := NewBuilder()
	buf := make([]byte, 1)
	currentPos := Coord{0, 0}
	for {
		_, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		switch buf[0] {
		case Wall:
			rb.AddWall(currentPos)
			currentPos.x += 1
		case Crate:
			rb.AddCrate(currentPos)
			currentPos.x += 1
		case Slot:
			rb.AddSlot(currentPos)
			currentPos.x += 1
		case Keeper:
			rb.SetKeeper(currentPos)
			currentPos.x += 1
		case Empty:
			currentPos.x += 1
		case '\n':
			currentPos.x = 0
			currentPos.y += 1
		default:
			return room{}, fmt.Errorf("Unexpected char: %v", buf[0])
		}
	}

	return rb.Finish()
}