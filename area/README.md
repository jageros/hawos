# area

#### 介绍

用于填写中国地区时的选择列表，包括省，市，区（县/镇）

#### 安装教程

``go get github.com/jageros/hawox``

#### 使用例子

```go
package main

import (
	"fmt"
	"github.com/jageros/hawox/area"
)

func main() {
	ps := area.GetProvinces()
	fmt.Printf("省份列表：%v\n\n", ps.Strings())

	cs := area.GetCities("广东省")
	fmt.Printf("广东省城市列表：%v\n\n", cs.Strings())

	as := area.GetCounties("广东省", "广州市")
	fmt.Printf("广东省广州市区列表：%v\n\n", as.Strings())
}
```

#### 使用说明

1. GetProvinces() 获取省份列表
2. GetCities(province string) []ICity / (p *Province) GetCities() []ICity 根据省份获取城市列表
3. GetCounties(province, city string) / (c *City) GetCounties() []ICounty 根据省份和城市获取区（县，镇列表）