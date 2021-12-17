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
	d := GetDate("1993-08-27")
	fmt.Println(d.NewEightWord(3).EWString())

	//d2 := GetDate("1996-12-25")
	//for i := 0; i < 24; i += 2 {
	//	fmt.Println(d2.LunarMonth+"月"+d2.LunarDay, d2.Animal, d2.Constellation(), d2.EightWords(i))
	//}
}


func arg(openId, title, updateTime, lastClose, todayOpen, curPrice, rice, max, min, count, total, fiveBuy, fiveSell string) map[string]interface{} {
	return map[string]interface{}{
		"touser":      openId,
		"template_id": "1EaaO8jlUJBrj5y7Iz5_lJVXoQq2ivTxH8JnAM_vqGU",
		"url":         "",

		"miniprogram": map[string]interface{}{
			"appid":    "wxba88e64e7342b027",
			"pagepath": "ab8130a7bf55b78992e3d17f59909e0a",
		},

		"data": map[string]interface{}{
			"first": map[string]interface{}{
				"value": title,
				"color": "#173177",
			},
			"keyword1": map[string]interface{}{
				"value": updateTime,
				"color": "#173177",
			},
			"keyword2": map[string]interface{}{
				"value": lastClose,
				"color": "#173177",
			},
			"keyword3": map[string]interface{}{
				"value": todayOpen,
				"color": "#173177",
			},
			"keyword4": map[string]interface{}{
				"value": curPrice,
				"color": "#173177",
			},
			"keyword5": map[string]interface{}{
				"value": rice,
				"color": "#173177",
			},
			"keyword6": map[string]interface{}{
				"value": max,
				"color": "#173177",
			},
			"keyword7": map[string]interface{}{
				"value": min,
				"color": "#173177",
			},
			"keyword8": map[string]interface{}{
				"value": count,
				"color": "#173177",
			},
			"keyword9": map[string]interface{}{
				"value": total,
				"color": "#173177",
			},
			"keyword10": map[string]interface{}{
				"value": fiveBuy,
				"color": "#173177",
			},
			"keyword11": map[string]interface{}{
				"value": fiveSell,
				"color": "#173177",
			},
			"remark": map[string]interface{}{
				"value": "欢迎再次购买！",
				"color": "#173177",
			},
		},
	}
}

/*
{
           "touser":"OPENID",
           "template_id":"ngqIpbwh8bUfcSsECmogfXcV14J0tQlEpBO27izEYtY",
           "url":"http://weixin.qq.com/download",
           "miniprogram":{
             "appid":"xiaochengxuappid12345",
             "pagepath":"index?foo=bar"
           },
           "data":{
                   "first": {
                       "value":"恭喜你购买成功！",
                       "color":"#173177"
                   },
                   "keyword1":{
                       "value":"巧克力",
                       "color":"#173177"
                   },
                   "keyword2": {
                       "value":"39.8元",
                       "color":"#173177"
                   },
                   "keyword3": {
                       "value":"2014年9月22日",
                       "color":"#173177"
                   },
                   "remark":{
                       "value":"欢迎再次购买！",
                       "color":"#173177"
                   }
           }
       }

*/
