/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    astro_test
 * @Date:    2021/11/30 15:50
 * @package: astro
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package astro

import (
	"fmt"
	"testing"
)

func Test_Lunar(t *testing.T) {
	//date := "1993-08-28"
	//nDate := Lunar(date)
	//animal := GetAnimal(date)
	//constellation := GetConstellation(date)
	//
	//fmt.Printf("新历日期：%s\n", date)
	//fmt.Printf("农历日期：%s\n", nDate)
	//fmt.Printf("生肖：%s\n", animal)
	//fmt.Printf("星座：%s\n", constellation)
	//s := Horoscope(2021)
	//s1 := GetAnimalByYear(2021)
	fmt.Println(yTianGanDiZhi("1900-08-27"))
}

// 0000 11010101010  0000