/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    main
 * @Date:    2022/3/14 2:37 下午
 * @package: mogoattr
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package main

import (
	"context"
	"fmt"
	"github.com/jager/hawox/contextx"
	"github.com/jager/hawox/mgoattr"
	"sync"
	"time"
)

func main() {
	ctx, cancel := contextx.Default()
	defer cancel()
	start := time.Now()
	mgoattr.Initialize(ctx, func(opt *mgoattr.Option) {
		opt.Addr = "119.29.105.154:9717"
		opt.DB = "test_attr"
	})

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		attr := mgoattr.NewAttrMgr("hello", i)
		wg.Add(1)
		ctx.Go(func(ctx context.Context) error {
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key_%d", j)
				val := int32(j + 1000)
				attr.SetInt32(key, val)
			}
			err := attr.Save(true)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
			return nil
		})
	}

	//for i := 0; i < 10; i++ {
	//	attr := mgoattr.NewAttrMgr("hello", i)
	//	err := attr.Load()
	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}
	//	//id := i
	//	wg.Add(1)
	//	ctx.Go(func(ctx context.Context) error {
	//		for j := 0; j < 100; j++ {
	//			key := fmt.Sprintf("key_%d", j)
	//			val := attr.GetInt32(key)
	//			attr.SetInt32(key, val+100000)
	//			//fmt.Printf("ID=%d Key=%d Get value=%d\n", id, j, val)
	//		}
	//		err = attr.Save(true)
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//		wg.Done()
	//		return nil
	//	})
	//}
	wg.Wait()
	cancel()
	fmt.Println(ctx.Wait(), time.Now().Sub(start).String())
}
