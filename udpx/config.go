/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    config
 * @Date:    2022/1/18 5:34 下午
 * @package: udpx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package udpx

import (
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
	"time"
)

type Config struct {
	ListenIp     string // 监听IP
	Port         int    // 监听端口
	Mode         string
	MaxByte      int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	CloseTimeout time.Duration
	RateTime     time.Duration

	OnMsgHandle func(addr *net.UDPAddr, data []byte) (resp []byte)
}

func (s *Config) addr() string {
	return s.ListenIp + ":" + strconv.Itoa(s.Port)
}

func defaultServer() *Config {
	return &Config{
		ListenIp:     "",
		Port:         8888,
		Mode:         gin.DebugMode,
		MaxByte:      2048,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		CloseTimeout: time.Second * 10,
		RateTime:     time.Millisecond,
	}
}
