/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    logger
 * @Date:    2021/6/11 10:52 上午
 * @package: http
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package httpx

//
//import (
//	"fmt"
//	"git.hawtech.cn/jager/hawox/logs"
//	"github.com/gin-gonic/gin"
//	"time"
//)
//
//func logger() gin.HandlerFunc {
//	return gin.LoggerWithConfig(gin.LoggerConfig{
//		Formatter: ginLogFormatter,
//		Output:    logs.Writer(),
//	})
//}
//
//// ginLogFormatter is the log format function Logger middleware uses.
//var ginLogFormatter = func(param gin.LogFormatterParams) string {
//	if param.Latency > time.Minute {
//		// Truncate in a golang < 1.8 safe way
//		param.Latency = param.Latency - param.Latency%time.Second
//	}
//	var formatStr string
//	if param.ErrorMessage != "" {
//		formatStr = fmt.Sprintf("Code=%d TakeTime=%v IP=%s Method=%s Path=%v ErrMsg=%s",
//			param.StatusCode,
//			param.Latency,
//			param.ClientIP,
//			param.Method,
//			param.Path,
//			param.ErrorMessage,
//		)
//	} else {
//		formatStr = fmt.Sprintf("Code=%d TakeTime=%v IP=%s Method=%s Path=%v",
//			param.StatusCode,
//			param.Latency,
//			param.ClientIP,
//			param.Method,
//			param.Path,
//		)
//	}
//	return formatStr
//}
