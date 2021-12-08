/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    animal
 * @Date:    2021/12/8 6:59 下午
 * @package: astro
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package astro

import "github.com/nosixtools/solarlunar"

var (
	animals = []string{"猴", "鸡", "狗", "猪", "鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊"}
)

func GetAnimal(date string) string {
	s := solarlunar.SolarToChineseLuanr(date)
	return s[6:9]
}

func GetAnimalByYear(year int) string {
	return animals[year%12]
}
