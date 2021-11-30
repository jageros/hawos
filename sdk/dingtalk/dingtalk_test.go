/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    dingtalk_test
 * @Date:    2021/11/30 16:50
 * @package: dingtalk
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package dingtalk

import (
	"github.com/jageros/hawox/logx"
	"testing"
)

func Test_SendMsg(t *testing.T) {
	SetConfig("xxx", "https://oapi.dingtalk.com/robot/send?access_token=xxx")
	err := SendMsg("啦啦啦啦啦啦啦啦啦啦啦啦啦啦啦了~~~")
	if err != nil {
		logx.Error(err)
	}
}
