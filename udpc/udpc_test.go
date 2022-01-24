/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    udpc_test
 * @Date:    2022/1/24 12:33 下午
 * @package: udpc
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package udpc

import (
	"fmt"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/udpx"
	"net"
	"testing"
)

func TestReq(t *testing.T) {
	ctx, cancel := contextx.Default()
	defer cancel()
	addr := &net.UDPAddr{IP: net.IPv4(119, 29, 105, 154), Port: 9055}
	OnRespMsgHandle(addr, func(msgType udpx.MsgType, data []byte) {
		fmt.Println(string(data))
	})
	addr2 := &net.UDPAddr{IP: net.IPv4(119, 29, 105, 154), Port: 9066}
	OnRespMsgHandle(addr2, func(msgType udpx.MsgType, data []byte) {
		fmt.Println(string(data))
	})
	OnGlobalRespHandle(func(addr *net.UDPAddr, msgType udpx.MsgType, data []byte) {
		fmt.Println(addr.String(), string(data))
	})
	ctx.Go(func(ctx contextx.Context) error {
		for i := 0; i < 100000; i++ {
			err := SendTextMsg(addr, []byte(fmt.Sprintf("Num=%d", i)))
			if err != nil {
				return err
			}
		}
		return nil
	})
	ctx.Go(func(ctx contextx.Context) error {
		for i := 100001; i < 200000; i++ {
			err := SendTextMsg(addr2, []byte(fmt.Sprintf("Num=%d", i)))
			if err != nil {
				return err
			}
		}
		return nil
	})

	ctx.Go(func(ctx contextx.Context) error {
		<-ctx.Done()
		return nil
	})

	err := ctx.Wait()
	fmt.Printf("stop: %v", err)
}
