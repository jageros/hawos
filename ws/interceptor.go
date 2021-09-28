package ws

import (
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/examples/protos/pb"
	"github.com/jageros/hawox/jwt"
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/protoc"
	"github.com/jageros/hawox/protos/pbf"
	"time"
)

func interceptor(uid string, arg *pbf.Request) (*pbf.Response, bool) {
	switch arg.MsgID {
	case pb.MsgID_C2S_PING.ID():
		return protoc.CallFunc(uid, arg, ping), true

	case pb.MsgID_C2S_AUTH_TOKEN.ID():
		return protoc.CallFunc(uid, arg, auth), true
	}
	return nil, false
}

func ping(agent *protoc.Agent, arg interface{}) (interface{}, errcode.IErr) {
	arg2 := arg.(*pb.Ping)
	logx.Debugf("Ping uid=%s App-Id=%s", arg2.Uid, arg2.AppId)
	return &pb.Pong{
		Timestamp: time.Now().Unix(),
	}, nil
}

func auth(agent *protoc.Agent, arg interface{}) (interface{}, errcode.IErr) {
	arg2 := arg.(*pb.AuthMsg)
	token, err := jwt.GenerateToken(arg2.Uid)
	if err != nil {
		return nil, errcode.WithErrcode(-40, err)
	}

	return &pb.AuthResp{
		Token: token,
	}, nil
}
