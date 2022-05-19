/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    mini
 * @Date:    2021/8/23 6:39 下午
 * @package: wechat
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package wxmini

import (
	"encoding/json"
	"fmt"
	"git.hawtech.cn/jager/hawox/encrypt"
	"git.hawtech.cn/jager/hawox/errcode"
	"git.hawtech.cn/jager/hawox/logx"
	"io/ioutil"
	"net/http"
)

const baseUrl = "https://api.weixin.qq.com/sns/jscode2session"

type UserSession struct {
	OpenId     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionId    string `json:"unionid"`     // 用户在开放平台的唯一标识符，在满足 UnionID 下发条件的情况下会返回，详见 UnionID 机制说明。
	Errcode    int    `json:"errcode"`     // 错误码
	Errmsg     string `json:"errmsg"`      // 错误信息
}

type UserInfo struct {
	NickName  string `json:"nickName"`  // 昵称
	Gender    int    `json:"gender"`    // 性别
	Language  string `json:"language"`  // 语言
	City      string `json:"city"`      // 城市
	Province  string `json:"province"`  // 省份
	Country   string `json:"country"`   // 国家
	AvatarUrl string `json:"avatarUrl"` // 头像链接
}

func (u *UserInfo) Region() string {
	if u.Country != "" && u.Province != "" && u.City != "" {
		return u.Country + u.Province + u.City
	}

	if u.Province != "" && u.City != "" {
		return u.Province + u.City
	}

	if u.Country != "" && u.Province != "" {
		return u.Country + u.Province
	}

	if u.Province != "" {
		return u.Province
	}

	if u.City != "" {
		return u.City
	}

	if u.Country != "" {
		return u.Country
	}
	return ""
}

type LoginReply struct {
	UInfo   *UserInfo
	Session *UserSession
}

func getSession(wxAppId, wxAppSecret, code string) (*UserSession, error) {
	url := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		baseUrl, wxAppId, wxAppSecret, code)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reply := &UserSession{}
	err = json.Unmarshal(body, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func getUserInfo(rawData string) (*UserInfo, error) {
	// 用户信息解码
	u := &UserInfo{}
	err := json.Unmarshal([]byte(rawData), u)
	if err != nil {
		logx.Err(err).Msg("login json.Unmarshal")
		return nil, err
	}
	return u, nil
}

func Login(wxAppId, wxAppSecret, code, rawData, signature string) (*LoginReply, error) {
	uSess, err := getSession(wxAppId, wxAppSecret, code)
	if err != nil || uSess == nil || uSess.Errcode != 0 {
		if uSess == nil {
			return nil, err
		}
		return nil, errcode.New(int32(uSess.Errcode), uSess.Errmsg)
	}
	sha1Str := fmt.Sprintf("%s%s", rawData, uSess.SessionKey)
	signature2 := encrypt.Sha1(sha1Str)
	if signature2 != signature {
		return nil, errcode.InvalidParam
	}
	uInfo, err := getUserInfo(rawData)
	if err != nil {
		return nil, err
	}
	reply := &LoginReply{
		UInfo:   uInfo,
		Session: uSess,
	}
	return reply, nil
}

func GetOpenId(wxAppId, wxAppSecret, code string) (openId string, err error) {
	uSess, err1 := getSession(wxAppId, wxAppSecret, code)
	if err1 != nil || uSess == nil || uSess.Errcode != 0 {
		if uSess == nil {
			err = err1
			return
		}
		logx.Error().Int("ErrCode", uSess.Errcode).Str("ErrMsg", uSess.Errmsg).Msg("GetMiniOpenId")
		err = errcode.New(int32(uSess.Errcode), uSess.Errmsg)
		return
	}
	openId = uSess.OpenId
	return
}

func GetUnionId(appId, appSecret, code string) (unionId string, err error) {
	uSess, err1 := getSession(appId, appSecret, code)
	if err1 != nil || uSess == nil || uSess.Errcode != 0 {
		if uSess == nil {
			err = err1
			return
		}
		err = errcode.New(int32(uSess.Errcode), "")
		return
	}
	unionId = uSess.UnionId
	return
}
