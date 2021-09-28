package ws

import (
	"github.com/jageros/hawox/protos/pbf"
)

func interceptor(uid string, arg *pbf.Request) (*pbf.Response, bool) {
	return nil, false
}
