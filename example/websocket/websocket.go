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
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jager/hawox/flags"
	"github.com/jager/hawox/httpx"
	"github.com/jager/hawox/logx"
	"github.com/jager/hawox/ws"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	const appName = "ws-demo"
	ctx, wait := flags.Parse(appName)
	defer wait()

	httpx.InitializeHttpServer(ctx, func(engine *gin.Engine) {
		r := engine.Group("/ws")
		ws.Init(ctx, r, func(opt ws.Options) {
			opt.SetRelativePath("/ss")
			opt.HandleConnect(func(sess *melody.Session) {
				logx.Debug().Msgf("Ws Connect with keys=%+v", sess.Keys)
			})
			opt.HandleDisconnect(func(sess *melody.Session) {
				logx.Info().Msg("Ws Disconnect.")
			})
			opt.HandleMessage(onMsg)
		})
	}, func(s *httpx.Option) {
		s.Mode = flags.Options.Mode
		s.ListenIp = "0.0.0.0"
		s.Port = 8088
	})
}

func onMsg(sess *melody.Session, bytes []byte) {
	m := map[string]interface{}{}
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		logx.Err(err).Msg("json.Unmarshal")
		return
	}
	logx.Info().Msgf("%+v", m)
	err = ws.Broadcast(bytes, nil)
	if err != nil {
		logx.Err(err).Msg("Ws Broadcast")
	}
}
