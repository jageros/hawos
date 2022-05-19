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
	"context"
	"github.com/jager/hawox/contextx"
	"github.com/jager/hawox/udpx"
	"net"
	"time"
)

var (
	conn       *net.UDPConn
	handle     Handle
	maxPkgRead = 4096
)

type Option struct {
	RAddr        *net.UDPAddr
	MaxPkgSize   int
	WriteTimeout time.Duration
}

type Handle func(rAddr *net.UDPAddr, msgType udpx.MsgType, payload []byte)

type SetConfig func(maxPkgSize int)

func InitConn(ctx contextx.Context, lAddr *net.UDPAddr, f Handle) error {
	var err error
	conn, err = net.ListenUDP("udp", lAddr)
	if err != nil {
		return err
	}
	handle = f
	ctx.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return conn.Close()
	})
	read(ctx)
	return nil
}

func read(ctx contextx.Context) {
	ctx.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				resp := make([]byte, 4096)
				_, rAddr, err := conn.ReadFromUDP(resp)
				if err != nil {
					return err
				}
				pkg := udpx.GetPackage()
				pkg.Unmarshal(resp)
				if handle != nil {
					handle(rAddr, pkg.Type, pkg.Payload)
				}
				udpx.PutPackage(pkg)
			}
		}
	})
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

func LocalAddr() net.Addr {
	return conn.LocalAddr()
}
