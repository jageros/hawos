/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    config
 * @Date:    2021/8/20 5:20 下午
 * @package: nsq
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package nsq

import (
	"github.com/jager/hawox/uuid"
	"time"
)

type Config struct {
	Addrs    string
	Topic    string
	Channel  string
	WarnTime time.Duration
}

func defaultConfig() *Config {
	return &Config{
		Addrs:    "127.0.0.1:4161",
		Topic:    "queue_msg",
		Channel:  uuid.NewGlobalId(),
		WarnTime: time.Millisecond * 200,
	}
}
