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
	ps := GetProvinces()
	fmt.Printf("省份列表：%v\n\n", ps.Strings())

	cs := GetCities("广东省")
	fmt.Printf("广东省城市列表：%v\n\n", cs.Strings())

	as := GetCounties("广东省", "广州市")
	fmt.Printf("广东省广州市区列表：%v\n\n", as.Strings())
}
