/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    wsapp
 * @Date:    2021/6/10 11:39 上午
 * @package: ws
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package protoc

import (
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/flags"
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/protos/meta"
	"github.com/jageros/hawox/protos/pbf"
)

var (
	agentRpcHandlers = make(map[int32]AgentRpcHandler)
	msgIds           []int32
)

type AgentRpcHandler func(agent *Agent, arg interface{}) (interface{}, errcode.IErr)
type AgentRespHandler func(agent *Agent, arg interface{}) errcode.IErr

func RegisterAgentRpcHandler(msgID meta.IMsgID, handler AgentRpcHandler) {
	agentRpcHandlers[msgID.ID()] = handler
	msgIds = append(msgIds, msgID.ID())
}

func MsgIDs() []int32 {
	return msgIds
}

// OnClientRpcCall req: *pbf.Response  resp: *pbf.Response
//func OnClientRpcCall(sessID string, data []byte, w io.Writer) {
//	arg := new(pbf.Response)
//	err := arg.Unmarshal(data)
//	if err != nil {
//		log.Infof("OnClientRpcCall Unmarshal arg err: %v", err)
//		return
//	}
//	pa := getAgent(sessID)
//	var resp = new(pbf.Response)
//	resp.MsgID = arg.MsgID
//	defer func() {
//		respData, err := resp.Marshal()
//		if err == nil {
//			w.Write(respData)
//		}
//		putAgent(pa)
//	}()
//	if handleFun, ok := agentRpcHandlers[arg.GetMsgID()]; ok {
//		im := GetMeta(arg.GetMsgID())
//		if im == nil {
//			log.Infof("Service NotFound MsgID=%d", arg.GetMsgID())
//			resp.Code = pbf.ErrCode_ServiceNotFound
//			return
//		}
//		arg2, err := im.DecodeArg(arg.Payload)
//		if err != nil {
//			resp.Code = pbf.ErrCode_InternalErr
//			return
//		}
//		reply, herr := handleFun(pa, arg2)
//		if herr != nil {
//			resp.Code = herr.ECode()
//			return
//		}
//		payload, err := im.EncodeReply(reply)
//		if err != nil {
//			resp.Code = pbf.ErrCode_InternalErr
//		} else {
//			resp.Code = pbf.ErrCode_Success
//			resp.Payload = payload
//		}
//	} else {
//		resp.Code = NoRpcHandleErr.ECode()
//	}
//}

// OnRouterRpcCall req: *pbf.ReqArg  resp: *pbf.RespMsg
func OnRouterRpcCall(arg *pbf.ReqArg) (*pbf.RespMsg, errcode.IErr) {

	pa := getAgent(arg.Uid)
	defer func() {
		putAgent(pa)
	}()

	var resp = new(pbf.RespMsg)
	resp.MsgID = arg.MsgID
	resp.AppId = flags.AppID()

	pa.Uid = arg.Uid

	if handleFun, ok := agentRpcHandlers[arg.GetMsgID()]; ok {
		im, err := meta.GetMeta(arg.GetMsgID())
		if err != nil {
			logx.Infof("Service NotFound MsgID=%d", arg.GetMsgID())
			return nil, errcode.MetaCoderNotFound
		}

		arg2, err := im.DecodeArg(arg.Payload)
		if err != nil {
			return nil, errcode.WithErrcode(-33, err)
		}

		reply, herr := handleFun(pa, arg2)
		if herr != nil {
			resp.Code = herr.Code()
			return resp, nil
		}

		payload, err := im.EncodeReply(reply)
		if err != nil {
			return nil, errcode.WithErrcode(-44, err)
		}

		resp.Code = errcode.Success.Code()
		resp.Payload = payload
		return resp, nil
	}

	resp.Code = errcode.ProtoMsgIdNoHandles.Code()
	return resp, nil
}

// CallFunc req: *pbf.Response  resp: *pbf.Response
func CallFunc(uid string, arg *pbf.Request, handle AgentRpcHandler) *pbf.Response {

	pa := getAgent(uid)
	defer func() {
		putAgent(pa)
	}()
	pa.Uid = uid

	var resp = new(pbf.Response)
	resp.MsgID = arg.MsgID

	im, err := meta.GetMeta(arg.GetMsgID())
	if err != nil {
		logx.Infof("Service NotFound MsgID=%d", arg.GetMsgID())
		resp.Code = errcode.MetaCoderNotFound.Code()
		return resp
	}

	arg2, err := im.DecodeArg(arg.Payload)
	if err != nil {
		logx.Infof("MsgID=%d req arg decode err: %v", arg.GetMsgID(), err)
		resp.Code = errcode.InternalErr.Code()
		return resp
	}

	reply, herr := handle(pa, arg2)
	if herr != nil {
		resp.Code = herr.Code()
		return resp
	}

	payload, err := im.EncodeReply(reply)
	if err != nil {
		logx.Infof("MsgID=%d resp arg encode err: %v", arg.GetMsgID(), err)
		resp.Code = errcode.InternalErr.Code()
		return resp
	}

	resp.Code = errcode.Success.Code()
	resp.Payload = payload
	return resp
}

func RespFun(uid string, arg *pbf.Response, handle AgentRespHandler) errcode.IErr {

	pa := getAgent(uid)
	defer func() {
		putAgent(pa)
	}()

	pa.Uid = uid

	im, err := meta.GetMeta(arg.GetMsgID())
	if err != nil {
		logx.Infof("Service NotFound MsgID=%d", arg.GetMsgID())
		return errcode.MetaCoderNotFound
	}

	arg2, err := im.DecodeReply(arg.Payload)
	if err != nil {
		logx.Infof("MsgID=%d req arg decode err: %v", arg.GetMsgID(), err)
		return errcode.InternalErr
	}

	return handle(pa, arg2)
}
