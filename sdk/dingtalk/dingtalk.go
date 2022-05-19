package dingtalk

import (
	"fmt"
	"git.hawtech.cn/jager/hawox/encrypt"
	"git.hawtech.cn/jager/hawox/errcode"
	"git.hawtech.cn/jager/hawox/httpc"
	"time"
)

var (
	secret_   = "SECff53e03a11a5754e3773d7d0xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	robotUrl_ = "https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
)

func SetConfig(secret, robotUrl string) {
	secret_ = secret
	robotUrl_ = robotUrl
}

type result struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func SendMsg(text string) error {
	timestamp := time.Now().UnixMilli()
	str := fmt.Sprintf("%d\n%s", timestamp, secret_)
	sign := encrypt.HmacSha256Base64(str, secret_)

	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", robotUrl_, timestamp, sign)
	arg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": text,
		},
	}
	resp := new(result)
	err := httpc.RequestWithInterface(httpc.POST, url, httpc.JSON, arg, nil, resp)
	if err != nil {
		return err
	}
	if resp.Errcode != 0 {
		return errcode.New(int32(resp.Errcode), resp.Errmsg)
	}
	return nil
}
