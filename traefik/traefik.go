/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    traefik
 * @Date:    2022/3/28 6:47 PM
 * @package: traefik
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package traefik

type Option struct {
	Router      string
	EntryPoints string
	Service     string
	Rule        string
}
