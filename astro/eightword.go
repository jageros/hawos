/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    eightword
 * @Date:    2021/12/14 5:39 下午
 * @package: astro
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package astro

import "fmt"

type EightWord struct {
	hour int
	*Date
}

func (ew *EightWord) Word() string {
	hGanZhi, _ := newGanZhi(ew.DayGanZhi[:3], ew.hour)
	return fmt.Sprintf("%s%s%s%s", ew.YearGanZhi, ew.MonthGanZhi, ew.DayGanZhi, hGanZhi)
}

func (ew *EightWord) Words() []string {
	hGanZhi, _ := newGanZhi(ew.DayGanZhi[:3], ew.hour)
	words := []string{ew.YearGanZhi[:3], ew.YearGanZhi[3:], ew.MonthGanZhi[:3], ew.MonthGanZhi[3:], ew.DayGanZhi[:3], ew.DayGanZhi[3:], hGanZhi[:3], hGanZhi[3:]}
	return words
}

func (ew *EightWord) HourGanZhi() string {
	hGanZhi, _ := newGanZhi(ew.DayGanZhi[:3], ew.hour)
	return hGanZhi
}

func (ew *EightWord) WuXingAttr() []string {
	return wuXingAttrs(ew.Words())
}

func (ew *EightWord) EWString() string {
	return fmt.Sprintf("%s %s %s %s月%s 八字：%v 五行：%v 缺：%v", ew.Key(), ew.Constellation(), ew.Animal, ew.LunarMonth, ew.LunarDay, ew.Words(), ew.WuXingAttr(), ew.MissWuXingAttr())
}

func (ew *EightWord) MissWuXingAttr() []string {
	return missWuXing(ew.WuXingAttr())
}
