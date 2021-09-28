/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    queues
 * @Date:    2021/7/21 10:45 上午
 * @package: queues
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package queues

import (
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/protos/pbf"
	"github.com/jageros/hawox/queues/kafka"
	"github.com/jageros/hawox/queues/nsq"
	"time"
)

type Config struct {
	Type     string
	Addrs    string
	Topic    string
	GroupId  string
	WarnTime time.Duration
}

func (c *Config) kafka() func(cfg *kafka.Config) {
	return func(cfg *kafka.Config) {
		if c.Addrs != "" {
			cfg.Addrs = c.Addrs
		}
		if c.Topic != "" {
			cfg.Topic = c.Topic
		}
		if c.GroupId != "" {
			cfg.GroupId = c.GroupId
		}
		if c.WarnTime > 0 {
			cfg.WarnTime = c.WarnTime
		}
	}
}

func (c *Config) nsq() func(cfg *nsq.Config) {
	return func(cfg *nsq.Config) {
		if c.Addrs != "" {
			cfg.Addrs = c.Addrs
		}
		if c.Topic != "" {
			cfg.Topic = c.Topic
		}
		if c.GroupId != "" {
			cfg.Channel = c.GroupId
		}
		if c.WarnTime > 0 {
			cfg.WarnTime = c.WarnTime
		}
	}
}

func defaultConfig() *Config {
	return &Config{
		Type: "kafka",
	}
}

type IProducer interface {
	PushProtoMsg(msgId int32, arg interface{}, target *pbf.Target) error
	Push(msg *pbf.QueueMsg) error
}

type IConsumer interface {
	OnMessageHandler(f func(msg *pbf.QueueMsg))
}

func NewProducer(ctx contextx.Context, opfs ...func(cfg *Config)) (IProducer, error) {
	op := defaultConfig()
	for _, opf := range opfs {
		opf(op)
	}
	switch op.Type {
	case "nsq":
		return nsq.NewProducer(ctx, op.nsq())
	case "kafka":
		return kafka.NewProducer(ctx, op.kafka())
	default:
		return nil, errcode.New(-1, "未知队列类型")
	}
}

func NewConsumer(ctx contextx.Context, opfs ...func(cfg *Config)) (IConsumer, error) {
	op := defaultConfig()
	for _, opf := range opfs {
		opf(op)
	}
	switch op.Type {
	case "nsq":
		return nsq.NewConsumer(ctx, op.nsq())
	case "kafka":
		return kafka.NewConsumer(ctx, op.kafka())
	default:
		return nil, errcode.New(-1, "未知队列类型")
	}
}
