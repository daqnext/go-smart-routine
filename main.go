package main

import (
	"fmt"
	"time"

	"github.com/daqnext/go-smart-routine/sr"
)

func divide(a, b int) int {
	return a / b
}

func main() {

	// it takes about 30 seconds for a restart of a panic routine
	x := 0
	srh, err := sr.New(sr.TYPE_PANIC_REDO, func(c interface{}) {
		fmt.Println("start of the program")
		if x == 0 {
			x++
			fmt.Println("here0")
			divide(10, 0)
		}
		fmt.Println("end of the program")
	})
	if err == nil {
		srh.Start()
	}

	fmt.Println("///////////after 35 seconds//////////////")
	time.Sleep(35 * time.Second)
	errors := srh.GetALLErrors()
	for key, element := range errors {
		fmt.Println("error key :", key)
		fmt.Println("errors  :", element)
	}

}
