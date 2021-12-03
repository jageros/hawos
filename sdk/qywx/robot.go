/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    robot
 * @Date:    2021/12/3 1:44 下午
 * @package: qywx
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package qywx

import (
	"encoding/json"
	"github.com/jageros/hawox/errcode"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	robotUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=4f8cd479-e8a6-4c5b-98fe-4fa8eadce176"
)

func SetRobotConfig(url string) {
	robotUrl = url
}

type respMsg struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	InvalidUser string `json:"invaliduser"`
}

type textMsg struct {
	Content       string   `json:"content"`
	MentionedList []string `json:"mentioned_list"`
	MobileList    []string `json:"mentioned_mobile_list"`
}

type groupMsg struct {
	MsgType string  `json:"msgtype"`
	Text    textMsg `json:"text"`
}

func SendMsg(msg string) error {
	gMsg := &groupMsg{
		MsgType: "text",
		Text: textMsg{
			Content: msg,
		},
	}
	body, err := json.Marshal(gMsg)
	if err != nil {
		return err
	}
	bodyReader := strings.NewReader(string(body))
	resp, err := http.Post(robotUrl, "application/json;charset=UTF-8", bodyReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	reply := &respMsg{}
	err = json.Unmarshal(respBody, reply)
	if err != nil {
		return err
	}
	if reply.Errcode != 0 {
		return errcode.New(int32(reply.Errcode), reply.Errmsg)
	}
	return nil
}
