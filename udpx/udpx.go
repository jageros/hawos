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
	"github.com/gin-gonic/gin"
	"github.com/jageros/hawox/contextx"
	"net"
	"strconv"
	"time"
)

type Server struct {
	ListenIp     string // 监听IP
	Port         int    // 监听端口
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	CloseTimeout time.Duration
	MaxPkgSize   int32
	OnMsgHandle  func(session *Session, data []byte) (resp []byte)

	conn   *net.UDPConn
	ctx    contextx.Context
	cancel contextx.CancelFunc
}

func (s *Server) addr() string {
	return s.ListenIp + ":" + strconv.Itoa(s.Port)
}

func defaultServer() *Server {
	return &Server{
		ListenIp:     "",
		Port:         8888,
		Mode:         gin.DebugMode,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		CloseTimeout: time.Second * 10,
		MaxPkgSize:   4096,
	}
}

func (s *Server) updateReadTime() error {
	return s.conn.SetReadDeadline(time.Now().Add(s.ReadTimeout))
}

func (s *Server) run() {
	s.ctx.Go(func(ctx contextx.Context) error {
		for {
			select {
			case <-ctx.Done():
				return s.conn.Close()
			default:

			}
		}
	})
}

func Init(ctx contextx.Context, ops ...func(opt *Server)) error {
	s := defaultServer()
	for _, op := range ops {
		op(s)
	}
	ctx_, cancel := ctx.WithCancel()
	s.ctx = ctx_
	s.cancel = cancel

	addr, err := net.ResolveUDPAddr("udp", s.addr())
	if err != nil {
		return err
	}
	s.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	s.run()
	return nil
}
