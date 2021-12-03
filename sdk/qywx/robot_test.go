/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    robot_test
 * @Date:    2021/12/3 2:21 下午
 * @package: qywx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package qywx

import (
	"fmt"
	"testing"
)

func Test_SendMsg(t *testing.T) {
	err := SendMsg("啦啦啦啦啦啦！")
	if err != nil {
		fmt.Println(err)
	}
}
