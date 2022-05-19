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
	"github.com/jager/hawox/contextx"
	"github.com/jager/hawox/errcode"
	"github.com/jager/hawox/example/mygame/internal/service/player"
	"github.com/jager/hawox/example/mygame/internal/session"
	"github.com/jager/hawox/example/mygame/protos/meta"
	"github.com/jager/hawox/example/mygame/protos/pb"
	"github.com/jager/hawox/httpx"
	"github.com/jager/hawox/logx"
	"github.com/jager/hawox/ws"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	logx.Init(logx.DebugLevel)
	ctx, cancel := contextx.Default()
	defer cancel()

	player.RegisterRpcHandle()

	httpx.InitializeHttpServer(ctx, func(engine *gin.Engine) {
		r := engine.Group("/ws")
		ws.Init(ctx, r, func(opt ws.Options) {
			opt.HandleMessageBinary(handleMsg)
		})
	}, func(s *httpx.Option) {
		s.Port = 10088
	})

	logx.Info().Ints32("MsgID", meta.AllRegisteredMsgid())
	logx.Err(ctx.Wait()).Msg("Application Stop!")
	logx.Sync()
}

func handleMsg(ss *melody.Session, bytes []byte) {
	arg := &pb.PkgMsg{}
	err := arg.Unmarshal(bytes)
	if err != nil {
		logx.Err(err).Msg("pb.PkgMsg Unmarshal")
		return
	}
	sess := session.New(1001, 1000, 1)

	resp, pbErr := onClientMsg(sess, arg)
	var reply = &pb.PkgMsg{
		Msgid: arg.Msgid,
	}
	if pbErr != nil {
		data, _ := pbErr.Marshal()
		reply.Type = pb.MsgType_Err
		reply.Payload = data
	} else if resp != nil {
		reply.Type = pb.MsgType_Reply
		reply.Payload = resp
	}

	if reply.Type != pb.MsgType_Unknown {
		data, _ := reply.Marshal()
		err = ss.WriteBinary(data)
		if err != nil {
			logx.Err(err).Msg("WriteBinary")
		}
	}
}

func onClientMsg(ss *session.Session, arg *pb.PkgMsg) ([]byte, *pb.ErrMsg) {
	var er = &pb.ErrMsg{
		Code: 200,
		Msg:  "successful",
	}

	resp, err := meta.Call(ss, arg.Msgid, arg.Payload)
	if err != nil {
		if ierr, ok := err.(errcode.IErr); ok {
			er.Code = ierr.Code()
			er.Msg = ierr.ErrMsg()
		} else {
			er.Code = -1000
			er.Msg = err.Error()
		}
		return nil, er
	}

	if resp != nil {
		return resp, nil
	} else if arg.Type == pb.MsgType_Req {
		return nil, er
	}
	return nil, nil
}
