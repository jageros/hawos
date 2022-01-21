/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    udpx
 * @Date:    2022/1/18 5:06 下午
 * @package: udpx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package udpx

import (
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/logx"
	"net"
	"strconv"
	"time"
)

type Option struct {
	ListenIp       string // 监听IP
	Port           int    // 监听端口
	WriteTimeout   time.Duration
	MaxPkgSize     int32
	OnMsgHandle    func(addr *net.UDPAddr, data []byte) (resp []byte)
	OnBinaryHandle func(addr *net.UDPAddr, data []byte) (resp []byte)
}

func (s *Option) addr() string {
	return s.ListenIp + ":" + strconv.Itoa(s.Port)
}

func defaultServer() *Option {
	return &Option{
		ListenIp:     "",
		Port:         9055,
		WriteTimeout: time.Second * 10,
		MaxPkgSize:   4096,
	}
}

func Init(ctx contextx.Context, ops ...func(opt *Option)) error {
	s := defaultServer()
	for _, op := range ops {
		op(s)
	}

	addr, err := net.ResolveUDPAddr("udp", s.addr())
	if err != nil {
		return err
	}

	ctx.Go(func(ctx contextx.Context) error {
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			return err
		}
		ctx.Go(func(ctx contextx.Context) error {
			<-ctx.Done()
			return conn.SetDeadline(time.Now())
		})
		for {
			select {
			case <-ctx.Done():
				return conn.Close()
			default:

				data := make([]byte, s.MaxPkgSize)
				_, rAddr, err := conn.ReadFromUDP(data)
				if err != nil {
					logx.Infof("conn.ReadFromUDP err=%v", err)
					continue
				}
				pkg := &Package{}
				pkg.Unmarshal(data)

				var resp []byte
				switch pkg.Type {
				case TextMessage:
					if s.OnMsgHandle != nil {
						resp = s.OnMsgHandle(rAddr, pkg.Payload)
					}
				case BinaryMessage:
					if s.OnBinaryHandle != nil {
						resp = s.OnBinaryHandle(rAddr, pkg.Payload)
					}
				}

				if len(resp) > 0 {
					rPkg := &Package{
						Type:    BinaryMessage,
						Payload: resp,
					}
					err = conn.SetWriteDeadline(time.Now().Add(s.WriteTimeout))
					if err == nil {
						_, err = conn.WriteToUDP(rPkg.Marshal(), rAddr)
						if err != nil {
							logx.Infof("conn.WriteToUDP err=%v", err)
						}
					} else {
						logx.Infof("conn.SetWriteDeadline err=%v", err)
					}
				}
			}
		}
	})

	return nil
}
