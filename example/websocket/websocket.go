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
		ws.Init(ctx, r, func(opt ws.Options) {
			opt.SetRelativePath("/ss")
			//opt.SetKeys("uid")
			//opt.SetAuth(func(c *gin.Context) {
			//	_, ok := httpx.DecodeUrlVal(c, "uid")
			//	if !ok {
			//		httpx.ErrInterrupt(c, errcode.InvalidParam)
			//		return
			//	}
			//})
		})
	}, func(s *httpx.Option) {
		s.Mode = flags.Options.Mode
		s.ListenIp = flags.Options.HttpIp
		s.Port = flags.Options.HttpPort
	})
	ws.OnConnect(func(sess *melody.Session) {
		logx.Infof("Ws Connect with keys=%+v", sess.Keys)
	})
	ws.OnMessage(func(sess *melody.Session, bytes []byte) {
		//msg := string(bytes)
		//uid, _ := sess.Get("uid")
		//logx.Infof("Recv Msg=%s", msg)
		//filter := ws.NewBlacklistFilter("uid", uid)
		//err := ws.Broadcast(msg, filter)
		//if err != nil {
		//	logx.Errorf("ws Broadcast err: %v", err)
		//}
		m := map[string]interface{}{}
		err := json.Unmarshal(bytes, &m)
		if err != nil {
			logx.Errorf("json.Unmarshal err:%v", err)
			return
		}
		logx.Infof("%+v", m)
		err = ws.Broadcast(bytes, nil)
		if err != nil {
			logx.Errorf("Broadcast err: %v", err)
		}
	})
	ws.OnDisConnect(func(session *melody.Session) {
		logx.Infof("Ws Disconnect.")
	})
}
