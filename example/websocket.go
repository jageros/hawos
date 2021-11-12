/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    websocket
 * @Date:    2021/11/11 11:26 上午
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
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/ws"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	const appName = "ws-demo"
	ctx, wait := flags.Parse(appName)
	defer wait()

	httpx.InitializeHttpServer(ctx, func(engine *gin.Engine) {
		r := engine.Group("/ws")
		ws.Init(ctx, r, "/demo")
	}, func(s *httpx.Server) {
		s.Mode = flags.Options.Mode
		s.ListenIp = flags.Options.HttpIp
		s.Port = flags.Options.HttpPort
	})
	ws.OnConnect(func(sess *melody.Session) {
		logx.Infof("Ws Connect with keys=%+v", sess.Keys)
	})
	ws.OnMessage(func(sess *melody.Session, bytes []byte) {
		logx.Infof("Recv Msg=%s", string(bytes))
	})
	ws.OnDisConnect(func(session *melody.Session) {
		logx.Infof("Ws Disconnect.")
	})
}
