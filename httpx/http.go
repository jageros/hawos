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
	"github.com/gin-contrib/cors"
	"github.com/jager/hawox/contextx"
	"github.com/jager/hawox/logx"
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

func InitializeHttpServer(ctx contextx.Context, registry func(engine *gin.Engine), opfs ...func(opt *Option)) {
	opt := defaultOption()

	for _, opf := range opfs {
		opf(opt)
	}
	opt.Mode = mode(opt.Mode)

	gin.SetMode(opt.Mode)
	engine := gin.New()
	gin.ForceConsoleColor()

	engine.Use(logx.GinLogger(), gin.Recovery(), cors.Default())

	if opt.RateTime > 0 {
		engine.Use(RateMiddleware(opt.RateTime))
	}

	registry(engine)

	addr := fmt.Sprintf("%s:%d", opt.ListenIp, opt.Port)
	s := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}
	ctx.Go(func(ctx context.Context) error {
		return s.ListenAndServe()
	})
	ctx.Go(func(ctx context.Context) error {
		<-ctx.Done()
		ctx2, cancel := context.WithTimeout(context.Background(), opt.CloseTimeout)
		defer cancel()
		return s.Shutdown(ctx2)
	})
}
