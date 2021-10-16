# go-smart-routine
#### a safe and smart way to deal with go routine
#### panic_redo type will make the routine redo again after 30 seconds if any panic happens
#### panic_return type will quit if any panic happens
#### all panics will be recorded with unix-timestamp


```go

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



```
