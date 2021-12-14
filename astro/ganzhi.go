/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    ganzhi
 * @Date:    2021/12/13 6:06 下午
 * @package: astro
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package astro

var (
	tianGan = map[int]string{1: "甲", 2: "乙", 3: "丙", 4: "丁", 5: "戊", 6: "己", 7: "庚", 8: "辛", 9: "壬", 0: "癸"}
	//diZhi      = map[int]string{1: "子", 2: "丑", 3: "寅", 4: "卯", 5: "辰", 6: "巳", 7: "午", 8: "未", 9: "申", 10: "酉", 11: "戌", 12: "亥"}
	tianGanNum = map[string]int{"甲": 1, "乙": 2, "丙": 3, "丁": 4, "戊": 5, "己": 6, "庚": 7, "辛": 8, "壬": 9, "癸": 10}
	diZhiNum   = map[string]int{"子": 1, "丑": 2, "寅": 3, "卯": 4, "辰": 5, "巳": 6, "午": 7, "未": 8, "申": 9, "酉": 10, "戌": 11, "亥": 12}
)

func newGanZhi(dayTianGan string, hour int) (ganzhi, animal string) {

	var dizhi string
	switch hour {
	case 23, 0:
		dizhi = "子"
		animal = "鼠"
	case 1, 2:
		dizhi = "丑"
		animal = "牛"
	case 3, 4:
		dizhi = "寅"
		animal = "虎"
	case 5, 6:
		dizhi = "卯"
		animal = "兔"
	case 7, 8:
		dizhi = "辰"
		animal = "龙"
	case 9, 10:
		dizhi = "巳"
		animal = "蛇"
	case 11, 12:
		dizhi = "午"
		animal = "马"
	case 13, 14:
		dizhi = "未"
		animal = "羊"
	case 15, 16:
		dizhi = "申"
		animal = "猴"
	case 17, 18:
		dizhi = "酉"
		animal = "鸡"
	case 19, 20:
		dizhi = "戌"
		animal = "狗"
	case 21, 22:
		dizhi = "亥"
		animal = "猪"
	}

	// 时天干＝(日天干×2＋时地支－2)%10
	stg := tianGan[(tianGanNum[dayTianGan]*2+diZhiNum[dizhi]-2)%10]
	ganzhi = stg + dizhi
	return
}
