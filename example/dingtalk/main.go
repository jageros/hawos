package main

import (
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/sdk/dingtalk"
)

func main() {
	err := dingtalk.SendMsg("what the fuck!\n啦啦啦啦\n啊哈哈哈哈~")
	if err != nil {
		logx.Error(err)
	}
}
