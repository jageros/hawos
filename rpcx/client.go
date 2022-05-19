/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    client
 * @Date:    2021/6/9 4:13 下午
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
	"google.golang.org/grpc/resolver"
	"sync"
	"time"

	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/registry"
	"github.com/jageros/hawox/resolver/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/connectivity"
)

type RpcFn func(cc *grpc.ClientConn)

var cli *client

type Config interface {
	SetCallTimeout(callTimeout time.Duration)
}

type client struct {
	ctx         contextx.Context
	d           registry.Discovery
	conns       map[string]*grpc.ClientConn // map[serverName]conn
	rw          *sync.RWMutex
	callTimeout time.Duration
	builder     resolver.Builder
}

func (c *client) SetCallTimeout(callTimeout time.Duration) {
	c.callTimeout = callTimeout
}

func newClient(ctx contextx.Context, d registry.Discovery, opfs ...func(cfg Config)) *client {
	cli_ := &client{
		ctx:         ctx,
		d:           d,
		conns:       map[string]*grpc.ClientConn{},
		rw:          &sync.RWMutex{},
		callTimeout: time.Second * 3,
	}

	for _, opf := range opfs {
		opf(cli_)
	}

	cli_.builder = discovery.NewBuilder(cli_.ctx, cli_.d)

	cli_.stopMonitor()

	return cli_
}

func InitClient(ctx contextx.Context, d registry.Discovery, opfs ...func(cfctx Config)) {
	cli = newClient(ctx, d, opfs...)
}

func (c *client) getConnByName(name string) (*grpc.ClientConn, error) {
	c.rw.RLock()
	conn, ok := c.conns[name]
	c.rw.RUnlock()
	if ok {
		return conn, nil
	}

	c.rw.Lock()
	defer c.rw.Unlock()

	// 为避免别的协程已经创建，进行检测后再创建
	conn, ok = c.conns[name]
	if ok {
		return conn, nil
	}

	target := fmt.Sprintf("%s:///%s", discovery.Name, name)

	ctx, cancel := context.WithTimeout(c.ctx, c.callTimeout)
	defer cancel()

	cc, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, roundrobin.Name)), // This sets the initial balancinctx policy.
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithResolvers(c.builder),
	)
	if err != nil {
		return nil, err
	}
	c.conns[name] = cc

	return cc, nil
}

//
//func (c *client) getConnByTarget(target string) (*grpc.ClientConn, error) {
//	c.rw.RLock()
//	conn, ok := c.conns[target]
//	c.rw.RUnlock()
//	if ok {
//		return conn, nil
//	}
//
//	c.rw.Lock()
//	defer c.rw.Unlock()
//
//	// 为避免别的协程已经创建，进行检测后再创建
//	conn, ok = c.conns[target]
//	if ok {
//		return conn, nil
//	}
//
//	ctx, cancel := context.WithTimeout(c.ctx, c.callTimeout)
//	defer cancel()
//
//	cc, err := grpc.DialContext(
//		ctx,
//		target,
//		grpc.WithInsecure(),
//		grpc.WithBlock(),
//	)
//	if err != nil {
//		return nil, err
//	}
//	c.conns[target] = cc
//
//	return cc, nil
//}

func (c *client) stopMonitor() {
	c.ctx.Go(func(ctx context.Context) error {
		<-ctx.Done()
		for _, cc := range c.conns {
			cc.Close()
		}
		return ctx.Err()
	})
}

// --- call ---

func (c *client) call(cc *grpc.ClientConn, rpcFn RpcFn) errcode.IErr {
	state := cc.GetState()
	if state != connectivity.Ready {
		ctx, cancel := context.WithTimeout(c.ctx, c.callTimeout)
		cc.WaitForStateChange(ctx, state)
		cancel()
	}
	if cc.GetState() == connectivity.Ready {
		rpcFn(cc)
		return nil
	} else {
		logx.Error().Str("target", cc.Target()).Msg("Service Conn NotReady!")
		errMsg := fmt.Sprintf("%s Service Conn NotReady!", cc.Target())
		return errcode.New(-22, errMsg)
	}
}

func CallByName(serviceName string, rpcFn RpcFn) errcode.IErr {
	cc, err := cli.getConnByName(serviceName)
	if err != nil {
		logx.Err(err).Msg("CallByName getConnByName")
		return errcode.WithErrcode(-11, err)
	}

	return cli.call(cc, rpcFn)
}

//func CallByTarget(target string, rpcFn RpcFn) errcode.IErr {
//	cc, err := cli.getConnByTarget(target)
//	if err != nil {
//		logx.Errorf("CallByTarget getConnByTarget err: %v", err)
//		return errcode.WithErrcode(-11, err)
//	}
//
//	return cli.call(cc, rpcFn)
//}
