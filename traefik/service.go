/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    service
 * @Date:    2022/3/28 6:59 PM
 * @package: traefik
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package traefik

import "fmt"

type API struct {
	Router            string   // 路由名称
	RouterEntryPoints []string // 入口端口
	RouterRule        string   // 路由规则
	Service           string   // 服务名称
	Url               string   // 服务地址
	PassHostHeader    bool     // 头部穿透
	HealthCheckPath   string   // 健康检测路径
}

func (a *API) KV() map[string]string {
	result := map[string]string{
		fmt.Sprintf(kRule, a.Router):          a.RouterRule,
		fmt.Sprintf(kService, a.Router):       a.Service,
		fmt.Sprintf(kServerUrl, a.Service, 0): a.Url,
	}
	return result
}
