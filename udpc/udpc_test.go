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
	"github.com/jager/hawox/contextx"
	"github.com/jager/hawox/udpx"
	"net"
	"testing"
)

func TestReq(t *testing.T) {
	ctx, cancel := contextx.Default()
	defer cancel()
	addr := &net.UDPAddr{IP: net.IPv4(119, 29, 105, 154), Port: 9055}

	c, err := New(ctx, func(opt *ClientOption) {
		opt.TargetAddr = addr
		//opt.OnMsgHandle = onMsgHandle
	})
	if err != nil {
		t.Error(err)
	}
	err = c.SendTextMsg([]byte("xxxxxxxxx"))
	if err != nil {
		t.Error(err)
	}
	err = ctx.Wait()
	fmt.Printf("stop: %v", err)
}

func onMsgHandle(msgType udpx.MsgType, data []byte) {
}
