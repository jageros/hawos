/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    udpx_test
 * @Date:    2022/1/20 4:05 下午
 * @package: udpx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package udpx

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Server2 struct {
	ListenIp   string // 监听IP
	Port       int    // 监听端口
	Mode       string
	MaxPkgSize int32
}

func TestUdpx(t *testing.T) {
	data, err := json.Marshal(&Server2{
		ListenIp:   "192.168.31.168",
		Port:       9055,
		Mode:       "debug",
		MaxPkgSize: 6094,
	})
	if err != nil {
		t.Error(err)
	}
	arg := &Package{
		Type:    BinaryMessage,
		Payload: data,
	}

	ddd := make([]byte, 4096)

	copy(ddd, arg.Marshal())

	arg2 := &Package{}
	arg2.UnMarshal(ddd)

	data2 := &Server2{}
	err = json.Unmarshal(arg2.Payload, data2)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("msgType=%d data=%+v\n", arg2.Type, data2)
}
