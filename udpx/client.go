/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    client
 * @Date:    2022/1/21 6:59 下午
 * @package: udpx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package udpx

import (
	"github.com/jageros/hawox/contextx"
	"net"
)

type Client struct {
	TargetAddr string
	conn       *net.UDPConn
}

func NewClient(ctx contextx.Context, opfs ...func(opt *Option)) {

}
