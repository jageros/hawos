/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    main
 * @Date:    2022/3/10 1:58 下午
 * @package: game
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/example/mygame/internal/session"
	"github.com/jageros/hawox/example/mygame/protos/meta"
	"github.com/jageros/hawox/example/mygame/protos/pb"
	"github.com/jageros/hawox/httpx"
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/ws"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	logx.Init(logx.DebugLevel)
	ctx, cancel := contextx.Default()
	defer cancel()
	httpx.InitializeHttpServer(ctx, func(engine *gin.Engine) {
		r := engine.Group("/ws")
		ws.Init(ctx, r)
	}, func(s *httpx.Server) {
		s.Port = 10088
	})
	ws.OnMessageBinary(func(ss *melody.Session, bytes []byte) {
		arg := &pb.PkgMsg{}
		err := arg.Unmarshal(bytes)
		if err != nil {
			logx.Error(err)
			return
		}
		sess := session.New(1001, 1000, 1)
		resp, err := meta.Call(sess, arg.Msgid, arg.Payload)
		if err != nil {
			if ierr, ok := err.(errcode.IErr); ok {
				reply := &pb.PkgMsg{
					Type: pb.MsgType_Err,
					Msgid: arg.Msgid,

				}
				data, _ := reply.Marshal()
				err = ss.WriteBinary(data)
				if err != nil {
					logx.Error(err)
				}
			}
		}
		return
	}
	if resp != nil {
		reply := &pb.RespMsg{
			Msgid:   arg.Msgid,
			Payload: resp,
		}
		data, _ := reply.Marshal()
		err = ss.WriteBinary(data)
		if err != nil {
			logx.Error(err)
		}
	}
})
player.RegisterRpcHandle()

logx.Infof("server stop with: %v", ctx.Wait())
logx.Sync()
}
