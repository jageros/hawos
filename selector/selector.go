/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    selector
 * @Date:    2021/7/9 10:47 上午
 * @package: selector
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package selector

import (
	"sync"
)

var (
	msgId2names map[int32]map[string]struct{}
	rw          sync.RWMutex
)

func init() {
	msgId2names = map[int32]map[string]struct{}{}
}

func Update(mds []*metaData) {
	rw.Lock()
	defer rw.Unlock()

	msgId2names = map[int32]map[string]struct{}{}
	for _, md := range mds {
		for _, id := range md.MsgIds {
			if mn, ok := msgId2names[id]; ok {
				mn[md.Name] = struct{}{}
			} else {
				msgId2names[id] = map[string]struct{}{md.Name: {}}
			}
		}
	}
}

func GetName(msgId int32) []string {
	rw.RLock()
	defer rw.RUnlock()
	var names []string
	if ms, ok := msgId2names[msgId]; ok {
		for n, _ := range ms {
			names = append(names, n)
		}
	}
	return names
}
