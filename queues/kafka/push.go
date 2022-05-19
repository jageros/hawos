/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    kafka
 * @Date:    2021/7/5 9:52 上午
 * @package: kafka
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package kafka

import (
	"context"
	"git.hawtech.cn/jager/hawox/contextx"
	"git.hawtech.cn/jager/hawox/logx"
	"github.com/Shopify/sarama"
	"strings"
)

type Producer struct {
	ctx contextx.Context
	pd  sarama.AsyncProducer
	cfg *Config
}

func NewProducer(g contextx.Context, opfs ...func(cfg *Config)) (*Producer, error) {
	pd := &Producer{
		ctx: g,
		cfg: defaultConfig(),
	}

	for _, opf := range opfs {
		opf(pd.cfg)
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V0_11_0_2
	addrs := strings.Split(pd.cfg.Addrs, ";")

	producer, err := sarama.NewAsyncProducer(addrs, config)
	if err != nil {
		return nil, err
	}
	pd.pd = producer

	pd.run()

	return pd, nil
}

func (p *Producer) Push(data []byte) error {
	msg_ := &sarama.ProducerMessage{
		Topic: p.cfg.Topic,
		Value: sarama.ByteEncoder(data),
	}
	p.ctx.Go(func(ctx context.Context) error {
		select {
		case p.pd.Input() <- msg_:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	return nil
}

func (p *Producer) run() {
	p.ctx.Go(func(ctx context.Context) error {
		var offset int64 = -1
		for {
			select {
			case <-ctx.Done():
				p.pd.AsyncClose()
				return ctx.Err()

			case errMsg := <-p.pd.Errors():
				if offset != errMsg.Msg.Offset {
					p.pd.Input() <- errMsg.Msg
					offset = errMsg.Msg.Offset
				}
				logx.Err(errMsg).Msg("Kafka Handle Err")

			case msg := <-p.pd.Successes():
				offset = msg.Offset
			}
		}
	})
}
