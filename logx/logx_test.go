/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    logs_test
 * @Date:    2022/3/21 4:57 PM
 * @package: logs
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package logx

import (
	"errors"
	"testing"
)

func Test_xxx(t *testing.T) {
	err := Init(func(opt *Option) {
		opt.LogPath = "./logfile.log"
	})
	err = errors.New("xxx")
	Err(err).Str("uid", "1001").Str("roomid", "r1001").Int64("gid", 1999).Msg("Hello World!")
	Sync()
}
