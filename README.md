# go-smart-routine
a safe and smart way to deal with go routine


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
	srh := sr.New_Redo(func() {
		fmt.Println("start of the program")
		if x == 0 {
			x++
			fmt.Println("here0")
			divide(10, 0)
		}
		fmt.Println("end of the program")
	})

	srh.Start()

	fmt.Println("///////////after 35 seconds//////////////")
	time.Sleep(35 * time.Second)
	errors := srh.GetALLErrors()
	for key, element := range errors {
		fmt.Println("error key :", key)
		fmt.Println("errors  :", element)
	}

}



```
