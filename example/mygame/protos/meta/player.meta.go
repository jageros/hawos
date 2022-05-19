// Code generated by metactl. DO NOT EDIT.
// source: player.proto

package meta

import (
	"errors"

	sess "git.hawtech.cn/jager/hawox/example/mygame/protos/meta/sess"

	pb "git.hawtech.cn/jager/hawox/example/mygame/protos/pb"	
)

//@ C2S_FETCH_CONFIG resp: Config
//------------------------------------------------------------------------------------------

var C2S_FETCH_CONFIG = &meta_C2S_FETCH_CONFIG{}

// implement IMeta

type meta_C2S_FETCH_CONFIG struct {
	handleFn func(ss sess.ISession) (resp *pb.Config, err error)
}

func (m *meta_C2S_FETCH_CONFIG) RegistryHandle(f func(ss sess.ISession) (resp *pb.Config, err error)) {
	m.handleFn = f
	registerMeta(m)
}

func (m *meta_C2S_FETCH_CONFIG) handle(ss sess.ISession, arg interface{}) (interface{}, error) {
	return m.handleFn(ss)
}

func (m *meta_C2S_FETCH_CONFIG) getMsgID() pb.MsgID {
	return pb.MsgID_C2S_FETCH_CONFIG
}

func (m *meta_C2S_FETCH_CONFIG) encodeArg(arg interface{}) ([]byte, error) {
		return nil, nil
}

func (m *meta_C2S_FETCH_CONFIG) decodeArg(data []byte) (interface{}, error) {
	return nil, nil
}

func (m *meta_C2S_FETCH_CONFIG) encodeReply(reply interface{}) ([]byte, error) {
	_reply, ok := reply.(*pb.Config)
	if !ok {
		p, ok := reply.([]byte)
		if ok {
			return p, nil
		}

		return nil, errors.New("C2S_FETCH_CONFIG_meta EncodeReply error type")
	}

	return _reply.Marshal()
}

func (m *meta_C2S_FETCH_CONFIG) decodeReply(data []byte) (interface{}, error) {
	reply := &pb.Config{}
	if err := reply.Unmarshal(data); err != nil {
		return nil, err
	} else {
		return reply, nil
	}
}

// ================== C2S_FETCH_CONFIG End ==================


//@ C2S_PLAYER_LOGIN req: LoginArg resp: LoginResp
//------------------------------------------------------------------------------------------

var C2S_PLAYER_LOGIN = &meta_C2S_PLAYER_LOGIN{}

// implement IMeta

type meta_C2S_PLAYER_LOGIN struct {
	handleFn func(ss sess.ISession, arg *pb.LoginArg) (resp *pb.LoginResp, err error)
}

func (m *meta_C2S_PLAYER_LOGIN) RegistryHandle(f func(ss sess.ISession, arg *pb.LoginArg) (resp *pb.LoginResp, err error)) {
	m.handleFn = f
	registerMeta(m)
}

func (m *meta_C2S_PLAYER_LOGIN) handle(ss sess.ISession, arg interface{}) (interface{}, error) {
	return m.handleFn(ss, arg.(*pb.LoginArg))
}

func (m *meta_C2S_PLAYER_LOGIN) getMsgID() pb.MsgID {
	return pb.MsgID_C2S_PLAYER_LOGIN
}

func (m *meta_C2S_PLAYER_LOGIN) encodeArg(arg interface{}) ([]byte, error) {
	_arg, ok := arg.(*pb.LoginArg)
	if !ok {
		p, ok := arg.([]byte)
		if ok {
			return p, nil
		}

		return nil, errors.New("C2S_PLAYER_LOGIN_meta EncodeArg error type")
	}

	return _arg.Marshal()
}

func (m *meta_C2S_PLAYER_LOGIN) decodeArg(data []byte) (interface{}, error) {
	arg := &pb.LoginArg{}
	if err := arg.Unmarshal(data); err != nil {
		return nil, err
	} else {
		return arg, nil
	}
}

