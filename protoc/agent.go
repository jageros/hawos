/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    agent
 * @Date:    2021/6/16 3:28 下午
 * @package: wsapp
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package protoc

import (
	"sync"
)

var agentPool *sync.Pool

func init() {
	agentPool = &sync.Pool{New: func() interface{} {
		return new(Agent)
	}}
}

type Agent struct {
	Uid     string
}

func getAgent(uid string) *Agent {
	ag := agentPool.Get().(*Agent)
	ag.Uid = uid
	return ag
}

func putAgent(ag *Agent) {
	agentPool.Put(ag)
}
