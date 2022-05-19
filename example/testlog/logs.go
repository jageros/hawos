/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    logs
 * @Date:    2022/3/22 10:48 AM
 * @package: logs
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jager/hawox/errcode"
	"github.com/jager/hawox/flags"
	"github.com/jager/hawox/httpx"
	"github.com/jager/hawox/logx"
)

const (
	appName = "tslog"
	port    = 7777
)

func main() {
	ctx, wait := flags.Parse(appName, func(opt *flags.Option) {
		opt.LogDir = "logs/test.log"
	})
	defer wait()

	httpx.InitializeHttpServer(ctx, handles, func(opt *httpx.Option) {
		opt.ListenIp = "0.0.0.0"
		opt.Port = port
		opt.Mode = flags.Options.Mode
	})

	logx.Debug().Str("service", "http").Str("listen_ip", "0.0.0.0").Int("port", port).Str("mode", flags.Options.Mode).Send()
}

func handles(engine *gin.Engine) {
	r := engine.Group("ts")
	r.GET("/err", func(ctx *gin.Context) {
		err1 := errcode.New(1, "err1")
		err2 := errcode.New(2, "err2")
		err3 := errcode.New(3, "err3")
		var err errcode.IErr
		if httpx.HasErr(ctx, err1, err2, err3, err) {
			return
		}
		httpx.PkgMsgWrite(ctx, nil)
	})

	r.GET("/say", func(ctx *gin.Context) {
		httpx.PkgMsgWrite(ctx, map[string]interface{}{"msg": "Hello world!"})
	})

	r.GET("/healthcheck", func(ctx *gin.Context) {
		httpx.PkgMsgWrite(ctx, nil)
	})
}
