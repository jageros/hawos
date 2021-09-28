/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    bytedance
 * @Date:    2021/9/26 11:34 上午
 * @package: bytedance
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package bytedance

import (
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/httpc"
)

const (
	url    = "https://developer.toutiao.com/api/apps/v2/jscode2session"
	appid  = "xxx"
	secret = "xxx"
)

type Response struct {
	ErrNo   int    `json:"err_no"`
	ErrTips string `json:"err_tips"`
	Data    struct {
		SessionKey      string `json:"session_key"`
		Openid          string `json:"openid"`
		AnonymousOpenid string `json:"anonymous_openid"`
		Unionid         string `json:"unionid"`
	} `json:"data"`
}

func Login(code, anyCode string) (resp Response, err error) {
	arg := map[string]interface{}{
		"appid":          appid,
		"secret":         secret,
		"anonymous_code": anyCode,
		"code":           code,
	}
	err = httpc.Request2(httpc.POST, url, httpc.JSON, arg, nil, &resp)
	return
}

func GetOpenid(code, anyCode string) (openid string, err error) {
	resp, err := Login(code, anyCode)
	if err != nil {
		return "", err
	}
	if resp.ErrNo != 0 {
		return "", errcode.New(int32(resp.ErrNo), resp.ErrTips)
	}
	return resp.Data.Openid, nil
}
