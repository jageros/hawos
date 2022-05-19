/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    publish
 * @Date:    2021/7/2 3:51 下午
 * @package: nsq
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package nsq

import (
	"context"
	"fmt"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/httpc"
	"github.com/jageros/hawox/logx"
	"github.com/nsqio/go-nsq"
	"math/rand"
	"strings"
	"sync"
)

type Producer struct {
	ctx contextx.Context
	opt *Config
	pd  *nsq.Producer
	cfg *nsq.Config
	clk *sync.Mutex
}

func (p *Producer) getNodeAddr() (string, error) {
	addrs := strings.Split(p.opt.Addrs, ";")
	idx := rand.Intn(len(addrs))
	url := fmt.Sprintf("http://%s/nodes", addrs[idx])
	resp, err := httpc.RequestReturnMap(httpc.GET, url, httpc.FORM, nil, nil)
	if err != nil {
		return "", err
	}
	pds := resp["producers"].([]interface{})
	pdn := len(pds)
	if pdn <= 0 {
		return "", errcode.New(101, "无可用NSQ节点")
	}
	idx = rand.Intn(len(pds))
	pd := pds[idx].(map[string]interface{})
	addr := fmt.Sprintf("%v:%v", pd["broadcast_address"], pd["tcp_port"])
	logx.Debug().Str("addr", addr).Msg("NsqAddr")
	return addr, nil
}

func (p *Producer) connectToNsqd() error {
	p.clk.Lock()
	defer p.clk.Unlock()

	if p.pd != nil {
		err := p.pd.Ping()
		if err == nil {
			return nil
		}
		p.pd.Stop()
	}

	addr, err := p.getNodeAddr()
	if err != nil {
		return err
	}
	pd, err := nsq.NewProducer(addr, p.cfg)
	if err != nil {
		return err
	}

	p.pd = pd
	return nil
}

func (p *Producer) Push(data []byte) error {
	err := p.pd.Publish(p.opt.Topic, data)
	if err != nil {
		err = p.connectToNsqd()
		if err != nil {
			return err
		}
		err = p.pd.Publish(p.opt.Topic, data)
	}
	return err
}

func NewProducer(g contextx.Context, opfs ...func(cfg *Config)) (*Producer, error) {
	p := &Producer{
		ctx: g,
		opt: defaultConfig(),
		clk: &sync.Mutex{},
	}

	for _, opf := range opfs {
		opf(p.opt)
	}

	p.cfg = nsq.NewConfig()
	err := p.connectToNsqd()
	p.run()
	return p, err
}

func (p *Producer) run() {
	p.ctx.Go(func(ctx context.Context) error {
		<-ctx.Done()
		p.pd.Stop()
		return ctx.Err()
	})
}
