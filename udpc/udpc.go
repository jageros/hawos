/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    udpc
 * @Date:    2022/1/24 11:50 上午
 * @package: udpc
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package udpc

import (
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/udpx"
	"net"
	"sync"
)

var (
	maxPkgRead   = 4096
	conn         *net.UDPConn
	handles      map[string]RespHandle
	globalHandle GlobalRespHandle
	rwMux        sync.RWMutex
)

func SetMaxPkgRead(n int) {
	maxPkgRead = n
}

type RespHandle func(msgType udpx.MsgType, data []byte)
type GlobalRespHandle func(addr *net.UDPAddr, msgType udpx.MsgType, data []byte)

func init() {
	var err error
	conn, err = net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		logx.Fatalf("udp client init err: %v", err)
	}
	handles = map[string]RespHandle{}
	go onMsg()
}

func onMsg() {
	for {
		data := make([]byte, maxPkgRead)
		_, addr, err := conn.ReadFromUDP(data)
		if err == nil {
			pkg := udpx.GetPackage()
			pkg.Unmarshal(data)
			rwMux.RLock()
			h, ok := handles[addr.String()]
			rwMux.RUnlock()
			if ok {
				h(pkg.Type, pkg.Payload)
			}
			if globalHandle != nil {
				globalHandle(addr, pkg.Type, pkg.Payload)
			}
		}
	}
}

func SendTextMsg(addr *net.UDPAddr, data []byte) error {
	pkg := udpx.GetPackage()
	pkg.Type = udpx.TextMessage
	pkg.Payload = data
	_, err := conn.WriteToUDP(pkg.Marshal(), addr)
	return err
}

func SendBinaryMsg(addr *net.UDPAddr, data []byte) error {
	pkg := udpx.GetPackage()
	pkg.Type = udpx.BinaryMessage
	pkg.Payload = data
	_, err := conn.WriteToUDP(pkg.Marshal(), addr)
	return err
}

func OnRespMsgHandle(addr *net.UDPAddr, h RespHandle) {
	rwMux.Lock()
	handles[addr.String()] = h
	rwMux.Unlock()
}

func OnGlobalRespHandle(h GlobalRespHandle) {
	globalHandle = h
}

func LocalAddr() net.Addr {
	return conn.LocalAddr()
}
