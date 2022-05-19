/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    ginhook
 * @Date:    2022/3/22 3:58 PM
 * @package: logs
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package logx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type nilWrite struct {
}

func (nw *nilWrite) Write(p []byte) (n int, err error) {
	return
}

func GinLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: ginLogFormatter,
		Output:    &nilWrite{},
	})
}

// ginLogFormatter is the log format function Logger middleware uses.
var ginLogFormatter = func(params gin.LogFormatterParams) string {
	lgev := Info()
	if params.StatusCode != http.StatusOK {
		lgev = Warn()
	}
	if params.ErrorMessage != "" {
		lgev = Warn().Str("gin-err", params.ErrorMessage)
	}
	lgev.Str("method", params.Method).Int("code", params.StatusCode).
		Str("ip", params.ClientIP).Str("path", params.Path).
		Str("take", params.Latency.String()).Send()

	return ""
}

//
//func (l *loggerSt) Write(p []byte) (n int, err error) {
//	n = len(p)
//	arg := paramPool.Get().(map[string]interface{})
//	err = json.Unmarshal(p, &arg)
//	if err != nil {
//		return
//	}
//	ev := l.lg.Info()
//	if arg["gin-err"].(string) != "" {
//		ev = l.lg.Warn()
//	}
//
//	for k, v := range arg {
//		switch v.(type) {
//		case float64:
//			ev = ev.Float64(k, v.(float64))
//		case string:
//			if v != "" {
//				ev = ev.Str(k, v.(string))
//			}
//
//		case map[string]interface{}:
//			if v != nil {
//				var d []byte
//				d, err = json.Marshal(v)
//				if err != nil {
//					return
//				}
//				ev = ev.RawJSON(k, d)
//			}
//
//		default:
//			if v != nil {
//				ev = ev.Str(k, fmt.Sprintf("%v", v))
//			}
//		}
//	}
//	paramPool.Put(arg)
//	ev.Send()
//	return
//}