func (m *meta_C2S_PLAYER_LOGIN) encodeReply(reply interface{}) ([]byte, error) {
	_reply, ok := reply.(*pb.LoginResp)
	if !ok {
		p, ok := reply.([]byte)
		if ok {
			return p, nil
		}

		return nil, errors.New("C2S_PLAYER_LOGIN_meta EncodeReply error type")
	}

	return _reply.Marshal()
}

func (m *meta_C2S_PLAYER_LOGIN) decodeReply(data []byte) (interface{}, error) {
	reply := &pb.LoginResp{}
	if err := reply.Unmarshal(data); err != nil {
		return nil, err
	} else {
		return reply, nil
	}
}

// ================== C2S_PLAYER_LOGIN End ==================


//@ C2S_PLAYER_PLAYING req: PlayingArg
//------------------------------------------------------------------------------------------

var C2S_PLAYER_PLAYING = &meta_C2S_PLAYER_PLAYING{}

// implement IMeta

type meta_C2S_PLAYER_PLAYING struct {
	handleFn func(ss sess.ISession, arg *pb.PlayingArg) (err error)
}

func (m *meta_C2S_PLAYER_PLAYING) RegistryHandle(f func(ss sess.ISession, arg *pb.PlayingArg) (err error)) {
	m.handleFn = f
	registerMeta(m)
}

func (m *meta_C2S_PLAYER_PLAYING) handle(ss sess.ISession, arg interface{}) (interface{}, error) {
	return nil, m.handleFn(ss, arg.(*pb.PlayingArg))
}

func (m *meta_C2S_PLAYER_PLAYING) getMsgID() pb.MsgID {
	return pb.MsgID_C2S_PLAYER_PLAYING
}

func (m *meta_C2S_PLAYER_PLAYING) encodeArg(arg interface{}) ([]byte, error) {
	_arg, ok := arg.(*pb.PlayingArg)
	if !ok {
		p, ok := arg.([]byte)
		if ok {
			return p, nil
		}

		return nil, errors.New("C2S_PLAYER_PLAYING_meta EncodeArg error type")
	}

	return _arg.Marshal()
}

func (m *meta_C2S_PLAYER_PLAYING) decodeArg(data []byte) (interface{}, error) {
	arg := &pb.PlayingArg{}
	if err := arg.Unmarshal(data); err != nil {
		return nil, err
	} else {
		return arg, nil
	}
}

func (m *meta_C2S_PLAYER_PLAYING) encodeReply(reply interface{}) ([]byte, error) {
	return nil, nil
}

func (m *meta_C2S_PLAYER_PLAYING) decodeReply(data []byte) (interface{}, error) {
	return nil, nil
}

// ================== C2S_PLAYER_PLAYING End ==================


//@ C2S_PLAYER_LOGOUT
//------------------------------------------------------------------------------------------

var C2S_PLAYER_LOGOUT = &meta_C2S_PLAYER_LOGOUT{}

// implement IMeta

type meta_C2S_PLAYER_LOGOUT struct {
	handleFn func(ss sess.ISession) (err error)
}

func (m *meta_C2S_PLAYER_LOGOUT) RegistryHandle(f func(ss sess.ISession) (err error)) {
	m.handleFn = f
	registerMeta(m)
}

func (m *meta_C2S_PLAYER_LOGOUT) handle(ss sess.ISession, arg interface{}) (interface{}, error) {
	return nil, m.handleFn(ss)
}

func (m *meta_C2S_PLAYER_LOGOUT) getMsgID() pb.MsgID {
	return pb.MsgID_C2S_PLAYER_LOGOUT
}

func (m *meta_C2S_PLAYER_LOGOUT) encodeArg(arg interface{}) ([]byte, error) {
		return nil, nil
}

func (m *meta_C2S_PLAYER_LOGOUT) decodeArg(data []byte) (interface{}, error) {
	return nil, nil
}

func (m *meta_C2S_PLAYER_LOGOUT) encodeReply(reply interface{}) ([]byte, error) {
	return nil, nil
}

func (m *meta_C2S_PLAYER_LOGOUT) decodeReply(data []byte) (interface{}, error) {
	return nil, nil
}

// ================== C2S_PLAYER_LOGOUT End ==================

