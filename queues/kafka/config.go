/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    kafka
 * @Date:    2021/8/20 3:54 下午
 * @package: kafka
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package kafka

import (
	"git.hawtech.cn/jager/hawox/uuid"
	"time"
)

type Config struct {
	Addrs    string
	Topic    string
	GroupId  string
	WarnTime time.Duration
}

func defaultConfig() *Config {
	return &Config{
		Addrs:    "",
		Topic:    "queue_msg",
		GroupId:  uuid.NewGlobalId(),
		WarnTime: time.Millisecond * 200,
	}
}
