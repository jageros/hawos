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
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/httpx"
	"gopkg.in/olahol/melody.v1"
	"time"
)

var ss *service

type Options interface {
	SetCallTimeout(t time.Duration)
	SetRelativePath(relativePath string)
	SetKeys(keys ...string)
	SetAuth(auth func(c *gin.Context))
}

type service struct {
	ctx          contextx.Context
	m            *melody.Melody
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

func OnMessage(fn func(session *melody.Session, bytes []byte)) {
	ss.m.HandleMessage(fn)
}

func OnMessageBinary(fn func(session *melody.Session, bytes []byte)) {
	ss.m.HandleMessageBinary(fn)
}

func OnConnect(fn func(*melody.Session)) {
	ss.m.HandleConnect(fn)
}

func OnDisConnect(fn func(*melody.Session)) {
	ss.m.HandleDisconnect(fn)
}

func Broadcast(data []byte, filter *Filter) error {
	if filter == nil {
		return ss.m.Broadcast(data)
	}

	return ss.m.BroadcastFilter(data, func(session *melody.Session) bool {

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
		return ss.m.BroadcastBinary(data)
	}

	return ss.m.BroadcastBinaryFilter(data, func(session *melody.Session) bool {

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
		ctx:          ctx,
		m:            melody.New(),
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

	err := s.m.HandleRequestWithKeys(c.Writer, c.Request, keys)
	if err != nil {
		httpx.ErrInterrupt(c, errcode.WithErrcode(-1, err))
	}
}
