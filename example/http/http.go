/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    http
 * @Date:    2021/11/11 11:10 上午
 * @package: main
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jageros/hawox/flags"
	"github.com/jageros/hawox/httpx"
)

func main() {
	const appName = "http-demo"
	ctx, wait := flags.Parse(appName)
	defer wait()

	httpx.InitializeHttpServer(ctx, func(engine *gin.Engine) {
		r := engine.Group("/api")
		r.GET("/sayhello", func(c *gin.Context) {
			httpx.PkgMsgWrite(c, map[string]interface{}{"say": "hello world!"})
		})
	}, func(s *httpx.Option) {
		s.Mode = flags.Options.Mode
		s.ListenIp = flags.Options.HttpIp
		s.Port = flags.Options.HttpPort
	})
}
