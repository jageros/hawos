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
	for _, p := range ps {
		pName := p.GetName()
		cs := p.GetCities()
		for _, c := range cs {
			cName := c.GetName()
			ccs := c.GetCounties()
			for _, cc := range ccs {
				name := cc.GetName()
				fmt.Println(pName, cName, name)
			}
		}
	}
}