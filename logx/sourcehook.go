/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    filehook
 * @Date:    2022/3/21 5:23 PM
 * @package: logs
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package logx

import (
	"github.com/rs/zerolog"
)

type sourcehook struct {
	source string
}

func (sh *sourcehook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if sh.source != "" {
		e.Str("source", sh.source)
	}
}
