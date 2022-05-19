/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    main
 * @Date:    2022/1/21 2:44 下午
 * @package: udp
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package main

import (
	"github.com/jager/hawox/flags"
	"github.com/jager/hawox/logx"
	"github.com/jager/hawox/udpx"
	"net"
)

func main() {
	const appName = "udp-demo"
	ctx, wait := flags.Parse(appName)
	defer wait()

	err := udpx.Init(ctx, func(opt *udpx.Option) {
		opt.OnBinaryHandle = onBinaryMsgMsg
	})
	if err != nil {
		logx.Fatal().Err(err).Msg("udp init")
	}
}

func onBinaryMsgMsg(addr *net.UDPAddr, data []byte) {}
