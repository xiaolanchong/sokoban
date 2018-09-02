
package main

import (
	"os"
	"fmt"
	"io"
)

func main() {
	var err error
	if (len(os.Args) < 2) {
		fmt.Printf("Usage: %s <input file>\n", os.Args[0])
		return
	}
	
	var f io.Reader
	f, err = os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Cannot open '%s': %v\n", os.Args[0], err)
		return
	}
	
	var st stage
	st, err = BuildTextStage(f)
	if err != nil {
		fmt.Printf("Cannot create stage: %v\n", err)
		return
	}
	
	tf, _ := NewTermFeedback(st.warehouse, st.slots)
	defer tf.Close()
	//Solve(st.warehouse, st.arrangement, st.slots, tf)
	SolveWithPathfinding(st.warehouse, st.arrangement, st.slots, tf)
}
