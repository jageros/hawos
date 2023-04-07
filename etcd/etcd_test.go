package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	cli, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	kv := clientv3.NewKV(cli)
	kv.Put(ctx, "name", "jager")
	rsp, err := kv.Get(ctx, "name")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(rsp.Kvs[0].Value))
}
