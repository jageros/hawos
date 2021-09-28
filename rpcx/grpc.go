/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    grpc
 * @Date:    2021/7/9 6:15 下午
 * @package: rpc
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package rpcx

import (
	"context"
	"github.com/jageros/hawox/protoc"
	"github.com/jageros/hawox/protos/pbf"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) ReqCall(ctx context.Context, arg *pbf.ReqArg) (*pbf.RespMsg, error) {
	return protoc.OnRouterRpcCall(arg)
}

func RegistryRpcServer(s *grpc.Server) {
	pbf.RegisterRouterServer(s, &server{})
}
