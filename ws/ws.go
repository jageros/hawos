/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    service
 * @Date:    2021/7/8 5:30 下午
 * @package: service
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/httpx"
	"github.com/jageros/hawox/jwt"
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/protos/pbf"
	"github.com/jageros/hawox/rpcx"
	"github.com/jageros/hawox/selector"
	"google.golang.org/grpc"
	"gopkg.in/olahol/melody.v1"
	"time"
)

var ss *service

type service struct {
	ctx         contextx.Context
	m           *melody.Melody
	callTimeout time.Duration
}

func OnMessage(fn func(session *melody.Session, bytes []byte)) {
	ss.m.HandleMessage(fn)
}

func OnMessageBinary(fn func(session *melody.Session, bytes []byte)) {
	ss.m.HandleMessageBinary(fn)
}

func OnConnect(fn func(*melody.Session)) {
	ss.m.HandleConnect(fn)
}

func OnDisConnect(fn func(*melody.Session)) {
	ss.m.HandleDisconnect(fn)
}

func Broadcast(data []byte, target *pbf.Target) error {
	if target == nil || target.GroupId["0"] {
		return ss.m.BroadcastBinary(data)
	}

	return ss.m.BroadcastBinaryFilter(data, func(session *melody.Session) bool {
		uidi, exist := session.Get("uid")
		if !exist {
			return false
		}
		uid := uidi.(string)
		roomIdi, exist := session.Get("roomId")
		if !exist {
			return false
		}
		roomId := roomIdi.(string)

		if target.UnlessUids[uid] {
			return false
		}
		if target.GroupId[roomId] {
			return true
		}
		if target.Uids[uid] {
			return true
		}
		return false
	})
}

func Init(ctx contextx.Context, r *gin.RouterGroup, relativePath string) {
	ss = &service{
		ctx:         ctx,
		m:           melody.New(),
		callTimeout: time.Second * 5,
	}
	ss.m.HandleMessageBinary(ss.handleMessage)
	r.GET(relativePath, jwt.CheckToken, ss.handler)
}

func (s *service) handler(c *gin.Context) {
	uid, ok := jwt.Uid(c)
	if !ok {
		return
	}
	roomId := c.GetHeader("X-RoomId")
	if roomId == "" {
		httpx.ErrInterrupt(c, errcode.InvalidParam)
		return
	}
	err := s.m.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{"uid": uid, "roomId": roomId})
	if err != nil {
		httpx.ErrInterrupt(c, errcode.WithErrcode(-1, err))
	}
}

func (s *service) handleMessage(session *melody.Session, bytes []byte) {
	start := time.Now()
	uid, exist := session.Get("uid")
	if !exist {
		return
	}
	arg := new(pbf.Request)
	err := arg.Unmarshal(bytes)
	if err != nil {
		logx.Errorf("OnClientRpcCall Unmarshal arg err: %v", err)
		return
	}
	var resp = new(pbf.Response)
	resp.MsgID = arg.MsgID
	defer func() {
		respData, err := resp.Marshal()
		if err == nil {
			err = session.WriteBinary(respData)
			if err != nil {
				logx.Error(err)
			}
		}
		take := time.Now().Sub(start)
		if take > time.Second {
			logx.Infof("Websocket Handle take: %s", take.String())
		}
		if resp.Code != errcode.Success.Code() {
			logx.Infof("Handle MsgID=%s Code=%d TakeTime=%s", resp.MsgID, resp.Code)
		}

	}()

	if err := frontendMiddleware(uid.(string), arg, session); err != nil {
		resp.Code = err.Code()
		return
	}

	if rs, ok := interceptor(uid.(string), arg); ok {
		resp.Code = rs.Code
		resp.Payload = rs.Payload
		return
	}

	tc, err := jwt.ParseToken(arg.Token)

	if err != nil {
		resp.Code = errcode.VerifyErr.Code()
		return
	}

	reqArg := &pbf.ReqArg{
		Uid:     tc.Uid,
		MsgID:   arg.MsgID,
		Payload: arg.Payload,
	}

	appNames := selector.GetName(arg.GetMsgID())

	if len(appNames) <= 0 {
		resp.Code = errcode.ServiceNotFound.Code()
		return
	}

	var ierr errcode.IErr
	var respMsgs []*pbf.RespMsg
	for _, appName := range appNames {
		key := fmt.Sprintf("appid_%s", appName)
		id, exist := session.Get(key)
		if exist {
			appName = fmt.Sprintf("%s/%s", appName, id)
		}
		err_ := rpcx.CallByName(appName, func(cc *grpc.ClientConn) {
			rt := pbf.NewRouterClient(cc)
			ctx, cancel := s.ctx.WithTimeout(s.callTimeout)
			defer cancel()
			respMsg2, err2 := rt.ReqCall(ctx, reqArg)
			if err2 == nil {
				respMsgs = append(respMsgs, respMsg2)
				if !exist {
					session.Set(key, respMsg2.AppId)
				}
			} else {
				logx.Errorf("Rpc Call Name=%s err=%v", appName, err2)
				err = err2
			}
		})
		if err_ != nil {
			ierr = err_
		}
	}

	if ierr != nil {
		resp.Code = ierr.Code()
		return
	}

	if err != nil {
		err2 := errcode.WithErrcode(-60, err)
		resp.Code = err2.Code()
		logx.Errorf("ReqCall appName=%v MsgId=%s err:%v", appNames, arg.MsgID, err2)
		return
	}

	resp.Code = respMsgs[0].Code
	resp.Payload = respMsgs[0].GetPayload()

	if err := backendMiddleware(uid.(string), resp, session); err != nil {
		resp.Code = err.Code()
		return
	}

}
