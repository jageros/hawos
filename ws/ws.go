/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    service
 * @Date:    2021/7/8 5:30 下午
 * @package: service
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/jager/hawox/contextx"
	"github.com/jager/hawox/errcode"
	"github.com/jager/hawox/httpx"
	"gopkg.in/olahol/melody.v1"
	"time"
)

var ss *service

type Options interface {
	SetCallTimeout(t time.Duration)
	SetRelativePath(relativePath string)
	SetKeys(keys ...string)
	SetAuth(auth func(c *gin.Context))

	HandleConnect(fn func(*melody.Session))
	HandleDisconnect(fn func(*melody.Session))
	HandlePong(fn func(*melody.Session))
	HandleMessage(fn func(*melody.Session, []byte))
	HandleMessageBinary(fn func(*melody.Session, []byte))
	HandleSentMessage(fn func(*melody.Session, []byte))
	HandleSentMessageBinary(fn func(*melody.Session, []byte))
	HandleError(fn func(*melody.Session, error))
	HandleClose(fn func(*melody.Session, int, string) error)
}

type service struct {
	*melody.Melody
	ctx          contextx.Context
	callTimeout  time.Duration
	relativePath string
	keys         []string
	auth         func(c *gin.Context)
}

type Filter struct {
	isRefused bool
	Keys      map[string]map[interface{}]struct{}
}

func (f *Filter) AddKeys(key string, vals ...interface{}) {
	if f.Keys == nil {
		f.Keys = map[string]map[interface{}]struct{}{
			key: {},
		}
	} else if _, ok := f.Keys[key]; !ok {
		f.Keys[key] = map[interface{}]struct{}{}
	}
	for _, v := range vals {
		f.Keys[key][v] = struct{}{}
	}
}

func newFilter(isRefused bool, key string, vals ...interface{}) *Filter {
	filter := &Filter{
		isRefused: isRefused,
		Keys: map[string]map[interface{}]struct{}{
			key: {},
		},
	}
	for _, v := range vals {
		filter.Keys[key][v] = struct{}{}
	}
	return filter
}

func NewBlacklistFilter(key string, vals ...interface{}) *Filter {
	return newFilter(true, key, vals...)
}

func NewWhitelistFilter(key string, vals ...interface{}) *Filter {
	return newFilter(false, key, vals...)
}

func (s *service) SetCallTimeout(t time.Duration) {
	s.callTimeout = t
}
func (s *service) SetRelativePath(relativePath string) {
	s.relativePath = relativePath
}
func (s *service) SetKeys(keys ...string) {
	s.keys = keys
}

func (s *service) SetAuth(auth func(c *gin.Context)) {
	s.auth = auth
}

func Broadcast(data []byte, filter *Filter) error {
	if filter == nil {
		return ss.Broadcast(data)
	}

	return ss.BroadcastFilter(data, func(session *melody.Session) bool {

		if filter.isRefused {
			for k, v := range session.Keys {
				if vs, ok := filter.Keys[k]; ok {
					if _, ok := vs[v]; ok {
						return false
					}
				}
			}
			return true
		}

		for k, v := range session.Keys {
			if vs, ok := filter.Keys[k]; ok {
				if _, ok := vs[v]; ok {
					return true
				}
			}
		}
		return false

	})
}

func BroadcastBinary(data []byte, filter *Filter) error {
	if filter == nil {
		return ss.BroadcastBinary(data)
	}

	return ss.BroadcastBinaryFilter(data, func(session *melody.Session) bool {

		if filter.isRefused {
			for k, v := range session.Keys {
				if vs, ok := filter.Keys[k]; ok {
					if _, ok := vs[v]; ok {
						return false
					}
				}
			}
		} else {
			for k, v := range session.Keys {
				if vs, ok := filter.Keys[k]; ok {
					if _, ok := vs[v]; ok {
						return true
					}
				}
			}
		}
		return false
	})
}

func Init(ctx contextx.Context, r *gin.RouterGroup, opfs ...func(opt Options)) {
	ss = &service{
		Melody:       melody.New(),
		ctx:          ctx,
		callTimeout:  time.Second * 5,
		relativePath: "/ws",
	}
	for _, opf := range opfs {
		opf(ss)
	}
	if ss.auth != nil {
		r.GET(ss.relativePath, ss.auth, ss.handler)
	} else {
		r.GET(ss.relativePath, ss.handler)
	}

}

func (s *service) handler(c *gin.Context) {
	var keys = map[string]interface{}{}
	for _, k := range s.keys {
		var vi interface{}
		v := c.GetHeader(k)
		if v != "" {
			vi = v
		} else {
			v, ok := httpx.DecodeUrlVal(c, k)
			if ok {
				vi = v
			} else {
				vi, _ = c.Get(k)
			}
		}
		if vi != nil {
			keys[k] = vi
		}
	}

	err := s.HandleRequestWithKeys(c.Writer, c.Request, keys)
	if err != nil {
		httpx.ErrInterrupt(c, errcode.WithErrcode(-1, err))
	}
}
