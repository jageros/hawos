/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    package
 * @Date:    2022/1/20 11:33 上午
 * @package: udpx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package udpx

import (
	"encoding/binary"
)

type MsgType int8

const (
	TextMessage   MsgType = 1
	BinaryMessage MsgType = 2
)

type Package struct {
	Type    MsgType
	Payload []byte
}

func (p *Package) Marshal() []byte {
	payloadLen := len(p.Payload)
	buff := make([]byte, payloadLen+5)
	buff[0] = byte(p.Type)
	binary.BigEndian.PutUint32(buff[1:], uint32(payloadLen))
	if len(p.Payload) > 0 {
		copy(buff[5:], p.Payload)
	}
	return buff
}

func (p *Package) Unmarshal(body []byte) {
	bodyLen := len(body)
	if bodyLen >= 1 {
		p.Type = MsgType(body[0])
	}
	var payloadLen = 0
	if bodyLen >= 5 {
		payloadLen = int(binary.BigEndian.Uint32(body[1:5]))
	}
	if bodyLen > 5 {
		p.Payload = make([]byte, payloadLen)
		copy(p.Payload, body[5:])
	}
}
