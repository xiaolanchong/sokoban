
package main

import (
	"time"
	"github.com/nsf/termbox-go"
)

const (
	backgroundColor = termbox.ColorBlack
	textColor = termbox.ColorWhite
	
	keeperSym = '@'
	crateSym  = 'O'
	slotSym   = 'x'
	wallSym   = '#'
)

type TermFeedback struct {
	warehouse	Warehouse
	slots		Entities
}

func render(c Coord, ch rune) {
	//w, h := termbox.Size()
	termbox.SetCell(int(c.x) + 1, int(c.y) + 1, ch, textColor, backgroundColor)
}

const  ( 
	pauseMs = 1
)

func (tf TermFeedback) Render(arr Arrangement) {
	termbox.Clear(backgroundColor, backgroundColor)
	
	for k, _ := range tf.slots {
		render(k, slotSym)
	}
	
	for k, _ := range tf.warehouse.walls {
		render(k, wallSym)
	}
	
	for k, _ := range arr.crates {
		render(k, crateSym)
	}
	
	render(arr.keeper, keeperSym)
	termbox.Flush()
	
	
	time.Sleep(pauseMs * time.Millisecond)
}

func (tf TermFeedback) Close() {
	termbox.Close()
}

func NewTermFeedback(warehouse Warehouse, slots Entities) (TermFeedback, error) {
	err := termbox.Init()
	if err != nil {
		return TermFeedback{}, err
	}
	tf := TermFeedback{warehouse, slots}
	return tf, nil
}
