/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    envelope
 * @Date:    2022/1/18 5:36 下午
 * @package: udpx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package udpx

type envelope struct {
	t      int
	msg    []byte
	filter filterFunc
}
