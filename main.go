
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
	f, err = os.Open(os.Args[0])
	if err != nil {
		fmt.Printf("Cannot open '%s': %v\n", os.Args[0], err)
		return
	}
	
	_, err = BuildTextRoom(f)
	if err != nil {
		fmt.Printf("Cannot create room: %v\n", err)
		return
	}
}