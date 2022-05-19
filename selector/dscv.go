/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    dscv
 * @Date:    2021/7/9 3:12 下午
 * @package: dscv
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package selector

import (
	"context"
	"fmt"
	"github.com/jager/hawox/contextx"
	"github.com/jager/hawox/logx"
	"github.com/jager/hawox/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

const Prefix = "/pbserver"

// Option is etcd registry option.
type Option func(o *Options)

type Options struct {
	Id        string
	Namespace string
	Name      string
	Ttl       time.Duration
}

// Namespace with registry namespance.
func Namespace(ns string) Option {
	return func(o *Options) { o.Namespace = ns }
}

func Name(id, name string) Option {
	return func(o *Options) {
		o.Id = id
		o.Name = name
	}
}

// RegisterTTL with Register TTL.
func RegisterTTL(ttl time.Duration) Option {
	return func(o *Options) { o.Ttl = ttl }
}

// Registry is etcd registry.
type Registry struct {
	ctx    contextx.Context
	opts   *Options
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	// watcher
	watchChan clientv3.WatchChan
	watcher   clientv3.Watcher
}

func (r *Registry) run() {
	r.ctx.Go(func(ctx context.Context) error {
		r.watchChan = r.watcher.Watch(ctx, r.opts.Namespace, clientv3.WithPrefix(), clientv3.WithRev(0))
		err := r.watcher.RequestProgress(ctx)
		if err != nil {
			return err
		}

		for {
			select {
			case <-ctx.Done():
				r.deregister()
				err := r.watcher.Close()
				if err != nil {
					return err
				}
				return ctx.Err()

			case <-r.watchChan:
				resp, err := r.kv.Get(ctx, r.opts.Namespace, clientv3.WithPrefix())
				if err != nil {
					return err
				}
				var items []*metaData
				for _, kv := range resp.Kvs {
					si, err := unmarshal(kv.Value)
					if err != nil {
						return err
					}
					items = append(items, si)
				}
				Update(items)
			}
		}
	})
}

func Initialize(ctx contextx.Context, client *clientv3.Client, opts ...Option) (r *Registry) {
	options := &Options{
		Id:        uuid.NewServiceId(),
		Name:      "selector",
		Namespace: Prefix,
		Ttl:       time.Second * 15,
	}
	for _, o := range opts {
		o(options)
	}

	r = &Registry{
		ctx:     ctx,
		opts:    options,
		client:  client,
		kv:      clientv3.NewKV(client),
		watcher: clientv3.NewWatcher(client),
	}
	r.run()
	return
}

func (r *Registry) Register(msgIds []int32) error {
	if len(msgIds) <= 0 {
		return nil
	}
	md := &metaData{
		ID:     r.opts.Id,
		Name:   r.opts.Name,
		MsgIds: msgIds,
	}
	key := fmt.Sprintf("%s/%s/%s", r.opts.Namespace, r.opts.Name, md.ID)
	value, err := marshal(md)
	if err != nil {
		return err
	}
	if r.lease != nil {
		r.lease.Close()
	}
	// 创建租约
	r.lease = clientv3.NewLease(r.client)
	grant, err := r.lease.Grant(r.ctx, int64(r.opts.Ttl.Seconds()))
	if err != nil {
		return err
	}

	// 键值注册
	_, err = r.client.Put(r.ctx, key, value, clientv3.WithLease(grant.ID))
	if err != nil {
		return err
	}

	// 保持心跳
	hb, err := r.client.KeepAlive(r.ctx, grant.ID)
	if err != nil {
		return err
	}
	r.ctx.Go(func(ctx context.Context) error {
		for {
			select {
			case _, ok := <-hb:
				if !ok {
					return nil
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
	logx.Info().Str("Namespace", r.opts.Namespace).Str("Name", md.Name).Str("ID", md.ID).Ints32("MsgIds", md.MsgIds).Msg("Register Service")
	return nil
}

func (r *Registry) deregister() {
	defer func() {
		if r.lease != nil {
			r.lease.Close()
		}
	}()
	key := fmt.Sprintf("%s/%s/%s", r.opts.Namespace, r.opts.Name, r.opts.Id)
	ctx, cancel := context.WithTimeout(r.ctx, time.Second*10)
	defer cancel()
	_, err := r.client.Delete(ctx, key)
	logx.Err(err).Str("Namespace", r.opts.Namespace).Str("Name", r.opts.Name).Str("ID", r.opts.Id).Msg("pdserver deregister")
}
