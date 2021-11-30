/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    tianapi
 * @Date:    2021/11/30 14:31
 * @package: sdk
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package tianapi

import (
	"fmt"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/httpc"
)

var (
	baseUrl = "http://api.tianapi.com/jiejiari/index"
	key     = "26feed49ca7ba5c62adc63ad08e96f3f"
)

type IDateType interface {
	Type() DateType
	String() string
	Info() string
}

type DateType uint8

const (
	UnknownType DateType = 0 // 未知类型
	WorkingDay  DateType = 1 // 工作日
	Holiday     DateType = 2 // 节假日
	Weekends    DateType = 3 // 双休日
	TakeWorking DateType = 4 // 调休日
)

func (dt DateType) Type() DateType {
	return dt
}

func (dt DateType) String() string {
	switch dt {
	case WorkingDay:
		return "工作日"
	case Holiday:
		return "节假日"
	case Weekends:
		return "双休日"
	case TakeWorking:
		return "调休日"
	}
	return "未知类型"
}

func (dt DateType) Info() string {
	return dt.String()
}

type Date struct {
	ty   DateType
	info string
}

func (d *Date) Type() DateType {
	return d.ty
}

func (d *Date) String() string {
	return d.ty.String()
}

func (d *Date) Info() string {
	return d.info
}

func newDate(dty, info string) *Date {
	ty := UnknownType
	switch {
	case dty == "调休":
		ty = TakeWorking
	case dty == "工作日":
		ty = WorkingDay
	case dty == "双休日":
		ty = Weekends
	case dty == "节假日":
		ty = Holiday
	}
	return &Date{
		ty:   ty,
		info: info,
	}
}

type dateMsg struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	NewsList []struct {
		Date       string   `json:"date"`
		DayCode    int      `json:"daycode"`
		Weekday    int      `json:"weekday"`
		CnWeekday  string   `json:"cnweekday"`
		LunarYear  string   `json:"lunaryear"`
		LunarMonth string   `json:"lunarmonth"`
		LunarDay   string   `json:"lunarday"`
		Info       string   `json:"info"`
		Start      int      `json:"start"`
		Now        int      `json:"now"`
		End        int      `json:"end"`
		Holiday    string   `json:"holiday"`
		Name       string   `json:"name"`
		EnName     string   `json:"enname"`
		IsNotWork  int      `json:"isnotwork"`
		Vacation   []string `json:"vacation"`
		Remark     []string `json:"remark"`
		Wage       int      `json:"wage"`
		Tip        string   `json:"tip"`
		Rest       string   `json:"rest"`
	} `json:"newslist"`
}

func CheckDateType(date string) (IDateType, error) {
	url := fmt.Sprintf("%s?key=%s&date=%s", baseUrl, key, date)
	var result = new(dateMsg)
	err := httpc.RequestWithInterface(httpc.GET, url, httpc.FORM, nil, nil, result)
	if err != nil {
		return UnknownType, err
	}
	if result.Code != 200 {
		return UnknownType, errcode.New(int32(result.Code), result.Msg)
	}
	if len(result.NewsList) <= 0 {
		return UnknownType, nil
	}
	dd := result.NewsList[0]

	if dd.Tip == "调休" {
		return newDate(dd.Tip, dd.Name), nil
	}

	return newDate(dd.Info, dd.Name), nil
}
