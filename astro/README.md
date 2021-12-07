# astro

#### 介绍

通过新历日期获取星座，农历以及生肖

#### 安装教程

``go get github.com/jageros/hawox``

#### 使用例子

```go
package main

import (
	"fmt"
	"github.com/jageros/hawox/astro"
)

func main() {
	date := "1993-08-28"
	nDate := astro.Lunar(date)
	animal := astro.GetAnimal(date)
	constellation := astro.GetConstellation(date)

	fmt.Printf("新历：%s\n", date)
	fmt.Printf("农历：%s\n", nDate)
	fmt.Printf("生肖：%s\n", animal)
	fmt.Printf("星座：%s\n", constellation)
}
```

#### 使用说明

1. astro.Lunar("xxxx-xx-xx") 新历日期转农历
2. astro.GetAnimal("xxxx-xx-xx") 新历日期获取生肖
3. astro.GetConstellation("xxxx-xx-xx") 新历日期获取星座