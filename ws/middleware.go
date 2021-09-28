package ws

import (
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/examples/protos/pb"
	"github.com/jageros/hawox/protoc"
	"github.com/jageros/hawox/protos/pbf"
	"gopkg.in/olahol/melody.v1"
)

func frontendMiddleware(uid string, arg *pbf.Request, session *melody.Session) errcode.IErr {
	return nil
}

func backendMiddleware(uid string, arg *pbf.Response, session *melody.Session) errcode.IErr {
	switch arg.MsgID {
	case pb.MsgID_C2S_ENTER_ROOM.ID(), pb.MsgID_C2S_CREATE_ROOM.ID():
		return protoc.RespFun(uid, arg, func(agent *protoc.Agent, arg_ interface{}) errcode.IErr {
			arg2 := arg_.(*pb.RoomInfo)
			session.Set("roomId", arg2.RoomId)
			return nil
		})
	}
	return nil
}
