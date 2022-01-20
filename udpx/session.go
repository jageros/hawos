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
	"sync"
)

type Session struct {
	addr    *net.UDPAddr
	Keys    map[string]interface{}
	readCh chan *Package
	writeCh chan *Package
	open    bool
	rwMutex *sync.RWMutex
}
