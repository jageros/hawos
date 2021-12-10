/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    date
 * @Date:    2021/12/10 5:55 下午
 * @package: astro
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package astro

import (
	"fmt"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/httpc"
	"regexp"
	"strconv"
	"strings"
)

var (
	baseUrl = "https://wannianrili.bmcx.com/%d-%s-01__wannianrili/"

	one    = regexp.MustCompile(`<divclass="wnrl_k_you"id="wnrl_k_you_id_(.*?)">(.*?)</div></div>`)
	date   = regexp.MustCompile(`<divclass="wnrl_k_you_id_biaoti">(.*?)</div>`)
	day    = regexp.MustCompile(`<divclass="wnrl_k_you_id_wnrl_riqi">(.*?)</div>`)
	nongli = regexp.MustCompile(`<divclass="wnrl_k_you_id_wnrl_nongli">(.*?)</div>`)
	ganzhi = regexp.MustCompile(`<divclass="wnrl_k_you_id_wnrl_nongli_ganzhi">(.*?)</div>`)

	jieris = regexp.MustCompile(`<spanclass="wnrl_k_you_id_wnrl_jieri_biaoti">节日</span><spanclass="wnrl_k_you_id_wnrl_jieri_neirong">(.*?)</span>`)
	jieri  = regexp.MustCompile(`<ahref="/(.*?)__jieri/"target="_blank">(.*?)</a>`)

	yis  = regexp.MustCompile(`<spanclass="wnrl_k_you_id_wnrl_yi_biaoti">宜</span><spanclass="wnrl_k_you_id_wnrl_yi_neirong">(.*?)</span>`)
	jis  = regexp.MustCompile(`<spanclass="wnrl_k_you_id_wnrl_ji_biaoti">忌</span><spanclass="wnrl_k_you_id_wnrl_ji_neirong">(.*?)</span>`)
	yiji = regexp.MustCompile(`<ahref="(.*?)"target="_blank"title="(.*?)">(.*?)</a>`)
)

type Date struct {
	Year         int      `json:"year"`
	Month        int      `json:"month"`
	Day          int      `json:"day"`
	LunarMonth   string   `json:"lunar_month"`
	LunarDay     string   `json:"lunar_day"`
	Week         int      `json:"week"`
	IsLargeMonth bool     `json:"is_large_month"`
	Animal       string   `json:"animal"`
	YearGanZhi   string   `json:"year_gan_zhi"`
	MonthGanZhi  string   `json:"month_gan_zhi"`
	DayGanZhi    string   `json:"day_gan_zhi"`
	Festivals    []string `json:"festivals"`
	Suitable     []string `json:"suitable"`
	Avoid        []string `json:"avoid"`
}

func week(w string) int {
	switch {
	case strings.HasSuffix(w, "一"):
		return 1
	case strings.HasSuffix(w, "二"):
		return 2
	case strings.HasSuffix(w, "三"):
		return 3
	case strings.HasSuffix(w, "四"):
		return 4
	case strings.HasSuffix(w, "五"):
		return 5
	case strings.HasSuffix(w, "六"):
		return 6
	case strings.HasSuffix(w, "日"):
		return 7
	}
	return 0
}

func newDate(date, day, lunarDate, ganzhi string, festivals, suitable, avoid []string) (*Date, error) {
	y, err := strconv.Atoi(date[:4])
	if err != nil {
		return nil, err
	}
	m, err := strconv.Atoi(date[7:9])
	if err != nil {
		return nil, err
	}

	d, err := strconv.Atoi(day)
	if err != nil {
		return nil, err
	}

	large := date[13:16] == "大"

	lunarStr := strings.Split(lunarDate, "月")
	dd := &Date{
		Year:         y,
		Month:        m,
		Day:          d,
		LunarMonth:   lunarStr[0],
		LunarDay:     lunarStr[1],
		Week:         week(date[17:26]),
		IsLargeMonth: large,
		Animal:       ganzhi[12:15],
		YearGanZhi:   ganzhi[0:6],
		MonthGanZhi:  ganzhi[21:27],
		DayGanZhi:    ganzhi[30:36],
		Festivals:    festivals,
		Suitable:     suitable,
		Avoid:        avoid,
	}
	return dd, nil
}

func Crawling(year, month int) ([]*Date, error) {
	if year < 1900 || year > 2100 {
		return nil, errcode.New(1, "非法年份")
	}
	if month < 1 || month > 12 {
		return nil, errcode.New(2, "非法月份")
	}
	var m string
	if month < 10 {
		m = fmt.Sprintf("0%d", month)
	} else {
		m = strconv.Itoa(month)
	}
	url := fmt.Sprintf(baseUrl, year, m)
	body, err := httpc.Request(httpc.GET, url, httpc.FORM, nil, nil)
	if err != nil {
		return nil, err
	}

	var result []*Date

	hh := strings.ReplaceAll(string(body), " ", "")
	hh = strings.ReplaceAll(hh, "\n", "")
	hh = strings.ReplaceAll(hh, "\r", "")
	hh = strings.ReplaceAll(hh, "\t", "")

	strs := one.FindAllStringSubmatch(hh, -1)
	for _, ss := range strs {
		var ftv, yis_, jis_ []string
		s1 := date.FindStringSubmatch(ss[2])
		s2 := day.FindStringSubmatch(ss[2])
		s3 := nongli.FindStringSubmatch(ss[2])
		s4 := ganzhi.FindStringSubmatch(ss[2])

		s5 := jieris.FindStringSubmatch(ss[2])
		if len(s5) >= 1 {
			s5ss := jieri.FindAllStringSubmatch(s5[1], -1)
			for _, s5s := range s5ss {
				ftv = append(ftv, s5s[2])
			}
		}

		s6 := yis.FindStringSubmatch(ss[2])
		s6ss := yiji.FindAllStringSubmatch(s6[1], -1)
		for _, s6s := range s6ss {
			yis_ = append(yis_, s6s[3])
		}

		s7 := jis.FindStringSubmatch(ss[2])
		s7ss := yiji.FindAllStringSubmatch(s7[1], -1)
		for _, s7s := range s7ss {
			jis_ = append(jis_, s7s[3])
		}

		dd, err := newDate(s1[1], s2[1], s3[1], s4[1], ftv, yis_, jis_)
		if err == nil {
			result = append(result, dd)
		}
	}
	return result, nil
}
