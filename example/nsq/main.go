/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    queues_test
 * @Date:    2022/5/17 14:44
 * @package: queues
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package main

import (
	"fmt"
	"github.com/jager/hawox/contextx"
	"github.com/jager/hawox/queues"
	"time"
)

const nsqAddr = "127.0.0.1:4161"

func main() {
	ctx, cancel := contextx.Default()
	defer cancel()
	c, err := queues.NewConsumer(ctx, func(cfg *queues.Config) {
		cfg.Addrs = nsqAddr
		cfg.Topic = "jageros"
		cfg.Type = "nsq"
	})
	if err != nil {
		fmt.Println(err)
		cancel()
	}
	c.OnMessageHandler(func(data []byte) {
		fmt.Println(string(data))
	})

	p, err := queues.NewProducer(ctx, func(cfg *queues.Config) {
		cfg.Addrs = nsqAddr
		cfg.Type = "nsq"
		cfg.Topic = "jageros"
	})

	if err != nil {
		fmt.Println(err)
		cancel()
	}

	fmt.Println("init done.")

	for i := 0; i < 100; i++ {
		err = p.Push([]byte("jageros cccc"))
		if err != nil {
			fmt.Println("piblic: err: ", err)
		} else {
			fmt.Println("public successful!")
		}
		time.Sleep(time.Second)
	}

	fmt.Println("End ", ctx.Wait())
}
