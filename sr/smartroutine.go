package sr

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"runtime/debug"
	"strings"
	"time"
)

const REDO_SECS = 30

const TYPE_PANIC_REDO = "panic_redo"
const TYPE_PANIC_RETURN = "panic_return"

type SmartR struct {
	todo    chan struct{}
	Context interface{}
	errors  map[string][]string
	Type    string
	routine func(context interface{})
	done    chan struct{}
}

func New(Type_ string, routine_ func(context interface{})) (*SmartR, error) {
	return NewWithContext(Type_, routine_, nil)
}

func NewWithContext(Type_ string, routine_ func(context interface{}), Context_ interface{}) (*SmartR, error) {

	if Type_ != TYPE_PANIC_REDO && Type_ != TYPE_PANIC_RETURN {
		return nil, errors.New("wrong type")
	}

	return &SmartR{
		todo:    make(chan struct{}, 1),
		Context: Context_,
		errors:  make(map[string][]string),
		Type:    Type_,
		routine: routine_,
		done:    make(chan struct{}),
	}, nil
}

func (sr *SmartR) recordPanicStack(panicstr string, stack string) {
	lines := strings.Split(stack, "\n")
	maxlines := len(lines)
	if maxlines >= 100 {
		maxlines = 100
	}

	errors := []string{panicstr}
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
	sr.errors[errhash] = errors
}

func (sr *SmartR) GetALLErrors() map[string][]string {
	return sr.errors
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
