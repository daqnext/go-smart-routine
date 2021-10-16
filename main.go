package main

import (
	"fmt"
	"time"

	localLog "github.com/daqnext/LocalLog/log"
	"github.com/daqnext/go-smart-routine/sr"
)

func divide(a, b int) int {
	return a / b
}

func main() {

	lg, err := localLog.New("logs", 10, 10, 10)
	if err != nil {
		panic(err)
	}

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
	}, lg).Start()

	fmt.Println("///////////after 35 seconds//////////////")
	time.Sleep(35 * time.Second)

	// if sr.PanicExist {
	// 	fmt.Println(sr.PanicJson.GetContentAsString())
	// 	sr.ClearPanics()
	// }

}
