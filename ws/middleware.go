package ws

import (
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/protos/pbf"
	"gopkg.in/olahol/melody.v1"
)

func frontendMiddleware(uid string, arg *pbf.Request, session *melody.Session) errcode.IErr {
	return nil
}

func backendMiddleware(uid string, arg *pbf.Response, session *melody.Session) errcode.IErr {
	return nil
}
