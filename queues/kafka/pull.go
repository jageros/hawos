/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    pull
 * @Date:    2021/7/5 6:49 下午
 * @package: kafka
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/protos/pbf"
	"strings"
	"time"
)

type Consumer struct {
	ctx     contextx.Context
	cg      sarama.ConsumerGroup
	cfg     *Config
	handler func(msg *pbf.QueueMsg)
}

func (c *Consumer) OnMessageHandler(f func(msg *pbf.QueueMsg)) {
	c.handler = f
}

func NewConsumer(ctx contextx.Context, opfs ...func(cfg *Config)) (*Consumer, error) {
	csr := &Consumer{
		ctx: ctx,
		cfg: defaultConfig(),
		handler: func(msg *pbf.QueueMsg) {
			logx.Warnf("Kafka Consumer receive msg, but handler not set")
		},
	}

	for _, opf := range opfs {
		opf(csr.cfg)
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	config.Version = sarama.V0_11_0_2
	addrs := strings.Split(csr.cfg.Addrs, ";")
	cli, err := sarama.NewClient(addrs, config)
	if err != nil {
		return nil, err
	}

	cg, err := sarama.NewConsumerGroupFromClient(csr.cfg.GroupId, cli)

	if err != nil {
		return nil, err
	}

	csr.cg = cg

	csr.run()

	return csr, nil
}

func (c *Consumer) run() {
	c.ctx.Go(func(ctx contextx.Context) error {
		for {
			select {
			case <-ctx.Done():
				if c.cg != nil {
					err := c.cg.Close()
					if err != nil {
						return err
					}
				}
				return ctx.Err()
			default:
				err := c.cg.Consume(ctx, []string{c.cfg.Topic}, c)
				if err != nil {
					logx.Errorf("Consume err: %v", err)
					return err
				}
			}
		}
	})

}

func (c *Consumer) Setup(assignment sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(assignment sarama.ConsumerGroupSession) error { return nil }
func (c *Consumer) ConsumeClaim(assignment sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		start := time.Now()
		if msg == nil {
			logx.Infof("kafka ConsumeClaim recv msg=nil")
			continue
		}
		kmsg := &pbf.QueueMsg{}
		err := kmsg.Unmarshal(msg.Value)
		if err != nil {
			logx.Errorf("kafka Unmarshal msg err=%v", err)
			continue
		}

		c.handler(kmsg)

		assignment.MarkMessage(msg, "") // 确认消息
		take := time.Now().Sub(start)
		if take >= c.cfg.WarnTime {
			logx.Warnf("kafka consume msg take %s", take.String())
		}
	}
	return nil
}
