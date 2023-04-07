package recovers

import (
	"fmt"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/logx"
	"reflect"
	"runtime"
)

func CatchPanic(f func() error) (err error) {
	defer func() {
		err1 := recover()
		if err1 != nil {
			fn := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
			logx.Errorf("%s panic: %v", fn, err1)
			err = errcode.New(1, fmt.Sprintf("%v", err1))
		}
	}()

	err = f()
	return
}

func Catch(f func()) {
	defer func() {
		err1 := recover()
		if err1 != nil {
			fn := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
			fmt.Println(fn, err1)
			logx.Errorf("%s panic: %v", fn, err1)
		}
	}()
	f()
	return
}

func Go(f func()) {
	go Catch(f)
}
