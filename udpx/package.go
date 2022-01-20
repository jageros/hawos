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

type MsgType int8

const (
	TextMessage   MsgType = 1
	BinaryMessage MsgType = 2
	PingMessage   MsgType = 3
	PongMessage   MsgType = 4
)

type Package struct {
	Type    MsgType
	Payload []byte
}

func (p *Package) Marshal() []byte {
	buff := make([]byte, len(p.Payload)+1)
	buff[0] = byte(p.Type)
	copy(buff[1:], p.Payload)
	return buff
}

func (p *Package) UnMarshal(body []byte) {
	p.Type = MsgType(body[0])
	p.Payload = make([]byte, len(body)-1)
	copy(p.Payload, body[1:])
}
