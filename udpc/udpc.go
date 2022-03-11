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
	"context"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/udpx"
	"net"
	"time"
)

type RespHandle func(msgType udpx.MsgType, data []byte)

type ClientOption struct {
	MaxPkgSize   int
	WriteTimeout time.Duration
	OnMsgHandle  RespHandle
	SrcAddr      *net.UDPAddr
	TargetAddr   *net.UDPAddr
}

type Client struct {
	opt    *ClientOption
	conn   *net.UDPConn
	cancel contextx.CancelFunc
}

func New(ctx contextx.Context, opfs ...func(opt *ClientOption)) (*Client, error) {
	opt := &ClientOption{
		MaxPkgSize:   4096,
		WriteTimeout: time.Second * 5,
		OnMsgHandle: func(msgType udpx.MsgType, data []byte) {
			if msgType == udpx.TextMessage {
				logx.Infof("RespMsg=%s", string(data))
			} else {
				logx.Infof("RespMsg is binary.")
			}
		},
		SrcAddr: &net.UDPAddr{IP: net.IPv4zero, Port: 59055},
	}

	for _, opf := range opfs {
		opf(opt)
	}
	if opt.TargetAddr == nil {
		return nil, errcode.New(1, "Target Addr is nil")
	}
	conn, err := net.DialUDP("udp", opt.SrcAddr, opt.TargetAddr)
	if err != nil {
		return nil, err
	}
	ctx.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return conn.Close()
	})

	ctx_, cancel := contextx.WithCancel(ctx)
	ctx_.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				data := make([]byte, maxPkgRead)
				_, _, err = conn.ReadFromUDP(data)
				if err != nil {
					return err
				}
				pkg := udpx.GetPackage()
				pkg.Unmarshal(data)
				opt.OnMsgHandle(pkg.Type, pkg.Payload)
			}
		}
	})
	c := &Client{
		opt:    opt,
		conn:   conn,
		cancel: cancel,
	}
	return c, nil
}

func (c *Client) SendTextMsg(data []byte) error {
	pkg := udpx.GetPackage()
	pkg.Type = udpx.TextMessage
	pkg.Payload = data
	err := c.conn.SetWriteDeadline(time.Now().Add(c.opt.WriteTimeout))
	if err != nil {
		return err
	}
	_, err = c.conn.Write(pkg.Marshal())
	return err
}

func (c *Client) SendBinaryMsg(data []byte) error {
	pkg := udpx.GetPackage()
	pkg.Type = udpx.BinaryMessage
	pkg.Payload = data
	err := c.conn.SetWriteDeadline(time.Now().Add(c.opt.WriteTimeout))
	if err != nil {
		return err
	}
	_, err = c.conn.Write(pkg.Marshal())
	return err
}

func (c *Client) Close() {
	c.cancel()
}
