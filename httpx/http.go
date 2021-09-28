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

package httpx

import (
	"context"
	"fmt"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/logx"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func mode(value string) string {
	switch value {
	case gin.DebugMode, gin.ReleaseMode, gin.TestMode:
		return value
	case logx.InfoLevel, logx.WarnLevel, logx.ErrorLevel, logx.PanicLevel:
		return gin.ReleaseMode
	}
	return gin.DebugMode
}

type Server struct {
	ListenIp     string // 监听IP
	Port         int    // 监听端口
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	CloseTimeout time.Duration
	RateTime     time.Duration
	engine       *gin.Engine
	svr          *http.Server
}

func defaultServer() *Server {
	return &Server{
		ListenIp:     "",
		Port:         8888,
		Mode:         gin.DebugMode,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		CloseTimeout: time.Second * 10,
		RateTime:     time.Millisecond,
	}
}

func InitializeHttpServer(ctx contextx.Context, registry func(engine *gin.Engine), opfs ...func(s *Server)) {
	s := defaultServer()

	for _, opf := range opfs {
		opf(s)
	}
	s.Mode = mode(s.Mode)

	gin.SetMode(s.Mode)
	engine := gin.New()
	gin.ForceConsoleColor()
	engine.Use(logger(), gin.Recovery(), cors.Default())
	if s.RateTime > 0 {
		engine.Use(RateMiddleware(s.RateTime))
	}
	registry(engine)
	s.engine = engine

	addr := fmt.Sprintf("%s:%d", s.ListenIp, s.Port)
	s.svr = &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}
	ctx.Go(func(ctx contextx.Context) error {
		return s.svr.ListenAndServe()
	})
	ctx.Go(func(ctx contextx.Context) error {
		<-ctx.Done()
		ctx2, cancel := context.WithTimeout(context.Background(), s.CloseTimeout)
		defer cancel()
		return s.svr.Shutdown(ctx2)
	})
}
