# attribute

#### 介绍

对MongoDB数据读写的封装，通过键值的结构存储数据

#### 安装教程

``go get github.com/jageros/hawox``

#### 使用例子

1. 存储数据

```go
package main

import (
	"github.com/jageros/hawox/attribute"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/logx"
)

func main() {
	ctx, cancel := contextx.Default()
	defer cancel()
	attribute.Initialize(ctx, func(opt *attribute.Option) {
		opt.Addr = "127.0.0.1:27017"
		opt.DBName = "attribute"
	})

	attr := attribute.NewAttrMgr("attr", "test")

	attrMap := attribute.NewMapAttr()
	attrMap.SetStr("name", "jackss")

	attr.SetMapAttr("info", attrMap)

	err := attr.Save(true)
	if err != nil {
		logx.Error(err)
	} else {
		logx.Info("save successful!")
	}

	err = ctx.Wait()
	logx.Infof("Stop With: %v", err)
}
```

2. 读取数据

```go
package main

import (
	"github.com/jageros/hawox/attribute"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/logx"
)

func main() {
	ctx, cancel := contextx.Default()
	defer cancel()
	attribute.Initialize(ctx, func(opt *attribute.Option) {
		opt.Addr = "127.0.0.1:27017"
		opt.DBName = "attribute"
	})

	attr := attribute.NewAttrMgr("attr", "test")
	err := attr.Load(true)
	if err != nil {
		logx.Error(err)
	} else {
		attrMap := attr.GetMapAttr("info")
		name := attrMap.GetStr("name")
		logx.Infof("Get Name=%s", name)
	}

	err = ctx.Wait()
	logx.Infof("Stop With: %v", err)
}
```

#### 使用说明

1. xxxx
2. xxxx
3. xxxx