package zlog

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func TestLog(t *testing.T) {
	lg := NewLogger(func(opt *Option) {
		opt.FileOut = true
		opt.Path = "logs/server.log"
		opt.MaxFileSize = 1
		opt.MaxBackups = 3
		opt.MaxAge = 1
		opt.Compress = true
		opt.Source = "test"
	})
	err := errors.New("xxxx")
	lg.Namespace("cccc").With("appid", 43245234234).With("sign", "dfasdfsdfsf").With("opt", &Option{
		Path:        "",
		Level:       "",
		MaxFileSize: 0,
		MaxBackups:  0,
		MaxAge:      0,
		Compress:    false,
		Caller:      false,
		StdOut:      false,
		FileOut:     false,
		Source:      "",
	}).With("err", err).Error("dddd")
	fmt.Println(lg.Sync())
}
