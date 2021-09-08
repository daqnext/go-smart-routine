package main

import (
	"fmt"
	"time"

	"github.com/daqnext/go-smart-routine/sr"
)

func divide(a, b int) int {
	return a / b
}

var Pncs = make(map[string][]string)

func main() {

	// it takes about 30 seconds for a restart of a panic routine
	x := 0
	sr.New_Panic_Redo(func() {
		fmt.Println("start of the program")
		if x == 0 {
			x++
			fmt.Println("here0")
			divide(10, 0)
		}
		fmt.Println("end of the program")
	}).Start()

	fmt.Println("///////////after 35 seconds//////////////")
	time.Sleep(35 * time.Second)

	panics := sr.GetPanics()
	for key, element := range panics {
		fmt.Println("panic key :", key)
		fmt.Println("panics  :", element)
	}
	sr.ClearPanics()
}
