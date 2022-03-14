/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    nsq
 * @Date:    2021/7/2 3:28 下午
 * @package: nsq
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package nsq

import (
	"context"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/logx"
	"github.com/nsqio/go-nsq"
	"strings"
	"time"
)

type Consumer struct {
	ctx     contextx.Context
	cfg     *Config
	csr     *nsq.Consumer
	handler func(data []byte)
	csrCnt  int64
}

func NewConsumer(ctx contextx.Context, opfs ...func(cfg *Config)) (*Consumer, error) {
	c := &Consumer{
		ctx: ctx,
		cfg: defaultConfig(),
		handler: func(data []byte) {
			logx.Warnf("Nsq Consumer receive msg, but handler not set!")
		},
	}

	for _, opf := range opfs {
		opf(c.cfg)
	}

	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second
	csr, err := nsq.NewConsumer(c.cfg.Topic, c.cfg.Channel, cfg)
	if err != nil {
		return nil, err
	}

	csr.SetLogger(nil, 0)
	csr.AddHandler(c)

	c.csr = csr

	addr := strings.Split(c.cfg.Addrs, ";")
	err = c.csr.ConnectToNSQLookupds(addr)

	c.run()

	return c, err
}

func (c *Consumer) OnMessageHandler(f func(data []byte)) {
	c.handler = f
}

func (c *Consumer) HandleMessage(msg *nsq.Message) error {
	start := time.Now()

	c.handler(msg.Body)

	take := time.Now().Sub(start)
	if take > c.cfg.WarnTime {
		logx.Warnf("Nsq Consumer Msg take: %s", take.String())
	}
	return nil
}

func (c *Consumer) run() {
	c.ctx.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			if c.csr != nil {
				c.csr.Stop()
			}
			return ctx.Err()
		case i := <-c.csr.StopChan:
			logx.Infof("Consumer StopChan=%d", i)
			return nil
		}
	})
}
