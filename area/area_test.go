/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    area_test
 * @Date:    2021/12/6 1:47 下午
 * @package: area
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package area

import (
	"fmt"
	"testing"
)

func Test_Areas(t *testing.T) {
	as := GetCounties("广东省", "广州市")
	for _, a := range as {
		fmt.Println(a.GetName(), a.GetCode())
	}
	//fmt.Printf("广东省广州市区列表：%v\n\n", as.Strings())
}
