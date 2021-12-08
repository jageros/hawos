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
)



func Lunar(date string) string {
	return solarlunar.SolarToChineseLuanr(date)
}


func Horoscope(year int) string {
	var eightWord string
	eightWord += tianGan[year%10]
	eightWord += diZhi[year%12]
	return eightWord
}
