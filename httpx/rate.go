/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    rate
 * @Date:    2021/9/7 10:32 上午
 * @package: httpx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package httpx

import (
	"github.com/gin-gonic/gin"
	"git.hawtech.cn/jager/hawox/errcode"
	"golang.org/x/time/rate"
	"time"
)

func RateMiddleware(rateTime time.Duration) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(rateTime), int(time.Second/rateTime))
	return func(c *gin.Context) {
		err := limiter.Wait(c)
		if err != nil {
			ErrInterrupt(c, errcode.Overload)
			return
		}
		c.Next()
	}
}
