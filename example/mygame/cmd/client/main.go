/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    main
 * @Date:    2022/3/11 10:07
 * @package: client
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package main

import (
	"github.com/jager/hawox/contextx"
	"github.com/jager/hawox/example/mygame/protos/pb"
	"github.com/jager/hawox/wsc"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := contextx.Default()
	defer cancel()
	wc := wsc.New(ctx)
	ss, err := wc.ConnectWithHeader("ws://127.0.0.1:10088/ws/gate", http.Header{"uid": {"1001"}})
	if err != nil {
		log.Fatal(err)
	}

	wc.HandleMessageBinary(func(ss *wsc.Session, bytes []byte) {
		resp := &pb.PkgMsg{}
		err = resp.Unmarshal(bytes)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(resp.String())
		}
	})

	tk := time.NewTicker(time.Second)
	var msgid int32 = 100
	msg := &pb.PkgMsg{
		Type: pb.MsgType_Req,
	}
loop:
	for {
		select {
		case <-tk.C:
			msg.Msgid = pb.MsgID(msgid)
			if msgid == 101 {
				req := &pb.LoginArg{
					Account:  "lhj168os@gmail.com",
					Password: "xxxxxxxx",
				}
				msg.Payload, _ = req.Marshal()
			}
			if msgid == 102 {
				req := &pb.PlayingArg{
					Vigor: 90,
					Angle: 108,
				}
				msg.Payload, _ = req.Marshal()
			}
			data, _ := msg.Marshal()
			err = ss.WriteBinary(data)
			if err != nil {
				log.Println(err)
			}

		case <-ctx.Done():
			break loop
		}
		msgid++
		if msgid > 103 {
			msgid = 100
		}
	}
}
