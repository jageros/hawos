/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    session
 * @Date:    2022/1/18 5:35 下午
 * @package: udpx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package udpx

import (
	"net"
	"net/http"
	"sync"
)

type Session struct {
	Response *http.Response
	Keys     map[string]interface{}
	conn     *net.UDPConn
	output   chan *envelope
	melody   *Melody
	open     bool
	rwMutex  *sync.RWMutex
}
