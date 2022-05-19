/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    wechat
 * @Date:    2021/8/23 6:33 下午
 * @package: wechat
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package wxapp

import (
	"encoding/json"
	"fmt"
	"github.com/jageros/hawox/errcode"
	"io/ioutil"
	"net/http"
)

const (
	appid     = ""
	appSecret = ""
)

type UserSession struct {
	AccessToken  string `json:"access_token"`  // 接口调用凭证
	ExpiresIn    int32  `json:"expires_in"`    // access_token 接口调用凭证超时时间，单位（秒）
	RefreshToken string `json:"refresh_token"` // 用户刷新 access_token
	OpenId       string `json:"openid"`        // 授权用户唯一标识
	Scope        string `json:"scope"`         // 用户授权的作用域，使用逗号（,）分隔
	Errcode      int    `json:"errcode"`       // 错误码
}

type UserInfo struct {
	NickName   string   `json:"nickname"`   // 昵称
	Sex        int      `json:"sex"`        // 普通用户性别，1 为男性，2 为女性
	City       string   `json:"city"`       // 城市
	Province   string   `json:"province"`   // 省份
	Country    string   `json:"country"`    // 国家，如中国为 CN
	HeadImgUrl string   `json:"headimgurl"` // 头像
	Privilege  []string `json:"privilege"`  // 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	OpenId     string   `json:"openid"`     // 普通用户的标识，对当前开发者帐号唯一
	UnionId    string   `json:"unionid"`    // 用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的 unionid 是唯一的。
}

type LoginReply struct {
	UInfo   *UserInfo
	Session *UserSession
}

func getSession(code string) (*UserSession, error) {
	sessUrl := "https://api.weixin.qq.com/sns/oauth2/access_token"
	url := fmt.Sprintf("%s?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		sessUrl, appid, appSecret, code)
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

/*
{
  "access_token": "ACCESS_TOKEN",
  "expires_in": 7200,
  "refresh_token": "REFRESH_TOKEN",
  "openid": "OPENID",
  "scope": "SCOPE"
}
*/

// GET https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=APPID&grant_type=refresh_token&refresh_token=REFRESH_TOKEN
func (s *UserSession) refreshToken() error {
	baseUrl := "https://api.weixin.qq.com/sns/oauth2/refresh_token"
	url := fmt.Sprintf("%s?appid=%sgrant_type=refresh_token&refresh_token=%s", baseUrl, appid, s.RefreshToken)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	reply := &UserSession{}
	err = json.Unmarshal(body, reply)
	if err != nil {
		return err
	}
	if reply.Errcode != 0 {
		return errcode.New(int32(reply.Errcode), "")
	}
	s.AccessToken = reply.AccessToken
	s.ExpiresIn = reply.ExpiresIn
	s.RefreshToken = reply.RefreshToken
	s.OpenId = reply.OpenId
	s.Scope = reply.Scope
	return nil
}

// GET https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID
func getUserInfo(accessToken, openid string) (*UserInfo, error) {
	baseUrl := "https://api.weixin.qq.com/sns/userinfo"
	url := fmt.Sprintf("%s?access_token=%s&openid=%s", baseUrl, accessToken, openid)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reply := &UserInfo{}
	err = json.Unmarshal(body, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func Login(code string) (*LoginReply, error) {
	uSess, err := getSession(code)
	if err != nil || uSess == nil || uSess.Errcode != 0 {
		if uSess == nil {
			return nil, err
		}
		return nil, errcode.New(int32(uSess.Errcode), "")
	}

	uInfo, err := getUserInfo(uSess.AccessToken, uSess.OpenId)
	if err != nil {
		return nil, err
	}
	reply := &LoginReply{
		UInfo:   uInfo,
		Session: uSess,
	}
	return reply, nil
}

func GetOpenId(code string) (openId string, err error) {
	uSess, err1 := getSession(code)
	if err1 != nil || uSess == nil || uSess.Errcode != 0 {
		if uSess == nil {
			err = err1
			return
		}
		err = errcode.New(int32(uSess.Errcode), "")
		return
	}
	openId = uSess.OpenId
	return
}

func GetUnionId(code string) (unionId string, err error) {
	resp, err := Login(code)
	if err != nil {
		return "", err
	}
	unionId = resp.UInfo.UnionId
	return
}
