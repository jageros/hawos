/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    qywx
 * @Date:    2021/8/19 2:41 下午
 * @package: qywx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package qywx

import (
	"encoding/json"
	"fmt"
	"github.com/jageros/hawox/logx"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/*
 * 企业微信api
 * 功能描述：给企业微信自建应用发送消息
 * 应用： 服务器崩溃通知
 */

const (
	corpid     = "wwbfcd9334d9b30c7c"                          // 企业id
	agentId    = 1000002                                       // 自建应用id
	corpsecret = "E5NiE8iEfwYSzG8ITxwZ2tw0D6T3rI_BtgXwsTtPoD4" // 自建应用secret

	getTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"                                              // 获取token的url
	sebdMsgUrl  = "https://qyapi.weixin.qq.com/cgi-bin/message/send"                                          // 发送信息url
	groupUrl    = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=4f8cd479-e8a6-4c5b-98fe-4fa8eadce176" // 测试地址

	contentType = "application/json;charset=UTF-8"

	_TokenExpiresErr = 42001
	_InvalidTokenErr = 40014

	timeFormat = "2006-01-02 15:04:05"
)

var (
	accessToken = ""
)

//
//type tokenInfo struct {
//	Errcode int    `json:"errcode"`
//	Errmsg  string `json:"errmsg"`
//	Token   string `json:"access_token"`
//	Expires int    `json:"expires_in"`
//}
//
//type reqMsg struct {
//	ToUsers                string      `json:"touser"`
//	MsgType                string      `json:"msgtype"`
//	AgentID                int         `json:"agentid"`
//	TextCard               TextCardMsg `json:"textcard"`
//	EnableIdTrans          int         `json:"enable_id_trans"`
//	EnableDuplicateCheck   int         `json:"enable_duplicate_check"`
//	DuplicateCheckInterval int         `json:"duplicate_check_interval"`
//}
//
//type TextCardMsg struct {
//	Title       string `json:"title"`
//	Description string `json:"description"`
//	Url         string `json:"url"`
//	BtnTxt      string `json:"btntxt"`
//}
//
//func NewReqMsg(title, msg, skipUrl, channel, version string, users ...string) *reqMsg {
//	user := "@all"
//	if len(users) > 0 {
//		user = users[0]
//	}
//	desc := fmt.Sprintf("<div class=\"gray\">%s</div> <div class=\"normal\">%s</div><div class=\"highlight\">渠道：%s</div><div class=\"highlight\">版本：%s</div>", time.Now().Format(consts.TimeFormat), msg, channel, version)
//	return &reqMsg{
//		ToUsers:                user,
//		MsgType:                "textcard",
//		AgentID:                agentId,
//		EnableIdTrans:          0,
//		EnableDuplicateCheck:   0,
//		DuplicateCheckInterval: 1800,
//		TextCard: TextCardMsg{
//			Title:       title,
//			Description: desc,
//			Url:         skipUrl,
//			BtnTxt:      "详情",
//		},
//	}
//}

type respMsg struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	InvalidUser string `json:"invaliduser"`
}

//func getToken() string {
//	url := fmt.Sprintf("%s?corpid=%s&corpsecret=%s", getTokenUrl, corpid, corpsecret)
//	resp, err := http.Get(url)
//	if err != nil {
//		panic(err)
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		panic(err)
//	}
//	reply := &tokenInfo{}
//	err = json.Unmarshal(body, reply)
//	if err != nil {
//		panic(err)
//	}
//	return reply.Token
//}
//
//func SendMsg(msg *reqMsg) error {
//	if accessToken == "" {
//		accessToken = getToken()
//	}
//	url := fmt.Sprintf("%s?access_token=%s", sebdMsgUrl, accessToken)
//	body, err := json.Marshal(msg)
//	if err != nil {
//		return err
//	}
//	bodyReader := strings.NewReader(string(body))
//	resp, err := http.Post(url, contentType, bodyReader)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//	respBody, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return err
//	}
//	reply := &respMsg{}
//	err = json.Unmarshal(respBody, reply)
//	if err != nil {
//		return err
//	}
//	if reply.Errcode == _InvalidTokenErr || reply.Errcode == _TokenExpiresErr {
//		accessToken = getToken()
//		return SendMsg(msg)
//	}
//	if reply.Errcode != 0 {
//		errmsg := fmt.Sprintf("ErrCode=%d ErrMsg=%s", reply.Errcode, reply.Errmsg)
//		return errors.New(errmsg)
//	}
//	return nil
//}

// ---------- ------------ ------- group msg-------- ------------- -------------

type textMsg struct {
	Content       string   `json:"content"`
	MentionedList []string `json:"mentioned_list"`
	MobileList    []string `json:"mentioned_mobile_list"`
}

type groupMsg struct {
	MsgType string  `json:"msgtype"`
	Text    textMsg `json:"text"`
}

func newGroupMsg(title, content string) *groupMsg {
	msg := fmt.Sprintf("%s\n\n%s\n%s", time.Now().Format(timeFormat), title, content)
	return &groupMsg{
		MsgType: "text",
		Text: textMsg{
			Content:    msg,
			MobileList: []string{"@all"},
		},
	}
}

func PostMsgToGroup(title, content string) {
	msg := newGroupMsg(title, content)
	body, err := json.Marshal(msg)
	if err != nil {
		logx.Infof("PostMsgToGroup json.Marshal err=%v", err)
		return
	}
	url := groupUrl
	bodyReader := strings.NewReader(string(body))
	resp, err := http.Post(url, contentType, bodyReader)
	if err != nil {
		logx.Infof("PostMsgToGroup http.Post err=%v", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logx.Infof("PostMsgToGroup ioutil.ReadAll err=%v", err)
		return
	}
	reply := &respMsg{}
	err = json.Unmarshal(respBody, reply)
	if err != nil {
		logx.Infof("PostMsgToGroup json.Unmarshal err=%v", err)
		return
	}
	if reply.Errcode != 0 {
		errmsg := fmt.Sprintf("ErrCode=%d ErrMsg=%s", reply.Errcode, reply.Errmsg)
		logx.Infof("PostMsgToGroup json.Unmarshal err=%v", errmsg)
	}
	return
}
