/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    http
 * @Date:    2021/5/28 2:44 下午
 * @package: http
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package mygin

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/jageros/hawox/logx"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func mode(value string) string {
	switch value {
	case gin.DebugMode, gin.ReleaseMode, gin.TestMode:
		return value
	case "info", "warn", "error", "panic", "fatal":
		return gin.ReleaseMode
	}
	return gin.DebugMode
}

type Option struct {
	ListenIp     string // 监听IP
	Port         int    // 监听端口
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	CloseTimeout time.Duration
	RateTime     time.Duration
	Logger       gin.HandlerFunc

	Handler func(engine *gin.Engine) error
}

func defaultOption() *Option {
	return &Option{
		ListenIp:     "",
		Port:         8888,
		Mode:         gin.DebugMode,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		CloseTimeout: time.Second * 10,
		RateTime:     time.Millisecond,
	}
}

type Server struct {
	opt *Option
	*gin.Engine
	*http.Server
}

func NewServer(opfs ...func(opt *Option)) *Server {
	opt := defaultOption()
	for _, opf := range opfs {
		opf(opt)
	}
	return NewServerWithConfig(opt)
}

func NewServerWithConfig(opt *Option) *Server {
	opt.Mode = mode(opt.Mode)
	gin.SetMode(opt.Mode)
	engine := gin.New()
	gin.ForceConsoleColor()
	engine.Use(logx.GinLogger(), gin.Recovery(), cors.Default())
	if opt.RateTime > 0 {
		engine.Use(RateMiddleware(opt.RateTime))
	}
	return &Server{Engine: engine, opt: opt}
}

func (s *Server) Start(ctx context.Context) error {
	if s.opt.Handler == nil {
		return errors.New("Handler is nil")
	}
	err := s.opt.Handler(s.Engine)
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%d", s.opt.ListenIp, s.opt.Port)
	svr := &http.Server{
		Addr:         addr,
		Handler:      s.Engine,
		ReadTimeout:  s.opt.ReadTimeout,
		WriteTimeout: s.opt.WriteTimeout,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}
	s.Server = svr
	return svr.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	ctx2, cancel := context.WithTimeout(ctx, s.opt.CloseTimeout)
	defer cancel()
	return s.Shutdown(ctx2)
}

//func (s *Server) Endpoint() (*url.URL, error) {
//
//}
