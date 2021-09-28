/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    rpc
 * @Date:    2021/6/8 3:59 下午
 * @package: rpc
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package rpcx

import (
	"context"
	"fmt"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/registry"
	"github.com/jageros/hawox/uuid"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	ctx    contextx.Context
	svr    *grpc.Server
	option *Option
}

type Option struct {
	ID           string
	Name         string
	Ip           string
	Port         int // 端口最大值：65535
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	CloseTimeout time.Duration
	Registrar    registry.Registrar // 服务注册接口
	Endpoint     string
}

func defaultOption() *Option {
	return &Option{
		ID:           uuid.NewAppid("server"),
		Name:         "server",
		Ip:           "127.0.0.1",
		Port:         9999,
		ReadTimeout:  time.Second * 2,
		WriteTimeout: time.Second * 2,
		CloseTimeout: time.Second * 10,
	}
}

func NewServer(ctx contextx.Context, opfs ...func(opt *Option)) *Server {
	s := &Server{
		ctx:    ctx,
		option: defaultOption(),
		svr:    grpc.NewServer(),
	}

	for _, opf := range opfs {
		opf(s.option)
	}

	s.run()

	return s
}

func (s *Server) RegistryService(registryFunc func(svr *grpc.Server)) {
	registryFunc(s.svr)
}

func (s *Server) run() {
	s.ctx.Go(func(ctx contextx.Context) error {
		addr := fmt.Sprintf("%s:%d", s.option.Ip, s.option.Port)
		li, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		err = s.register()
		if err != nil {
			return err
		}
		return s.svr.Serve(li)
	})

	s.ctx.Go(func(ctx contextx.Context) error {
		<-ctx.Done()
		err := s.deregister()
		s.svr.GracefulStop()
		if err != nil {
			return err
		}
		return ctx.Err()
	})
}

func (s *Server) buildServiceInstance() *registry.ServiceInstance {
	if s.option.Endpoint == "" {
		s.option.Endpoint = fmt.Sprintf("%s:%d", s.option.Ip, s.option.Port)
	}

	return &registry.ServiceInstance{
		ID:       s.option.ID,
		Name:     s.option.Name,
		Type:     "grpc",
		Version:  "1.0",
		Endpoint: s.option.Endpoint,
	}
}

func (s *Server) register() error {
	if s.option.Registrar != nil {
		return s.option.Registrar.Register(s.ctx, s.buildServiceInstance())
	}
	return nil
}

func (s *Server) deregister() error {
	if s.option.Registrar != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.option.CloseTimeout)
		defer cancel()
		return s.option.Registrar.Deregister(ctx, s.buildServiceInstance())
	}
	return nil
}
