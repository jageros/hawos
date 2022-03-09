/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    meta_test
 * @Date:    2022/3/9 10:53 上午
 * @package: template
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package metatemp

import (
	"fmt"
	"testing"
)

func Test_Meta(t *testing.T) {
	ss := []string{"  //@ C2S_AUTH_TOKEN	req: AuthMsg    resp: AuthResp  	", "//@ C2S_PING req: Ping  resp: Pong"}
	s, msgids, err := GenMetaFile("github.com/jageros/hawos/protos/pb", "MsgID", ss)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(s, msgids)
}
