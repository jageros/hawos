/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    test
 * @Date:    2022/3/29 17:31
 * @package: test
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package main

import (
	"fmt"
	"sync"
	"time"
)

type Cond struct {
	*sync.Cond
}

func (c *Cond) Wait() {
	c.L.Lock()
	c.L.Unlock()
}

func NewCond() *Cond {
	l := &sync.RWMutex{}
	c := sync.NewCond(l)
	return &Cond{c}
}

func main() {
	w := sync.WaitGroup{}
	cc := NewCond()
	for i := 0; i < 10; i++ {
		w.Add(1)
		go func(n int) {
			defer w.Done()
			cc.Wait()
			time.Sleep(time.Second)
			fmt.Println(time.Now().String(), "xxxx", n)
		}(i)
	}
	time.Sleep(time.Second)
	cc.Broadcast()
	w.Wait()
}
