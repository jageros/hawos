/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    kcp
 * @Date:    2021/11/4 1:34 下午
 * @package: kcp
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package kcp

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jageros/hawox/contextx"
	"github.com/xtaci/kcp-go"
	"net"
	"sync"
	"time"
)

var (
	session map[string]net.Conn
	mx      sync.Mutex
)

type Option struct {
	Addr         string
	Block        kcp.BlockCrypt
	DataShards   int
	ParityShards int
	Secret       string
	HandleConn   func(ctx context.Context, conn net.Conn) error
}

func defaultOption() *Option {
	op := &Option{
		Addr:         "0.0.0.0:9066",
		DataShards:   10,
		ParityShards: 3,
		Secret:       "cd803706cdc822e372fd7c73c0f109b9",
	}
	block, err := kcp.NewSalsa20BlockCrypt([]byte(op.Secret))
	if err == nil {
		op.Block = block
	}
	return op
}

func Server(ctx contextx.Context, opfs ...func(opt *Option)) {
	op := defaultOption()
	for _, opf := range opfs {
		opf(op)
	}
	listen, err := kcp.ListenWithOptions(op.Addr, op.Block, op.DataShards, op.ParityShards)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.Go(func(ctx_ context.Context) error {
		for {
			select {
			case <-ctx_.Done():
				return listen.Close()
			default:
			}
			conn, err := listen.Accept()
			if err != nil {
				return err
			}
			if op.HandleConn != nil {
				ctx.Go(func(ctx_ context.Context) error {
					return op.HandleConn(ctx_, conn)
				})
			}
		}
	})
}

func handleConnS(conn net.Conn) {
	for {
		fmt.Println("recv -----> ")
		datas := bytes.NewBuffer(nil)
		var buf [512]byte

		n, err := conn.Read(buf[0:])
		fmt.Println(n)
		datas.Write(buf[0:n])
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Print("datas : ")
		fmt.Println(string(datas.Bytes()))

	}
}

func send2Client(conn net.Conn) {
	for {
		conn.Write([]byte("xxxxxxxxxxxxxxxxx"))
		time.Sleep(time.Second)
	}
}

func Client() {
	conn, err := kcp.Dial("127.0.0.1:10086")
	if err != nil {
		fmt.Println("----", err)
		return
	}
	go handleConnC(conn)
	select {}
	//for {
	//	fmt.Println("send ------> ")
	//	ret, err2 := conn.Write([]byte("hello kcp!!"))
	//	if err2 != nil {
	//		fmt.Println(err2)
	//	} else {
	//		fmt.Println(ret)
	//	}
	//	time.Sleep(time.Second)
	//}
}

func handleConnC(conn net.Conn) {
	for {
		fmt.Println("recv -----> ")
		datas := bytes.NewBuffer(nil)
		var buf [512]byte

		n, err := conn.Read(buf[0:])
		datas.Write(buf[0:n])
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Print("datas : ")
		fmt.Println(datas.Bytes())
	}
}
