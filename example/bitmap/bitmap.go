/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    bitmap
 * @Date:    2021/11/24 4:48 下午
 * @package: bitmap
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package main

import (
	"fmt"
	"github.com/jageros/hawox/bitmap"
	"github.com/jageros/hawox/contextx"
	"runtime"
)

func main() {
	ctx, cancel := contextx.Default()
	defer cancel()

	ch := make(chan int)
	ctx.Go(func(ctx contextx.Context) error {
		bm := bitmap.New()

		for i := 23214235; i < 23214235+10000000000; i++ {
			bm.Add(i)
		}

		fmt.Println(bm.Len())

		bm.Clear()
		ch <- 1
		return nil
	})

	ctx.Go(func(ctx contextx.Context) error {
		select {
		case <-ctx.Done():
			return nil
		case <-ch:
			runtime.GC()
			fmt.Println("Clear")
		}
		return nil
	})

	ctx.Go(func(ctx contextx.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})

	err := ctx.Wait()
	fmt.Println(err)
}
