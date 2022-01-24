/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    default
 * @Date:    2022/1/24 4:24 下午
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
)

var (
	conn       *net.UDPConn
	handles    func(addr *net.UDPAddr, msgType udpx.MsgType, data []byte)
	maxPkgRead = 4096
)

func OnMsgHandle(f func(addr *net.UDPAddr, msgType udpx.MsgType, data []byte)) {
	handles = f
}

func DefaultConn() *net.UDPConn {
	if conn != nil {
		return conn
	}
	var err error
	conn, err = net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		logx.Fatalf("udp client init err: %v", err)
	}
	go onMsg(conn)
	return conn
}

func onMsg(con *net.UDPConn) {
	for {
		data := make([]byte, maxPkgRead)
		_, addr, err := con.ReadFromUDP(data)
		if err == nil {
			pkg := udpx.GetPackage()
			pkg.Unmarshal(data)
			if handles != nil {
				handles(addr, pkg.Type, pkg.Payload)
			}
		}
	}
}

func SendTextMsg(addr *net.UDPAddr, data []byte) error {
	conn2 := DefaultConn()
	pkg := udpx.GetPackage()
	pkg.Type = udpx.TextMessage
	pkg.Payload = data
	_, err := conn2.WriteToUDP(pkg.Marshal(), addr)
	return err
}

func SendBinaryMsg(addr *net.UDPAddr, data []byte) error {
	conn2 := DefaultConn()
	pkg := udpx.GetPackage()
	pkg.Type = udpx.BinaryMessage
	pkg.Payload = data
	_, err := conn2.WriteToUDP(pkg.Marshal(), addr)
	return err
}

func LocalAddr() net.Addr {
	conn2 := DefaultConn()
	return conn2.LocalAddr()
}
