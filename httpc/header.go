/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    header
 * @Date:    2021/8/18 6:22 下午
 * @package: httpc
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package httpc

import "net/http"

func SetHeader(req *http.Request, arg map[string]string) {
	for key, val := range arg {
		req.Header.Set(key, val)
	}
}
