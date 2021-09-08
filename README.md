# go-smart-routine
#### a safe and smart way to deal with go routine
#### panic_redo type will make the routine redo again after 30 seconds if any panic happens
#### panic_return type will quit if any panic happens
#### all panics will be recorded with unix-timestamp


```go
import (
	"github.com/daqnext/go-smart-routine/sr"
)
```

### use-case:

```go



func divide(a, b int) int {
	return a / b
}



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


```
