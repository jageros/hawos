/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    constellation
 * @Date:    2021/9/20 12:16 上午
 * @package: constellation
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package astro

import (
	"github.com/nosixtools/solarlunar"
	"strconv"
	"strings"
)

var (
	unknownConstellation = "未知星座"
	constellations       = []string{"摩羯座", "水瓶座", "双鱼座", "白羊座", "金牛座", "双子座", "巨蟹座", "狮子座", "处女座", "天秤座", "天蝎座", "射手座", "摩羯座"}
	splitDays            = []int{20, 19, 21, 20, 21, 22, 23, 23, 23, 24, 23, 22}
)

func GetAnimal(date string) string {
	s := solarlunar.SolarToChineseLuanr(date)
	return s[6:9]
}

func GetConstellation(date string) string {
	strs := strings.Split(date, "-")
	if len(strs) < 3 {
		return unknownConstellation
	}
	month, err := strconv.Atoi(strs[1])
	if err != nil {
		return unknownConstellation
	}
	day, err := strconv.Atoi(strs[2])
	if err != nil {
		return unknownConstellation
	}
	if month < 1 || month > 12 || day < 1 || day > 31 {
		return unknownConstellation
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
