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

type Options interface {
	SetCallTimeout(t time.Duration)
	SetRelativePath(relativePath string)
	SetKeys(keys []string)
	SetAuth(auth func(c *gin.Context))
}

type service struct {
	ctx          contextx.Context
	m            *melody.Melody
	callTimeout  time.Duration
	relativePath string
	keys         []string
	auth         func(c *gin.Context)
}

type Filter struct {
	All        bool
	UnlessKeys map[string]map[interface{}]struct{}
	Keys       map[string]map[interface{}]struct{}
}

func (f *Filter) AddKeys(key string, vals ...interface{}) {
	if f.Keys == nil {
		f.Keys = map[string]map[interface{}]struct{}{
			key: {},
		}
	} else if _, ok := f.Keys[key]; !ok {
		f.Keys[key] = map[interface{}]struct{}{}
	}
	for _, v := range vals {
		f.Keys[key][v] = struct{}{}
	}
}

func (f *Filter) AddUnlessKeys(key string, vals ...interface{}) {
	if f.UnlessKeys == nil {
		f.UnlessKeys = map[string]map[interface{}]struct{}{
			key: {},
		}
	} else if _, ok := f.UnlessKeys[key]; !ok {
		f.UnlessKeys[key] = map[interface{}]struct{}{}
	}
	for _, v := range vals {
		f.UnlessKeys[key][v] = struct{}{}
	}
}

func NewFilter(all bool) *Filter {
	filter := &Filter{
		All:        all,
		UnlessKeys: map[string]map[interface{}]struct{}{},
		Keys:       map[string]map[interface{}]struct{}{},
	}
	return filter
}

func (s *service) SetCallTimeout(t time.Duration) {
	s.callTimeout = t
}
func (s *service) SetRelativePath(relativePath string) {
	s.relativePath = relativePath
}
func (s *service) SetKeys(keys []string) {
	s.keys = keys
}

func (s *service) SetAuth(auth func(c *gin.Context)) {
	s.auth = auth
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

func Broadcast(data []byte, filter *Filter) error {
	if filter == nil || filter.All {
		return ss.m.BroadcastBinary(data)
	}

	return ss.m.BroadcastBinaryFilter(data, func(session *melody.Session) bool {

		for k, v := range session.Keys {
			if vs, ok := filter.UnlessKeys[k]; ok {
				if _, ok := vs[v]; ok {
					return false
				}
			}
		}

		for k, v := range session.Keys {
			if vs, ok := filter.Keys[k]; ok {
				if _, ok := vs[v]; ok {
					return true
				}
			}
		}

		return false
	})
}

func Init(ctx contextx.Context, r *gin.RouterGroup, opfs ...func(opt Options)) {
	ss = &service{
		ctx:         ctx,
		m:           melody.New(),
		callTimeout: time.Second * 5,
		//auth:         jwt.CheckToken,
		//Keys:         []string{"uid", "roomId"},
		relativePath: "/gate",
	}
	for _, opf := range opfs {
		opf(ss)
	}
	ss.m.HandleMessageBinary(ss.handleMessage)
	if ss.auth != nil {
		r.GET(ss.relativePath, ss.auth, ss.handler)
	} else {
		r.GET(ss.relativePath, ss.handler)
	}

}

func (s *service) handler(c *gin.Context) {
	var keys = map[string]interface{}{}
	for _, k := range s.keys {
		var vi interface{}
		v := c.GetHeader(k)
		if v != "" {
			vi = v
		} else {
			vi, _ = c.Get(k)
		}
		if vi != nil {
			keys[k] = v
		}
	}

	err := s.m.HandleRequestWithKeys(c.Writer, c.Request, keys)
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
