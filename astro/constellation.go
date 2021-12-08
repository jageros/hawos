/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    constellation
 * @Date:    2021/12/8 6:54 下午
 * @package: astro
 * @Version: v1.0.0
 *
 * @Description: 星座系列
 *
 */

package astro

import (
	"strconv"
	"strings"
)

var (
	constellations = []string{"摩羯座", "水瓶座", "双鱼座", "白羊座", "金牛座", "双子座", "巨蟹座", "狮子座", "处女座", "天秤座", "天蝎座", "射手座", "摩羯座"}
	splitDays      = []int{20, 19, 21, 20, 21, 22, 23, 23, 23, 24, 23, 22}
)

func GetConstellation(date string) string {
	errMsg := "日期格式不正确"
	strs := strings.Split(date, "-")
	if len(strs) < 3 {
		return errMsg
	}
	month, err := strconv.Atoi(strs[1])
	if err != nil {
		return errMsg
	}
	day, err := strconv.Atoi(strs[2])
	if err != nil {
		return errMsg
	}
	if month < 1 || month > 12 || day < 1 || day > 31 {
		return errMsg
	}
	var index = month - 1
	if day >= splitDays[index] {
		index++
	}
	if index >= 12 {
		index = 0
	}
	return constellations[index]
}
