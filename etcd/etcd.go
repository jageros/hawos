/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    etcd
 * @Date:    2021/6/11 6:02 下午
 * @package: etcd
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type OpFn func(config *clientv3.Config)

func New(opfs ...OpFn) (*clientv3.Client, error) {
	config := &clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	}
	for _, opf := range opfs {
		opf(config)
	}
	cli, err := clientv3.New(*config)
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), config.DialTimeout)
	defer cancel()
	_, err = cli.Status(timeoutCtx, config.Endpoints[0])
	if err != nil {
		return nil, err
	}
	return cli, nil
}
