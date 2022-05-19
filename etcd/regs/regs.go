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

package regs

import (
	"context"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	prefix = "global"
	ttl    = time.Second * 15
)

type Options interface {
	SetTTL(t time.Duration)
	SetNamespace(namespace string)
}

// Registry is etcd registry.
type Registry struct {
	ctx       contextx.Context
	ttl       time.Duration
	namespace string

	// etcd
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	// watcher
	watchChan clientv3.WatchChan
	watcher   clientv3.Watcher

	values      map[string]string
	keys        []string
	keyIndex    int64
	keyIndexMux sync.Mutex
	rw          sync.RWMutex
}

func (r *Registry) key(key string) string {
	if strings.HasPrefix(key, "/") {
		return r.namespace + key
	} else {
		return r.namespace + "/" + key
	}
}

func (r *Registry) SetNamespace(namespace string) {
	r.namespace = namespace
}

func (r *Registry) SetTTL(t time.Duration) {
	r.ttl = t
}

func New(ctx contextx.Context, client *clientv3.Client, opts ...func(opt Options)) (r *Registry) {
	r = &Registry{
		ctx:       ctx,
		ttl:       ttl,
		namespace: prefix,
		client:    client,
		kv:        clientv3.NewKV(client),
		watcher:   clientv3.NewWatcher(client),
		values:    map[string]string{},
	}

	for _, o := range opts {
		o(r)
	}

	r.update()
	r.watch()
	return
}

func (r *Registry) set(key, value string) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.values[key] = value
	var keys []string
	for key := range r.values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	r.keys = keys
}

func (r *Registry) del(key string) {
	r.rw.Lock()
	defer r.rw.Unlock()
	delete(r.values, key)
	var keys []string
	for key := range r.values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	r.keys = keys
}

func (r *Registry) update() {
	r.ctx.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			resp, err := r.kv.Get(r.ctx, r.namespace, clientv3.WithPrefix())
			if err != nil {
				return err
			}

			r.rw.Lock()
			defer r.rw.Unlock()
			r.values = map[string]string{}
			for _, kv := range resp.Kvs {
				r.values[string(kv.Key)] = string(kv.Value)
			}
			var keys []string
			for key := range r.values {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			r.keys = keys
			return nil
		}
	})
}

func (r *Registry) watch() {
	r.ctx.Go(func(ctx context.Context) error {
		r.watchChan = r.watcher.Watch(ctx, r.namespace, clientv3.WithPrefix(), clientv3.WithRev(0))
		err := r.watcher.RequestProgress(ctx)
		if err != nil {
			return err
		}

		for {
			select {
			case <-ctx.Done():
				err := r.watcher.Close()
				if err != nil {
					return err
				}
				return ctx.Err()

			case wc := <-r.watchChan:
				for _, ev := range wc.Events {
					switch ev.Type {
					case clientv3.EventTypePut:
						r.set(string(ev.Kv.Key), string(ev.Kv.Value))
					case clientv3.EventTypeDelete:
						r.del(string(ev.Kv.Key))
					}
					logx.Info().Str("type", ev.Type.String()).Str(string(ev.Kv.Key), string(ev.Kv.Value)).Msg("RegistryEvent")
				}
			}
		}
	})
}

func (r *Registry) Register(key, value string) error {
	key = r.key(key)
	if r.lease != nil {
		r.lease.Close()
	}
	// 创建租约
	r.lease = clientv3.NewLease(r.client)
	grant, err := r.lease.Grant(r.ctx, int64(r.ttl.Seconds()))
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
	//logx.Infof("Put etcd key=%s value=%s", key, value)
	return nil
}

func (r *Registry) Deregister(key string) error {
	key = r.key(key)
	ctx, cancel := context.WithTimeout(context.Background(), r.ttl)
	defer cancel()
	_, err := r.client.Delete(ctx, key)
	if err != nil {
		return err
	} else {
		//logx.Infof("Del etcd key=%s", key)
	}
	if r.lease != nil {
		err = r.lease.Close()
	}
	return err
}

func (r *Registry) RegisterWithAutoDeregister(key, value string) error {
	err := r.Register(key, value)
	if err != nil {
		return err
	}
	r.ctx.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return r.Deregister(key)
	})
	return nil
}

func (r *Registry) Get(key string) string {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.values[r.key(key)]
}

func (r *Registry) GetWithPrefix(prefix string) []string {
	prefix = r.key(prefix)
	r.rw.RLock()
	defer r.rw.RUnlock()
	var values []string
	for key, val := range r.values {
		if strings.HasPrefix(key, prefix) {
			values = append(values, val)
		}
	}
	return values
}

func (r *Registry) GetWithPrefixRandom(prefix string) string {
	prefix = r.key(prefix)
	r.rw.RLock()
	defer r.rw.RUnlock()
	var values []string
	for key, val := range r.values {
		if strings.HasPrefix(key, prefix) {
			values = append(values, val)
		}
	}
	return values[rand.Intn(len(values))]
}

func (r *Registry) GetWithPrefixRoundRobin(prefix string) string {
	prefix = r.key(prefix)
	r.rw.RLock()
	defer r.rw.RUnlock()
	r.keyIndexMux.Lock()
	index := r.keyIndex
	if int(r.keyIndex) >= len(r.keys) {
		r.keyIndex = 1
		index = 0
	} else {
		r.keyIndex++
	}
	r.keyIndexMux.Unlock()

	return r.values[r.keys[index]]
}

func (r *Registry) RangeKV(prefix string, f func(key, value string)) {
	prefix = r.key(prefix)
	r.rw.RLock()
	defer r.rw.RUnlock()
	for key, val := range r.values {
		if strings.HasPrefix(key, prefix) {
			f(key, val)
		}
	}
}

// ====== global ========

var regs *Registry

func InitGlobalRegs(ctx contextx.Context, c *clientv3.Client, opfs ...func(opt Options)) {
	regs = New(ctx, c, opfs...)
}
func Register(key, value string) error {
	return regs.Register(key, value)
}

func Deregister(key string) error {
	return regs.Deregister(key)
}

func Get(key string) string {
	return regs.Get(key)
}

func GetWithPrefix(prefix string) []string {
	return regs.GetWithPrefix(prefix)
}

func RangKV(prefix string, f func(key, value string)) {
	regs.RangeKV(prefix, f)
}
