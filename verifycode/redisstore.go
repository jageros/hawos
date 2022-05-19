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
)

type redisStore struct {
}

// Set sets the digits for the captcha id.
func (c *redisStore) Set(id string, value string) error {
	return redis.SetString(id, value)
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

//Verify captcha's answer directly
func (c *redisStore) Verify(id, answer string, clear bool) bool {
	value := c.Get(id, clear)
	return value == answer
}
