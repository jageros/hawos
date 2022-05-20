/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    redisstore
 * @Date:    2022/5/16 13:39
 * @package: verifycode
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package verifycode

import (
	"github.com/jageros/hawox/redis"
	"time"
)

type redisStore struct {
	expireTime time.Duration
}

// Set sets the digits for the captcha id.
func (c *redisStore) Set(id string, value string) error {
	return redis.Set(id, value, c.expireTime)
}

// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (c *redisStore) Get(id string, clear bool) string {
	value, _ := redis.GetString(id)
	if clear {
		redis.Del(id)
	}
	return value
}

func (c *redisStore) Del(id string) {
	redis.Del(id)
}

//Verify captcha's answer directly
func (c *redisStore) Verify(id, answer string, clear bool) bool {
	value, _ := redis.GetString(id)
	ok := value == answer
	if ok && clear {
		redis.Del(id)
	}
	return ok
}
