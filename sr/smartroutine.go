package sr

import (
	"crypto/md5"
	"encoding/hex"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var panics map[string][]string = make(map[string][]string)

const REDO_SECS = 30
const TYPE_PANIC_REDO = "panic_redo"
const TYPE_PANIC_RETURN = "panic_return"

type SmartR struct {
	todo    chan struct{}
	Context interface{}
	Type    string
	routine func(context interface{})
	done    chan struct{}
}

func New_Panic_Redo(routine_ func()) *SmartR {
	return newWithContext(TYPE_PANIC_REDO, nil, func(c interface{}) {
		routine_()
	})
}

func New_Panic_Return(routine_ func()) *SmartR {
	return newWithContext(TYPE_PANIC_RETURN, nil, func(c interface{}) {
		routine_()
	})
}

func New_Panic_RedoWithContext(Context_ interface{}, routine_ func(c interface{})) *SmartR {
	return newWithContext(TYPE_PANIC_REDO, Context_, routine_)
}

func New_Panic_ReturnWithContext(Context_ interface{}, routine_ func(c interface{})) *SmartR {
	return newWithContext(TYPE_PANIC_RETURN, Context_, routine_)
}

func newWithContext(Type_ string, Context_ interface{}, routine_ func(context interface{})) *SmartR {

	return &SmartR{
		todo:    make(chan struct{}, 1),
		Context: Context_,
		//panics:     make(map[string][]string),
		//panicExist: false,
		Type:    Type_,
		routine: routine_,
		done:    make(chan struct{}),
	}
}

func (sr *SmartR) recordPanicStack(panicstr string, stack string) {
	lines := strings.Split(stack, "\n")
	maxlines := len(lines)
	if maxlines >= 100 {
		maxlines = 100
	}

	errors := []string{panicstr}
	errors = append(errors, strconv.FormatInt(time.Now().Unix(), 10))
	errstr := panicstr

	if maxlines >= 3 {
		for i := 2; i < maxlines; i = i + 2 {
			fomatstr := strings.ReplaceAll(lines[i], "	", "")
			errstr = errstr + "#" + fomatstr
			errors = append(errors, fomatstr)
		}
	}

	h := md5.New()
	h.Write([]byte(errstr))
	errhash := hex.EncodeToString(h.Sum(nil))
	panics[errhash] = errors
}

func GetPanics() map[string][]string {
	return panics
}

func ClearPanics() {
	panics = make(map[string][]string)
}

func (sr *SmartR) Start() {

	go func() {
		sr.todo <- struct{}{}
		for {
			select {
			case <-sr.todo:
				go func() {
					defer func() {
						if err := recover(); err != nil {
							//record panic
							sr.recordPanicStack(err.(error).Error(), string(debug.Stack()))
							//check redo
							if sr.Type == TYPE_PANIC_REDO {
								time.Sleep(REDO_SECS * time.Second)
								sr.todo <- struct{}{}
							} else {
								sr.done <- struct{}{}
							}
						}
					}()
					sr.routine(sr.Context)
					sr.done <- struct{}{}
				}()
			case <-sr.done:
				return
			}
		}
	}()

}
