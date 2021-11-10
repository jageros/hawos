/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    packet
 * @Date:    2021/11/4 5:53 下午
 * @package: packet
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package packet

import (
	errCode "github.com/jageros/hawox/errcode"
	"sync"
)

var (
	packetTooBigErr   = errCode.New(5001, "packetTooBigErr")
	unknownMsgTypeErr = errCode.New(5003, "unknownMsgType")
)

var packetPool *sync.Pool

func init() {
	packetPool = &sync.Pool{
		New: func() interface{} {
			return &Packet{}
		},
	}
}

const (
	SeqSize     = 4
	MsgIDSize   = 4
	MsgTypeSize = 1
	ErrcodeSize = 4
	PktSizeSize = 4
)

type MessageType = byte

const (
	MsgReq MessageType = 1 + iota // need reply
	MsgReply
	MsgPush // no reply
	MsgErr  // error reply
	MsgPing
	MsgPong
)

var (
	PingPacket = &Packet{
		msgType: MsgPing,
	}
	PongPacket = &Packet{
		msgType: MsgPong,
	}
)

type Packet struct {
	seq     uint32
	msgID   int32
	msgType MessageType
	errcode int32
	payload []byte
}

func GetPacket() *Packet {
	ip := packetPool.Get()
	pk := ip.(*Packet)
	pk.Reset()
	return pk
}

func PutPacket(pk *Packet) {
	if pk == PingPacket || pk == PongPacket {
		return
	}
	packetPool.Put(pk)
}

func (p *Packet) Reset() {
	p.seq = 0
	p.msgID = 0
	p.msgType = 0
	p.errcode = 0
	p.payload = nil
}

func (p *Packet) GetSeq() uint32 {
	return p.seq
}

func (p *Packet) GetMsgID() int32 {
	return p.msgID
}

func (p *Packet) GetMsgType() MessageType {
	return p.msgType
}

func (p *Packet) GetErrcode() int32 {
	return p.errcode
}

func (p *Packet) GetPayload() []byte {
	return p.payload
}

func (p *Packet) SetPayload(payload []byte) {
	p.payload = payload
}
