package registry

import (
	"context"
	"fmt"
	"github.com/jageros/hawox/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

const Prefix = "/microservices"

var (
	_ Registrar = &Registry{}
	_ Discovery = &Registry{}
)

// Option is etcd registry option.
type Option func(o *options)

type options struct {
	ctx       context.Context
	namespace string
	ttl       time.Duration
}

// Context with registry context.
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Namespace with registry namespance.
func Namespace(ns string) Option {
	return func(o *options) { o.namespace = ns }
}

// RegisterTTL with register ttl.
func RegisterTTL(ttl time.Duration) Option {
	return func(o *options) { o.ttl = ttl }
}

// Registry is etcd registry.
type Registry struct {
	opts   *options
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

// New creates etcd registry
func New(client *clientv3.Client, opts ...Option) (r *Registry) {
	options := &options{
		ctx:       context.Background(),
		namespace: Prefix,
		ttl:       time.Second * 15,
	}
	for _, o := range opts {
		o(options)
	}
	return &Registry{
		opts:   options,
		client: client,
		kv:     clientv3.NewKV(client),
	}
}

// Register the registration.
func (r *Registry) Register(ctx context.Context, service *ServiceInstance) error {
	key := service.Key(r.opts.namespace)
	value, err := marshal(service)
	if err != nil {
		return err
	}
	if r.lease != nil {
		r.lease.Close()
	}
	r.lease = clientv3.NewLease(r.client)
	grant, err := r.lease.Grant(ctx, int64(r.opts.ttl.Seconds()))
	if err != nil {
		return err
	}
	_, err = r.client.Put(ctx, key, value, clientv3.WithLease(grant.ID))
	if err != nil {
		return err
	}
	hb, err := r.client.KeepAlive(ctx, grant.ID)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case _, ok := <-hb:
				if !ok {
					return
				}
			case <-r.opts.ctx.Done():
				return
			}
		}
	}()
	logx.Infof("register Service Type=%s Name=%s ID=%s", service.Type, service.Name, service.ID)
	return nil
}

// Deregister the registration.
func (r *Registry) Deregister(ctx context.Context, service *ServiceInstance) error {
	defer func() {
		if r.lease != nil {
			r.lease.Close()
		}
	}()
	key := service.Key(r.opts.namespace)
	_, err := r.client.Delete(ctx, key)
	logx.Infof("deregister service Type=%s Name=%s ID=%s", service.Type, service.Name, service.ID)
	return err
}

// GetService return the service instances in memory according to the service name.
func (r *Registry) GetService(ctx context.Context, name string) ([]*ServiceInstance, error) {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, name)
	resp, err := r.kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var items []*ServiceInstance
	for _, kv := range resp.Kvs {
		si, err := unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		items = append(items, si)
	}
	return items, nil
}

// Watch creates a watcher according to the service name.
func (r *Registry) Watch(ctx context.Context, name string) (Watcher, error) {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, name)
	return newWatcher(ctx, key, r.client), nil
}